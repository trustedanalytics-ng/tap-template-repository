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

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-template-repository/model"
)

type TemplateRepository interface {
	GenerateParsedTemplate(templateId, uuid string) (model.Template, error)
}

type TemplateRepositoryConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func NewTemplateRepositoryBasicAuth(address, username, password string) (*TemplateRepositoryConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithBasicAuth()
	if err != nil {
		return nil, err
	}
	return &TemplateRepositoryConnector{address, username, password, client}, nil
}

func NewTemplateRepositoryCa(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TemplateRepositoryConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TemplateRepositoryConnector{address, username, password, client}, nil
}

func (t *TemplateRepositoryConnector) GenerateParsedTemplate(templateId, uuid string) (model.Template, error) {
	template := model.Template{}

	url := fmt.Sprintf("%s/api/v1/parsed_template/%s?serviceId=%s", t.Address, templateId, uuid)
	status, body, err := brokerHttp.RestGET(url, &brokerHttp.BasicAuth{t.Username, t.Password}, t.Client)
	err = json.Unmarshal(body, &template)
	if err != nil {
		return template, err
	}
	if status != http.StatusOK {
		return template, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return template, nil
}
