package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ApiConnector struct {
	BasicAuth *BasicAuth
	OAuth2    *OAuth2
	Client    *http.Client
	Url       string
}

type CallFunc func(url string, requestBody, authHeader string, client *http.Client) (int, []byte, error)

func GetModel(apiConnector ApiConnector, expectedStatus int, result interface{}) (int, error) {
	return callModelTemplateWithBody(RestGETWithBody, apiConnector, "", expectedStatus, result)
}

func PatchModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) (int, error) {
	return callModelTemplateWithBody(RestPATCH, apiConnector, requestBody, expectedStatus, result)

}

func PostModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) (int, error) {
	return callModelTemplateWithBody(RestPOST, apiConnector, requestBody, expectedStatus, result)
}

func PutModel(apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) (int, error) {
	return callModelTemplateWithBody(RestPUT, apiConnector, requestBody, expectedStatus, result)
}

func DeleteModel(apiConnector ApiConnector, expectedStatus int) (int, error) {
	return callModelTemplateWithBody(RestDELETE, apiConnector, "", expectedStatus, "")
}

func callModelTemplateWithBody(callFunc CallFunc, apiConnector ApiConnector, requestBody interface{}, expectedStatus int, result interface{}) (status int, err error) {
	requestBodyByte := []byte{}

	if requestBody != "" {
		requestBodyByte, err = json.Marshal(requestBody)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	authHeader := ""
	if apiConnector.OAuth2 != nil {
		authHeader = GetOAuth2Header(apiConnector.OAuth2)
	}
	if apiConnector.BasicAuth != nil {
		authHeader = GetBasicAuthHeader(apiConnector.BasicAuth)
	}

	status, body, err := callFunc(apiConnector.Url, string(requestBodyByte), authHeader, apiConnector.Client)
	if err != nil {
		return status, err
	}

	if result != "" {
		err = json.Unmarshal(body, result)
		if err != nil {
			return status, err
		}
	}

	if status != expectedStatus {
		return status, getWrongStatusError(status, expectedStatus, string(body))
	}
	return status, nil
}

func getWrongStatusError(status, expectedStatus int, body string) error {
	return errors.New(fmt.Sprintf("Bad response status: %d, expected status was: % d. Response body: %s", status, expectedStatus, body))
}
