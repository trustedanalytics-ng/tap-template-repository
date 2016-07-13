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
	"testing"

	"github.com/gocraft/web"

	"encoding/json"
	"github.com/smartystreets/goconvey/convey"
	"github.com/trustedanalytics/tapng-template-repository/model"
	TestUtils "github.com/trustedanalytics/tapng-template-repository/test"
	catalog_mock "github.com/trustedanalytics/tapng-template-repository/test/catalog"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router, c Context) {
	c = Context{}
	router = web.New(c)
	return router, c
}

func TestGenerateParsedTemplate(t *testing.T) {
	var (
		originalGetRawTemplate                    = getRawTemplate
		originalGetAvailableTemplates             = getAvailableTemplates
		originalGetTemplateMetadataById           = getTemplateMetadataById
		originalGetParsedTemplate                 = getParsedTemplate
		originalAddAndRegisterCustomTemplate      = addAndRegisterCustomTemplate
		originalRemoveAndUnregisterCustomTemplate = removeAndUnregisterCustomTemplate
	)

	defer func() {
		getRawTemplate = originalGetRawTemplate
		getAvailableTemplates = originalGetAvailableTemplates
		getTemplateMetadataById = originalGetTemplateMetadataById
		getParsedTemplate = originalGetParsedTemplate
		addAndRegisterCustomTemplate = originalAddAndRegisterCustomTemplate
		removeAndUnregisterCustomTemplate = originalRemoveAndUnregisterCustomTemplate
	}()

	router, context := prepareMocksAndRouter(t)
	router.Get("/api/v1/parsed_template/:templateId/", context.GenerateParsedTemplate)

	convey.Convey("Test Generate Parsed Template", t, func() {
		convey.Convey("No templateId provided", func() {
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/?serviceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			TestUtils.AssertResponse(response, "Not Found", 404)
		})
		convey.Convey("No serviceId provided", func() {
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/1", nil, router)
			TestUtils.AssertResponse(response, "templateId and uuid can't be empty!", 500)
		})
		convey.Convey("Template with templateId not found", func() {
			getTemplateMetadataById = catalog_mock.GetTemplateMetadataByIdNilMock
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/templateId?serviceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			TestUtils.AssertResponse(response, "Can't find template by id: templateId", 500)
		})
		convey.Convey("Non existing templateId provided", func() {
			getTemplateMetadataById = catalog_mock.GetTemplateMetadataByIdTemplateMock
			getParsedTemplate = catalog_mock.GetParsedTemplate
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/non_existing?serviceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			var template model.Template
			json.Unmarshal(response.Body.Bytes(), &template)
			convey.So(template.Id, convey.ShouldEqual, "non_existing")
			convey.So(response.Code, convey.ShouldEqual, 200)
		})

	})
}
