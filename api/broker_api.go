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

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tapng-go-common/logger"
	"github.com/trustedanalytics/tapng-go-common/util"
	"github.com/trustedanalytics/tapng-template-repository/catalog"
	"github.com/trustedanalytics/tapng-template-repository/model"
)

type Context struct {
	Template catalog.TemplateApi
}

var logger = logger_wrapper.InitLogger("api")

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result := []model.Template{}
	templatesMetadata := c.Template.GetAvailableTemplates()
	for _, templateMetadata := range templatesMetadata {
		template, err := c.Template.GetRawTemplate(templateMetadata, catalog.CatalogPath)
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
	uuid := req.URL.Query().Get("serviceId")
	if templateId == "" || uuid == "" {
		util.Respond500(rw, errors.New("templateId and uuid can't be empty!"))
		return
	}

	templateMetadata := c.Template.GetTemplateMetadataById(templateId)
	if templateMetadata == nil {
		util.Respond500(rw, errors.New(fmt.Sprintf("Can't find template by id: %s", templateId)))
		return
	}

	template, err := c.Template.GetParsedTemplate(templateMetadata, catalog.CatalogPath, uuid, "defaultOrg", "defaultSpace")
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, template, http.StatusOK)
}

func (c *Context) CreateCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := model.Template{}

	err := util.ReadJson(req, &reqTemplate)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	if reqTemplate.Id == "" {
		util.Respond500(rw, errors.New("Teplate Id can not be empty!"))
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
	templateId := req.PathParams["templateId"]
	if templateId == "" {
		util.Respond500(rw, errors.New("templateId can not be empty!"))
		return
	}

	templateMetadata := c.Template.GetTemplateMetadataById(templateId)
	if templateMetadata == nil {
		util.Respond500(rw, errors.New("Template not exist!"))
		return
	}

	template, err := c.Template.GetRawTemplate(templateMetadata, catalog.CatalogPath)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, template, http.StatusOK)
}

func (c *Context) DeleteCustomTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	if templateId == "" {
		util.Respond500(rw, errors.New("templateId can not be empty!"))
		return
	}

	err := c.Template.RemoveAndUnregisterCustomTemplate(templateId)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusNoContent)
}
