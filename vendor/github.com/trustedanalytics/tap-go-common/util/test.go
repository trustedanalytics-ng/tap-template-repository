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

package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gocraft/web"
	"github.com/smartystreets/goconvey/convey"
)

func SendRequest(rType, path string, body []byte, r *web.Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(rType, path, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func PrepareAndValidateRequest(v interface{}, t *testing.T) []byte {
	byteBody, marshalError := json.Marshal(v)
	if marshalError != nil {
		t.Fatal("Marshal request error: ", marshalError)
	}
	return byteBody
}

func AssertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		convey.So(strings.TrimSpace(string(rr.Body.Bytes())), convey.ShouldContainSubstring, body)
	}
	convey.So(rr.Code, convey.ShouldEqual, code)
}
