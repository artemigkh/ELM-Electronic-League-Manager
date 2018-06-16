package routesTest

import (
	"esports-league-manager/Backend/Server/routes"
	"github.com/gin-gonic/gin"
	"testing"
	"github.com/kataras/iris/core/errors"
	"bytes"
	"encoding/json"
)

type userCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUserCreateRequestBody(email, pass string) *bytes.Buffer {
	reqBody := userCreateRequest{
		Email:    email,
		Password: pass,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

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
func (u *mockUsersDAO) GetAuthenticationInformation(email string) (int, string, string, error) {
	return 0, "", "", nil
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
func (u *mockUsersDAOCreateUser) GetAuthenticationInformation(email string) (int, string, string, error) {
	return 0, "", "", nil
}


func testCreateNewUserPasswordTooShort(t *testing.T, pass string) {
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", pass),
		"passwordTooShort", "POST", 400)
}

func testCreateNewUserMalformedEmail(t *testing.T, email string) {
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody(email, "abcd1234"),
		"emailMalformed", "POST", 400)
}

func testCreateNewUserEmailInUse(t *testing.T) {
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
		"emailInUse", "POST", 400)
}

func testCreateNewUserDatabaseError(t *testing.T) {
	routes.UsersDAO = &mockUsersDAO{
		t: t,
		e: errors.New("fake database error"),
	}
	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
		500, "POST")
}

func testCorrectUserCreation(t *testing.T) {
	mockDAO := &mockUsersDAOCreateUser{
		t: t,
		UserCreated: false,
	}
	routes.UsersDAO = mockDAO

	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
		200, "POST")
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
		testCreateNewUserPasswordTooShort(t, "1234567")
	})
	t.Run("passwordTooShort=Length=7B", func(t *testing.T) {
		testCreateNewUserPasswordTooShort(t, "123456 ")
	})
	t.Run("passwordTooShortLength=0", func(t *testing.T) {
		testCreateNewUserPasswordTooShort(t, "")
	})

	//these do not have to be comprehensive as they are handled by a third party lib
	t.Run("malformedEmail=InvalidCharacters", func(t *testing.T) {
		testCreateNewUserMalformedEmail(t, "ç$€§/az@gmail.com")
	})
	t.Run("malformedEmail=NoName", func(t *testing.T) {
		testCreateNewUserMalformedEmail(t, "@gmail.com")
	})

	t.Run("emailInUse", testCreateNewUserEmailInUse)
	t.Run("databaseError", testCreateNewUserDatabaseError)
	t.Run("correctUserCreation", testCorrectUserCreation)
}
