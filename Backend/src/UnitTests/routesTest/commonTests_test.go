package routesTest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testParams struct {
	Error        string
	ResponseBody *bytes.Buffer
}

func httpTest(t *testing.T, body *bytes.Buffer, reqType, url string, responseCode int, params testParams) {
	//enable shorthand of omitting a body by passing nil
	if body == nil {
		body = new(bytes.Buffer)
	}

	//set up HTTP request
	req, _ := http.NewRequest(reqType, url, body)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	//get response from server and check that it's expected
	router.ServeHTTP(res, req)

	if res.Code != responseCode {
		var errorJson errorResponse
		err := json.Unmarshal(res.Body.Bytes(), &errorJson)
		var error string
		if err != nil {
			error = "unknown"
		} else {
			error = errorJson.Error
		}
		t.Errorf("Response code should be %v, was: %v. Error was: %v", responseCode, res.Code, error)
	}

	//check for errors if provided
	if params.Error != "" {
		var errorJson errorResponse
		err := json.Unmarshal(res.Body.Bytes(), &errorJson)
		if err != nil {
			println(err.Error())
			t.Error("Response body was not a json of an error form")
		}
		if errorJson.Error != params.Error {
			t.Errorf("Error should be %s, got %s", params.Error, errorJson.Error)
		}
	}

	//check response body if provided
	if params.ResponseBody != nil {
		if bytes.Compare(params.ResponseBody.Bytes(), res.Body.Bytes()) != 0 {
			t.Error("Expected JSON body and actual JSON body did not match")
		}
	}
}

func responseCodeTest(t *testing.T, body *bytes.Buffer, responseCode int, reqType, url string) *httptest.ResponseRecorder {
	//set up HTTP request
	req, _ := http.NewRequest(reqType, url, body)

	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	//get response from server and check that it's expected
	router.ServeHTTP(res, req)

	if res.Code != responseCode {
		var errorJson errorResponse
		err := json.Unmarshal(res.Body.Bytes(), &errorJson)
		var error string
		if err != nil {
			error = "unknown"
		} else {
			error = errorJson.Error
		}
		t.Errorf("Response code should be %v, was: %v. Error was: %v", responseCode, res.Code, error)
	}
	return res
}

func responseCodeAndErrorJsonTest(t *testing.T, body *bytes.Buffer, error, reqType, url string, responseCode int) {
	res := responseCodeTest(t, body, responseCode, reqType, url)

	//check we got the expected error
	var errorJson errorResponse
	err := json.Unmarshal(res.Body.Bytes(), &errorJson)
	if err != nil {
		t.Error("Response body was not a json of an error form")
	}
	if errorJson.Error != error {
		t.Errorf("Error should be %s, got %s", error, errorJson.Error)
	}
}

func testResponseAndCode(t *testing.T, reqBody, resBody *bytes.Buffer, reqType, url string, responseCode int) {
	res := responseCodeTest(t, reqBody, responseCode, reqType, url)
	if bytes.Compare(resBody.Bytes(), res.Body.Bytes()) != 0 {
		t.Error("Expected JSON body and actual JSON body did not match")
	}
}
