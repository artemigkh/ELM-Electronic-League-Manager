package routesTest

import (
	"esports-league-manager/Backend/Server/routes"
	"github.com/gin-gonic/gin"
	"testing"
	"bytes"
	"encoding/json"
	"esports-league-manager/mocks"
	"github.com/stretchr/testify/mock"
	"errors"
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

func testCreateNewUserMalformedBody(t *testing.T) {
	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})
}

func testCreateNewUserPasswordTooShort(t *testing.T, pass string) {
	httpTest(t, createUserCreateRequestBody("test@test.com", pass), "POST", "/",
		400, testParams{Error: "passwordTooShort"})
}

func testCreateNewUserMalformedEmail(t *testing.T, email string) {
	httpTest(t, createUserCreateRequestBody(email, "12345678"), "POST", "/",
		400, testParams{Error: "emailMalformed"})
}

func testCreateNewUserEmailInUse(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 400, testParams{Error: "emailInUse"})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func testCreateNewUserDatabaseError(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(false, nil)
	mockUsersDao.On("CreateUser", "test@test.com", mock.Anything, mock.Anything).
		Return(errors.New("fake db error"))

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func testCorrectUserCreation(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(false, nil)
	mockUsersDao.On("CreateUser", "test@test.com", mock.Anything, mock.Anything).
		Return(nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func Test_CreateNewUser(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_createNewUser)

	t.Run("malformedBody", testCreateNewUserMalformedBody)

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
