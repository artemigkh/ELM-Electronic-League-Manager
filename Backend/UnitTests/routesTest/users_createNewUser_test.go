package routesTest

import (
	"esports-league-manager/Backend/Server/routes"
	"github.com/gin-gonic/gin"
	"testing"
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

//func testCreateNewUserEmailInUse(t *testing.T) {
//	responseCodeAndErrorJsonTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
//		"emailInUse", "POST", "/", 400)
//}

//func testCreateNewUserDatabaseError(t *testing.T) {
//	routes.UsersDAO = &mockUsersDAO{
//		t: t,
//		e: errors.New("fake database error"),
//	}
//	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
//		500, "POST", "/")
//}
//
//func testCorrectUserCreation(t *testing.T) {
//	mockDAO := &mockUsersDAOCreateUser{
//		t: t,
//		UserCreated: false,
//	}
//	routes.UsersDAO = mockDAO
//
//	responseCodeTest(t, createUserCreateRequestBody("test@test.com", "abcd1234"),
//		200, "POST", "/")
//	if !mockDAO.UserCreated {
//		t.Error("User creation DAO function was not called")
//	}
//}

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

	//t.Run("emailInUse", testCreateNewUserEmailInUse)
	//t.Run("databaseError", testCreateNewUserDatabaseError)
	//t.Run("correctUserCreation", testCorrectUserCreation)
}
