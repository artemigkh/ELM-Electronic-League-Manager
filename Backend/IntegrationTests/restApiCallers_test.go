package IntegrationTests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"io/ioutil"
)

func makeApiCall(t *testing.T, bodyJson map[string]interface{}, reqType, url string, responseCode int) *bytes.Buffer {
	var body *bytes.Buffer
	//enable shorthand of omitting a body by passing nil
	if bodyJson == nil {
		body = new(bytes.Buffer)
	} else {
		bodyBytes, err := json.Marshal(bodyJson)
		if err != nil {
			t.Error("Could not Marshall dict into json")
		}
		body = bytes.NewBuffer(bodyBytes)
	}

		//set up HTTP request
	req, _ := http.NewRequest(reqType, url, body)

	req.Header.Set("Content-Type", "application/json")

	//get response from server and check that it's expected
	res, err := client.Do(req)
	if err != nil {
		t.Error("Request to server failed")
	}

	var bodyBytes *bytes.Buffer
	if res.StatusCode != responseCode {
		var errorJson errorResponse
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error("Could read response body")
		}

		err = json.Unmarshal(bodyBytes, &errorJson)
		var errorMsg string
		if err != nil {
			errorMsg = "unknown"
		} else {
			errorMsg = errorJson.Error
		}
		t.Errorf("Response code should be %v, was: %v. Error was: %v", responseCode, res.StatusCode, errorMsg)
	}

	return bodyBytes
}

func makeApiCallAndGetId(t *testing.T, bodyJson map[string]interface{},
	reqType, url string, responseCode int) int {
	body := makeApiCall(t, bodyJson, reqType, url, responseCode)

	var id idResponse
	err := json.Unmarshal(body.Bytes(), id)
	if err != nil {
		t.Error("Could not unmarshall response into an id")
	}

	return id.Id
}