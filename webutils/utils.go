/**
 * Copyright (c) 2015 Intel Corporation
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

package webutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"


	"github.com/gocraft/web"
)

type TapJWTToken struct {
	Jti       string   `json:"jti"`
	Sub       string   `json:"sub"`
	Scope     []string `json:"scope"`
	ClientId  string   `json:"client_id"`
	Cid       string   `json:"cid"`
	Azp       string   `json:"azp"`
	GrantType string   `json:"grant_type"`
	UserId    string   `json:"user_id"`
	Username  string   `json:"user_name"`
	Email     string   `json:"email"`
	RevSig    string   `json:"rev_sig"`
	Iat       int64    `json:"iat"`
	Exp       int64    `json:"exp"`
	Iss       string   `json:"iss"`
	Zid       string   `json:"zid"`
	Aud       []string `json:"aud"`
}

func ReadJson(req *web.Request, retstruct interface{}) error {
	var err error
	body, err := ioutil.ReadAll(req.Body)
	log.Println("Request body:", string(body))
	if err != nil {
		log.Println("Error reading request body:", err)
		return err
	}
	b := []byte(body)
	err = json.Unmarshal(b, &retstruct)
	if err != nil {
		log.Println("Error parsing request body json:", err)
		return err
	}
	log.Println("Request JSON parsed as: ", retstruct)
	return nil
}

func WriteJson(rw web.ResponseWriter, response interface{}, status_code int) error {
	b, err := json.Marshal(&response)
	if err != nil {
		log.Println("Error marshalling response:", err)
		return err
	}
	rw.Header().Set("Content-Type", "application/json")
	log.Println("Responding with status", status_code, " and JSON:", string(b))
	rw.WriteHeader(status_code)
	fmt.Fprintf(rw, "%s", string(b))
	return nil
}


