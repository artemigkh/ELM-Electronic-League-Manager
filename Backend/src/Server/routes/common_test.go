package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *gin.Engine

const (
	testUserId   int = iota
	testLeagueId int = iota
	testTeamId   int = iota
	testPlayerId int = iota
	testGameId   int = iota
)

var testSessionErrMsg = "Session Test Error"
var testSessionErr = errors.New(testSessionErrMsg)

var testBindErrMsg = "malformedInput"
var testInvalidUrlMsg = "IdMustBeInteger"

var testAccessErrMsg = "Permissions Test Error"
var testPermissionsErr = errors.New(testAccessErrMsg)

var testValidatorErrMsg = "Validator Test Error"
var testValidatorErr = errors.New(testValidatorErrMsg)

var testDaoErrMsg = "Dao Test Error"
var testDaoErr = errors.New(testDaoErrMsg)

type errorResponse struct {
	Error string `json:"error"`
}

type userIdResponse struct {
	UserId int `json:"userId"`
}

type leagueIdResponse struct {
	LeagueId int `json:"leagueId"`
}

type teamIdResponse struct {
	TeamId int `json:"teamId"`
}

type playerIdResponse struct {
	PlayerId int `json:"playerId"`
}

type gameIdResponse struct {
	GameId int `json:"gameId"`
}

func httpBody(body interface{}) *bytes.Buffer {
	bodyBytes, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyBytes)
}

type httpTestRunner interface {
	RunHttpTest()
}

type httpTest struct {
	T            *testing.T
	RequestData  interface{}
	ResponseData interface{}
	Type         string
	Url          string
	ResponseCode int
	Error        string
}

func (args httpTest) RunHttpTest() {
	// Get request body as bytes
	var requestBody *bytes.Buffer
	if args.RequestData == nil {
		requestBody = new(bytes.Buffer)
	} else {
		requestBody = httpBody(args.RequestData)
	}

	//set up HTTP request
	req, _ := http.NewRequest(args.Type, args.Url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	//get response from server and check that it's expected
	router.ServeHTTP(res, req)

	if res.Code != args.ResponseCode {
		var errorJson errorResponse
		err := json.Unmarshal(res.Body.Bytes(), &errorJson)
		var error string
		if err != nil {
			error = "unknown"
		} else {
			error = errorJson.Error
		}
		args.T.Errorf("Response code should be %v, was: %v. Error was: %v", args.ResponseCode, res.Code, error)
	}

	//check for errors if provided
	if args.Error != "" {
		var errorJson errorResponse
		err := json.Unmarshal(res.Body.Bytes(), &errorJson)
		if err != nil {
			println(err.Error())
			args.T.Error("Response body was not a json of an error form")
		}
		if errorJson.Error != args.Error {
			args.T.Errorf("Error should be %s, got %s", args.Error, errorJson.Error)
		}
	}

	//check response body if provided
	if args.ResponseData != nil {
		responseBody := httpBody(args.ResponseData)
		if bytes.Compare(responseBody.Bytes(), res.Body.Bytes()) != 0 {
			args.T.Error(fmt.Sprintf("Expected JSON body and actual JSON body did not match. "+
				"Exepcted: %v. Actual: %v", string(responseBody.Bytes()), string(res.Body.Bytes())))
		}
	}
}
