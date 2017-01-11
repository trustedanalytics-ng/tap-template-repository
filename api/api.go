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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gocraft/web"

	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
	util "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-template-repository/catalog"
	"github.com/trustedanalytics/tap-template-repository/model"
)

type Context struct {
	TemplateApi catalog.TemplateApi
}

var logger, _ = commonLogger.InitLogger("api")

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result := []model.RawTemplate{}
	for _, templatePath := range c.TemplateApi.GetAvailableTemplates() {
		template, err := c.TemplateApi.GetRawTemplate(templatePath)
		if err != nil {
			util.Respond500(rw, err)
			return
		}
		result = append(result, template)
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GenerateParsedTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	instanceId := req.URL.Query().Get("instanceId")
	planName := req.URL.Query().Get("planName")

	err := validateTemplateId(templateId)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = validateUuid(instanceId)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	templatePath := c.TemplateApi.GetTemplatePath(templateId)
	if templatePath == "" {
		util.Respond404(rw, errors.New(fmt.Sprintf("can't find template by id: %s", templateId)))
		return
	}

	rawTemplate, err := c.TemplateApi.GetRawTemplate(templatePath)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	template, err := c.TemplateApi.GetParsedTemplate(rawTemplate, prepareReplacements(req.URL.Query(), instanceId), planName)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, *template, http.StatusOK)
}

func prepareReplacements(query url.Values, instanceId string) map[string]string {
	replacements := make(map[string]string)
	replacements[model.GetPlaceholderWithDollarPrefix(model.PlaceholderInstanceID)] = instanceId

	for key, _ := range query {
		replacements["$"+key] = query.Get(key)
	}

	return model.GetMapWithDefaultReplacementsIfKeyNotExists(replacements)
}

func (c *Context) CreateCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	rawTemplate := model.RawTemplate{}

	err := util.ReadJson(req, &rawTemplate)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	//Validate if template can be properely parsed, use fake-test-instance-id as this is required during regular creation
	template, err := c.TemplateApi.GetParsedTemplate(rawTemplate, prepareReplacements(req.URL.Query(), "fake-test-instance-id"), model.EMPTY_PLAN_NAME)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	templateId := template.Id
	if err = validateTemplateId(templateId); err != nil {
		util.Respond400(rw, err)
		return
	}

	if _, templateExist := c.TemplateApi.GetAvailableTemplates()[templateId]; templateExist {
		util.Respond409(rw, errors.New(fmt.Sprintf("Template with Id: %s already exists!", templateId)))
		return
	}

	if err := c.TemplateApi.AddCustomTemplate(rawTemplate, templateId); err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusCreated)
}

func (c *Context) GetRawTemplate(rw web.ResponseWriter, req *web.Request) {
	templateID := req.PathParams["templateId"]
	err := validateTemplateId(templateID)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	templatePath := c.TemplateApi.GetTemplatePath(templateID)
	if templatePath == "" {
		util.Respond404(rw, errors.New("template doesn't exist!"))
		return
	}

	rawTemplate, err := c.TemplateApi.GetRawTemplate(templatePath)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, rawTemplate, http.StatusOK)
}

func (c *Context) DeleteCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	templateID := req.PathParams["templateId"]
	err := validateTemplateId(templateID)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	status, err := c.TemplateApi.RemoveAndUnregisterCustomTemplate(templateID)
	if err != nil {
		util.GenericRespond(status, rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusNoContent)
}
