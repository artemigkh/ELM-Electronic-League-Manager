package dataHelpersTest

import (
	"Server/databaseAccess"
	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"mocks"
	"testing"
)

type validateLeague struct {
	suite.Suite
	mockLeaguesDao *mocks.LeaguesDAO
	assert         *assert.Assertions
}

func (s *validateLeague) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *validateLeague) SetupTest() {
	s.mockLeaguesDao = new(mocks.LeaguesDAO)
	s.mockLeaguesDao.On("IsNameInUse", 0, mock.AnythingOfType("string")).
		Return(false, nil)
	databaseAccess.Leagues = s.mockLeaguesDao
}

func (s *validateLeague) TestNameTooShort() {
	league := databaseAccess.LeagueCore{
		Name: "a",
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.NameTooShort)
}

func (s *validateLeague) TestNameTooLong() {
	league := databaseAccess.LeagueCore{
		Name: randomdata.RandStringRunes(51),
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.NameTooLong)
}

func (s *validateLeague) TestNameNotUnique() {
	name := randomdata.RandStringRunes(50)
	s.mockLeaguesDao = new(mocks.LeaguesDAO)
	s.mockLeaguesDao.On("IsNameInUse", 17, name).Return(true, nil)
	databaseAccess.Leagues = s.mockLeaguesDao
	league := databaseAccess.LeagueCore{
		Name: name,
	}

	valid, problem, err := league.Validate(17)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.LeagueNameInUse)
	s.mockLeaguesDao.AssertExpectations(s.T())
}

func (s *validateLeague) TestNameDbError() {
	name := randomdata.RandStringRunes(50)
	s.mockLeaguesDao = new(mocks.LeaguesDAO)
	s.mockLeaguesDao.On("IsNameInUse", 0, name).
		Return(false, errors.New("fake db error"))
	databaseAccess.Leagues = s.mockLeaguesDao
	league := databaseAccess.LeagueCore{
		Name: name,
	}

	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Equal(err.Error(), "fake db error")
	s.assert.Equal(problem, "")
}

func (s *validateLeague) TestNameDescriptionTooLong() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(501),
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.DescriptionTooLong)
}

func (s *validateLeague) TestGameInvalid() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(0),
		Game:        "nonexistentgame",
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.InvalidGame)
}

func (s *validateLeague) TestLeaguePermissionsWrong() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(500),
		Game:        "genericsport",
		PublicView:  false,
		PublicJoin:  true,
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.LeaguePermissionsWrong)
}

func (s *validateLeague) TestLeagueTimestampsSignup() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(500),
		Game:        "genericsport",
		PublicView:  true,
		PublicJoin:  true,
		SignupStart: 2,
		SignupEnd:   1,
		LeagueStart: 3,
		LeagueEnd:   4,
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.TimeOutOfOrder)
}

func (s *validateLeague) TestLeagueTimestampsCompetition() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(500),
		Game:        "genericsport",
		PublicView:  true,
		PublicJoin:  true,
		SignupStart: 1,
		SignupEnd:   2,
		LeagueStart: 4,
		LeagueEnd:   3,
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.TimeOutOfOrder)
}

func (s *validateLeague) TestLeagueTimestampsWrongPeriod() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(500),
		Game:        "genericsport",
		PublicView:  true,
		PublicJoin:  true,
		SignupStart: 3,
		SignupEnd:   4,
		LeagueStart: 1,
		LeagueEnd:   2,
	}
	valid, problem, err := league.Validate(0)

	s.assert.False(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, databaseAccess.PeriodOutOfOrder)
}

func (s *validateLeague) TestLeagueCorrect() {
	league := databaseAccess.LeagueCore{
		Name:        randomdata.RandStringRunes(2),
		Description: randomdata.RandStringRunes(500),
		Game:        "basketball",
		PublicView:  false,
		PublicJoin:  false,
		SignupStart: 1,
		SignupEnd:   2,
		LeagueStart: 3,
		LeagueEnd:   4,
	}
	valid, problem, err := league.Validate(0)

	s.assert.True(valid)
	s.assert.Nil(err)
	s.assert.Equal(problem, "")
}

func TestValidateLeague(t *testing.T) {
	suite.Run(t, new(validateLeague))
}
