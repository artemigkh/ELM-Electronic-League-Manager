package routesTest

import (
	"net/http"
	"bytes"
	"net/http/httptest"
	"encoding/json"
	"testing"
)

func responseCodeTest(t *testing.T, email, pass string, responseCode int) *httptest.ResponseRecorder{
	//set up HTTP request
	reqBody := userCreateRequest{
		Email:    email,
		Password: pass,
	}
	reqBodyB, _ := json.Marshal(&reqBody)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBodyB))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	//get response from server and check that it's expected
	router.ServeHTTP(res, req)
	if res.Code != responseCode {
		t.Errorf("Response code should be %v, was: %v", responseCode, res.Code)
	}

	return res
}

func jsonErrorTest(t *testing.T, email, pass, error string) {
	res := responseCodeTest(t, email, pass, 400)

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