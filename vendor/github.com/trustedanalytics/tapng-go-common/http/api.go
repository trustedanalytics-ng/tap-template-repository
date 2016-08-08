package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ApiConnector struct {
	BasicAuth *BasicAuth
	Client    *http.Client
	Url       string
}

type CallFunc func(url string, requestBody string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error)

func GetModel(apiConnector ApiConnector, expectedStatus int, result interface{}) error {
	return callModelTemplateWithBody(RestGETWithBody, apiConnector, "", expectedStatus, result)
}

func PatchModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) error {
	return callModelTemplateWithBody(RestPATCH, apiConnector, requestBody, expectedStatus, result)

}

func PostModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) error {
	return callModelTemplateWithBody(RestPOST, apiConnector, requestBody, expectedStatus, result)
}

func PutModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) error {
	return callModelTemplateWithBody(RestPUT, apiConnector, requestBody, expectedStatus, result)
}

func DeleteModel(apiConnector ApiConnector, expectedStatus int) error {
	return callModelTemplateWithBody(RestDELETE, apiConnector, "", expectedStatus, "")
}

func callModelTemplateWithBody(callFunc CallFunc, apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) error {

	requestBodyByte := []byte{}
	var err error

	if requestBody != "" {
		requestBodyByte, err = json.Marshal(requestBody)
		if err != nil {
			return err
		}
	}

	status, body, err := callFunc(apiConnector.Url, string(requestBodyByte), apiConnector.BasicAuth, apiConnector.Client)
	if err != nil {
		return err
	}

	if result != "" {
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
	}

	if status != expectedStatus {
		return getWrongStatusError(status, expectedStatus, string(body))
	}
	return nil
}


func getWrongStatusError(status, expectedStatus int, body string) error {
	return errors.New(fmt.Sprintf("Bad response status: %d, expected status was: % d. Resposne body: %s", status, expectedStatus, body))
}
