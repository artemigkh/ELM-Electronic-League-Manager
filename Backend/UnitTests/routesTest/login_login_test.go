package routesTest

import (
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"errors"
)

//set up mocks
type mockUsersDAOLogin struct {
	t *testing.T
	e error
	emailInUse bool
	id int
	salt string
	storedHash string
}
func (u *mockUsersDAOLogin) CreateUser(email, salt, hash string) error {
	return nil
}
func (u *mockUsersDAOLogin) IsEmailInUse(email string) (bool, error) {
	return u.emailInUse, nil
}
func (u *mockUsersDAOLogin) GetAuthenticationInformation(email string) (int, string, string, error) {
	return u.id, u.salt, u.storedHash, u.e
}

type mockSessionsLoginError struct {
	t *testing.T
}
func (s *mockSessionsLoginError) AuthenticateAndGetUserID(ctx *gin.Context) (int, error) {
	return -1, nil
}
func (s *mockSessionsLoginError) LogIn(ctx *gin.Context, userID int) error {
	return errors.New("fake cookie error")
}

type mockSessionsLoginCorrect struct {
	t *testing.T
	LoggedIn bool
}
func (s *mockSessionsLoginCorrect) AuthenticateAndGetUserID(ctx *gin.Context) (int, error) {
	return 0, nil
}
func (s *mockSessionsLoginCorrect) LogIn(ctx *gin.Context, userID int) error {
	if userID != 517 {
		s.t.Errorf("login Id should be returned value (517), was %v", userID)
	}
	s.LoggedIn = true
	return nil
}

func testCreateNewUserNotEmailInUse(t *testing.T) {
	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: false,
		id: 1,
		salt: "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
	}
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
		"invalidLogin", "POST", 400)
}

func testBadPassword(t *testing.T) {
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "123456789lo7s"),
		"invalidLogin", "POST", 400)
}

func testBadHash(t *testing.T) {
	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: false,
		id: 1,
		salt: "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "005579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
	}
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "12345678"),
"invalidLogin", "POST", 400)
}

func testBadSalt(t *testing.T) {
	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: true,
		id: 1,
		salt: "0569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
	}
	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "12345678"),
		"invalidLogin", "POST", 400)
}

func testSessionError(t *testing.T) {
	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: true,
		id: 1,
		salt: "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
	}
	routes.ElmSessions = &mockSessionsLoginError{
		t: t,
	}
	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "12345678"), 500, "POST")
}

func testCorrectLogin(t *testing.T) {
	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: true,
		id: 517,
		salt: "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
	}
	 mockSessions := &mockSessionsLoginCorrect{
		t: t,
		LoggedIn: false,
	}
	routes.ElmSessions = mockSessions
	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "12345678"), 200, "POST")
	if mockSessions.LoggedIn == false {
		t.Error("Session login function was not called")
	}
}

func Test_Login(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_login)

	routes.UsersDAO = &mockUsersDAOLogin{
		t: t,
		e: nil,
		emailInUse: true,
		id: 1,
		salt: "f569dcbc75aa0d39462c00db0cdec7c5fae4e19bb3210838e8de0c843578a424",
		storedHash: "f05579a607c6847e35b2555453e2335b759b013fc313bf384f6b35810929b04bda7b448e6a9c784fea1015b06f74fff3784347d7e48df4af9e125d63499d3097",
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

	t.Run("loginEmailNotInUse", testCreateNewUserNotEmailInUse)
	t.Run("loginBadPassword", testBadPassword)
	t.Run("loginBadHash", testBadHash)
	t.Run("loginBadSalt", testBadSalt)
	t.Run("loginSessionError", testSessionError)
	t.Run("loginCorrect", testCorrectLogin)
}
