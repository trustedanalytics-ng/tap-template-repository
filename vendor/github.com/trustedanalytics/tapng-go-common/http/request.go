/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type BasicAuth struct {
	User     string
	Password string
}

func RestGETWithBody(url string, body string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("GET", url, body, "application/json", basicAuth, client)
}

func RestGET(url string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("GET", url, "", "application/json", basicAuth, client)
}

func RestPUT(url, body string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("PUT", url, body, "application/json", basicAuth, client)
}

func RestPOST(url, body string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("POST", url, body, "application/json", basicAuth, client)
}

func RestDELETE(url, body string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("DELETE", url, body, "", basicAuth, client)
}

func RestPATCH(url, body string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	return makeRequest("PATCH", url, body, "application/json-patch+json", basicAuth, client)
}

func makeRequest(reqType, url, body, contentType string, basicAuth *BasicAuth, client *http.Client) (int, []byte, error) {
	logger.Info("Doing:  ", reqType, url)

	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(reqType, url, bytes.NewBuffer([]byte(body)))
	} else {
		req, _ = http.NewRequest(reqType, url, nil)
	}
	AddBasicAuth(req, basicAuth)
	SetContentType(req, contentType)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ERROR: Make http request "+reqType, err)
		return -1, nil, err
	}
	ret_code := resp.StatusCode
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("ERROR: Make http request "+reqType, err)
		return -1, nil, err
	}

	if resp.Header.Get("Content-Type") == "application/octet-stream" {
		logger.Info("CODE:", ret_code, "BODY: [ Binary Data ]", resp.ContentLength)
	} else {
		logger.Info("CODE:", ret_code, "BODY:", string(data))
	}

	return ret_code, data, nil
}

func AddBasicAuth(req *http.Request, basicAuth *BasicAuth) {
	if basicAuth != nil {
		req.SetBasicAuth(basicAuth.User, basicAuth.Password)
	}
}

func SetContentType(req *http.Request, contentType string) {
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
}

func DownloadBinary(url string, basicAuth *BasicAuth, client *http.Client, dest io.Writer) (int64, error) {
	return binaryStreamRequest(url, basicAuth, client, dest)
}

func binaryStreamRequest(url string, basicAuth *BasicAuth, client *http.Client, dest io.Writer) (int64, error) {
	logger.Info("Doing:  ", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("ERROR: Make http request ", err)
		return -1, err
	}

	if basicAuth != nil {
		req.SetBasicAuth(basicAuth.User, basicAuth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ERROR: Make http request ", err)
		return -1, err
	}

	defer resp.Body.Close()
	_, err = io.CopyN(dest, resp.Body, resp.ContentLength)
	if err != nil {
		logger.Error("ERROR: Make http request ", err)
		return -1, err
	}

	logger.Info("CODE:", resp.StatusCode, "BODY: [ Binary Data ] Size:", resp.ContentLength)
	return resp.ContentLength, nil
}
