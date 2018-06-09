package routesTest

import (
	"esports-league-manager/Backend/Server/routes"
	"github.com/gin-gonic/gin"
	"testing"
	"github.com/kataras/iris/core/errors"
)

type userCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var router *gin.Engine


//set up mock user daos
type mockUsersDAO struct {
	t *testing.T
	e error
}
func (u *mockUsersDAO) InsertUser(email, salt, hash string) error {
	return nil
}
func (u *mockUsersDAO) IsEmailInUse(email string) (bool, error) {
	return true, u.e
}


type mockUsersDAOCreateUser struct {
	t *testing.T
	UserCreated bool
}
func (u *mockUsersDAOCreateUser) InsertUser(email, salt, hash string) error {
	if len(salt) != 64 {
		u.t.Errorf("Salt is incorrect length. Got %v, expected 64", len(salt))
	}
	if len(hash) != 128 {
		u.t.Errorf("Hash is incorrect length. Got %v, expected 128", len(salt))
	}
	u.UserCreated = true
	return nil
}
func (u *mockUsersDAOCreateUser) IsEmailInUse(email string) (bool, error) {
	return false, nil
}


func testPasswordTooShort(t *testing.T, pass string) {
	jsonErrorTest(t, "test@test.com", pass, "passwordTooShort")
}

func testMalformedEmail(t *testing.T, email string) {
	jsonErrorTest(t, email, "abcd1234", "emailMalformed")
}

func testEmailInUse(t *testing.T) {
	jsonErrorTest(t, "test@test.com", "abcd1234", "emailInUse")
}

func testDatabaseEror(t *testing.T) {
	routes.UsersDAO = &mockUsersDAO{
		t: t,
		e: errors.New("fake database error"),
	}
	responseCodeTest(t, "test@test.com", "abcd1234", 500)
}

func testCorrectUserCreation(t *testing.T) {
	mockDAO := &mockUsersDAOCreateUser{
		t: t,
		UserCreated: false,
	}
	routes.UsersDAO = mockDAO

	responseCodeTest(t, "test@test.com", "abcd1234", 200)
	if !mockDAO.UserCreated {
		t.Error("User creation DAO function was not called")
	}
}

func Test_CreateNewUser(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_createNewUser)
	routes.UsersDAO = &mockUsersDAO{
		t: t,
		e: nil,
	}

	t.Run("passwordTooShort=Length7A", func(t *testing.T) {
		testPasswordTooShort(t, "1234567")
	})
	t.Run("passwordTooShort=Length=7B", func(t *testing.T) {
		testPasswordTooShort(t, "123456 ")
	})
	t.Run("passwordTooShortLength=0", func(t *testing.T) {
		testPasswordTooShort(t, "")
	})

	//these do not have to be comprehensive as they are handled by a third party lib
	t.Run("malformedEmail=InvalidCharacters", func(t *testing.T) {
		testMalformedEmail(t, "ç$€§/az@gmail.com")
	})
	t.Run("malformedEmail=NoName", func(t *testing.T) {
		testMalformedEmail(t, "@gmail.com")
	})

	t.Run("emailInUse", testEmailInUse)
	t.Run("databaseError", testDatabaseEror)
	t.Run("correctUserCreation", testCorrectUserCreation)
}
