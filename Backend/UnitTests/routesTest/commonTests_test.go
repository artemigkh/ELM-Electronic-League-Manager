package routesTest

import (
	"net/http"
	"bytes"
	"net/http/httptest"
	"encoding/json"
	"testing"
)

func responseCodeTest(t *testing.T, body *bytes.Buffer, responseCode int, reqType string) *httptest.ResponseRecorder{
	//set up HTTP request
	req, _ := http.NewRequest(reqType, "/", body)


	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	//get response from server and check that it's expected
	router.ServeHTTP(res, req)

	if res.Code != responseCode {
		t.Errorf("Response code should be %v, was: %v", responseCode, res.Code)
	}
	return res
}

func responseCodeAndErrorJsonTest(t *testing.T, body *bytes.Buffer, error, reqType string, responseCode int) {
	res := responseCodeTest(t, body, responseCode, reqType)

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

func testResponseAndCode(t *testing.T, reqBody, resBody *bytes.Buffer, reqType string, responseCode int) {
	res := responseCodeTest(t, reqBody, responseCode, reqType)
	if bytes.Compare(resBody.Bytes(), res.Body.Bytes()) != 0 {
		t.Error("Expected JSON body and actual JSON body did not match")
	}
}