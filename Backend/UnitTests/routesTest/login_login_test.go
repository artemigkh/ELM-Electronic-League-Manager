package routesTest

import (
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"bytes"
	"encoding/json"
	"esports-league-manager/mocks"
	"github.com/stretchr/testify/mock"
	"strings"
	"errors"
)

func createLoginRequestBody(email, pass string) *bytes.Buffer {
	reqBody := userCreateRequest{
		Email:    email,
		Password: pass,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

func testLoginMalformedBody(t *testing.T) {
	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})
}

func testLoginPasswordTooShort(t *testing.T, pass string) {
	httpTest(t, createLoginRequestBody("test@test.com", pass), "POST", "/",
		400, testParams{Error: "passwordTooShort"})
}

func testLoginMalformedEmail(t *testing.T, email string) {
	httpTest(t, createLoginRequestBody(email, "12345678"), "POST", "/",
		400, testParams{Error: "emailMalformed"})
}

func testCreateNewUserNotEmailInUse(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(false, nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 400, testParams{Error: "invalidLogin"})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

var salt = "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424"
var storedHash = "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097"

func testBadPassword(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)
	mockUsersDao.On("GetAuthenticationInformation", "test@test.com").
		Return(1, salt, storedHash, nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678o3294"),
		"POST", "/", 400, testParams{Error: "invalidLogin"})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func testBadHash(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)
	mockUsersDao.On("GetAuthenticationInformation", "test@test.com").
		Return(1, salt, strings.Replace(storedHash, "f", "0", 1), nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678o3294"),
		"POST", "/", 400, testParams{Error: "invalidLogin"})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func testBadSalt(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)
	mockUsersDao.On("GetAuthenticationInformation", "test@test.com").
		Return(1, strings.Replace(salt, "f", "0", 1), storedHash, nil)

	routes.UsersDAO = mockUsersDao

	httpTest(t, createLoginRequestBody("test@test.com", "12345678o3294"),
		"POST", "/", 400, testParams{Error: "invalidLogin"})

	mock.AssertExpectationsForObjects(t, mockUsersDao)
}

func testSessionError(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)
	mockUsersDao.On("GetAuthenticationInformation", "test@test.com").
		Return(1, salt, storedHash, nil)

	mockSession := new(mocks.SessionManager)
	mockSession.On("LogIn", mock.Anything, 1).Return(errors.New("fake session error"))

	routes.UsersDAO = mockUsersDao
	routes.ElmSessions = mockSession

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockUsersDao, mockSession)
}

func testCorrectLogin(t *testing.T) {
	mockUsersDao := new(mocks.UsersDAO)
	mockUsersDao.On("IsEmailInUse", "test@test.com").Return(true, nil)
	mockUsersDao.On("GetAuthenticationInformation", "test@test.com").
		Return(1, salt, storedHash, nil)

	mockSession := new(mocks.SessionManager)
	mockSession.On("LogIn", mock.Anything, 1).Return(nil)

	routes.UsersDAO = mockUsersDao
	routes.ElmSessions = mockSession

	httpTest(t, createLoginRequestBody("test@test.com", "12345678"),
		"POST", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockUsersDao, mockSession)
}

func Test_Login(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_login)

	t.Run("malformedBody", testLoginMalformedBody)

	t.Run("passwordTooShort=Length7A", func(t *testing.T) {
		testLoginPasswordTooShort(t, "1234567")
	})
	t.Run("passwordTooShort=Length=7B", func(t *testing.T) {
		testLoginPasswordTooShort(t, "123456 ")
	})
	t.Run("passwordTooShortLength=0", func(t *testing.T) {
		testLoginPasswordTooShort(t, "")
	})

	//these do not have to be comprehensive as they are handled by a third party lib
	t.Run("malformedEmail=InvalidCharacters", func(t *testing.T) {
		testLoginMalformedEmail(t, "ç$€§/az@gmail.com")
	})
	t.Run("malformedEmail=NoName", func(t *testing.T) {
		testLoginMalformedEmail(t, "@gmail.com")
	})

	t.Run("loginEmailNotInUse", testCreateNewUserNotEmailInUse)
	t.Run("loginBadPassword", testBadPassword)
	t.Run("loginBadHash", testBadHash)
	t.Run("loginBadSalt", testBadSalt)
	t.Run("loginSessionError", testSessionError)
	t.Run("loginCorrect", testCorrectLogin)
}
