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
	"os"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-template-repository/catalog"
	"github.com/trustedanalytics/tap-template-repository/model"
)

type Context struct {
	Template catalog.TemplateApi
}

var logger = logger_wrapper.InitLogger("api")

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result := []model.Template{}
	templatesMetadata := c.Template.GetAvailableTemplates()
	for _, templateMetadata := range templatesMetadata {
		template, err := c.Template.GetRawTemplate(templateMetadata, catalog.TemplatesPath)
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

	templateMetadata := c.Template.GetTemplateMetadataById(templateId)
	if templateMetadata == nil {
		util.Respond404(rw, errors.New(fmt.Sprintf("Can't find template by id: %s", templateId)))
		return
	}

	template, err := c.Template.GetParsedTemplate(templateMetadata, catalog.TemplatesPath, prepareReplacements(req.URL.Query(), instanceId))
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, template, http.StatusOK)
}

func prepareReplacements(query url.Values, instanceId string) map[string]string {
	additionalReplacements := make(map[string]string)
	additionalReplacements[model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_INSTANCE_ID)] = instanceId
	additionalReplacements[model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_DOMAIN_NAME)] = os.Getenv("DOMAIN")
	additionalReplacements[model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_ORG)] = "defaultOrg"
	additionalReplacements[model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_SPACE)] = "defaultSpace"

	for key, _ := range query {
		additionalReplacements["$"+key] = query.Get(key)
	}
	return additionalReplacements
}

func (c *Context) CreateCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := model.Template{}

	err := util.ReadJson(req, &reqTemplate)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = validateTemplateId(reqTemplate.Id)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	if c.Template.GetTemplateMetadataById(reqTemplate.Id) != nil {
		logger.Warning(fmt.Sprintf("Template with Id: %s already exists!", reqTemplate.Id))
		util.WriteJson(rw, "", http.StatusConflict)
		return
	}

	err = c.Template.AddAndRegisterCustomTemplate(reqTemplate)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusCreated)
}

func (c *Context) GetCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	templateID := req.PathParams["templateId"]
	err := validateTemplateId(templateID)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	templateMetadata := c.Template.GetTemplateMetadataById(templateID)
	if templateMetadata == nil {
		util.Respond404(rw, errors.New("Template doesn't exist!"))
		return
	}

	template, err := c.Template.GetRawTemplate(templateMetadata, catalog.TemplatesPath)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, template, http.StatusOK)
}

func (c *Context) DeleteCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	templateID := req.PathParams["templateId"]
	err := validateTemplateId(templateID)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.Template.RemoveAndUnregisterCustomTemplate(templateID)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusNoContent)
}
