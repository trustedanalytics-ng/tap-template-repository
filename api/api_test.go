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
	"testing"

	"github.com/gocraft/web"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-template-repository/catalog"
	"github.com/trustedanalytics/tap-template-repository/model"
)

const (
	templateId   string = "test"
	templatePath string = "path"
	instanceId   string = "a5740d8a-9f4b-4711-a1a0-eae62db54474"
	planName     string = "samplePlan"
)

var rawTemplate = model.RawTemplate{model.RAW_TEMPLATE_ID_FIELD: templateId}
var template = model.Template{
	Id: templateId,
}

func prepareMocksAndRouter(t *testing.T) (router *web.Router, c Context, templateMock *catalog.MockTemplateApi) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	templateMock = catalog.NewMockTemplateApi(mockCtrl)
	c = Context{templateMock}
	router = web.New(c)
	return
}

func TestTemplates(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Get("/api/v1/templates", context.Templates)

	convey.Convey("Test Templates", t, func() {
		convey.Convey("No templates available", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetAvailableTemplates().Return(make(map[string]string)),
			)
			response := util.SendRequest("GET", "/api/v1/templates", nil, router)
			util.AssertResponse(response, "[]", 200)
		})
		convey.Convey("Template retrieval failed", func() {
			availableTemplates := make(map[string]string)
			availableTemplates[templateId] = templatePath
			gomock.InOrder(
				templateMock.EXPECT().GetAvailableTemplates().Return(availableTemplates),
				templateMock.EXPECT().GetRawTemplate(availableTemplates[templateId]).Return(rawTemplate, errors.New("failed")),
			)
			response := util.SendRequest("GET", "/api/v1/templates", nil, router)
			util.AssertResponse(response, "failed", 500)
		})
		convey.Convey("All avialable templates successfully retrieved", func() {
			availableTemplates := make(map[string]string)
			availableTemplates[templateId] = templatePath
			gomock.InOrder(
				templateMock.EXPECT().GetAvailableTemplates().Return(availableTemplates),
				templateMock.EXPECT().GetRawTemplate(availableTemplates[templateId]).Return(rawTemplate, nil),
			)
			response := util.SendRequest("GET", "/api/v1/templates", nil, router)
			util.AssertResponse(response, templateId, 200)
		})
	})
}

func TestGenerateParsedTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Get("/api/v1/parsed_template/:templateId/", context.GenerateParsedTemplate)

	convey.Convey("Test Generate Parsed Template", t, func() {
		convey.Convey("No templateId provided", func() {
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/parsed_template//?instanceId=%s", instanceId), nil, router)
			util.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("No instanceId provided", func() {
			response := util.SendRequest("GET", "/api/v1/parsed_template/1", nil, router)
			util.AssertResponse(response, "uuid can't be empty!", 400)
		})
		convey.Convey("Template with templateId not found", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(""),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/parsed_template/%s?instanceId=%s", templateId, instanceId), nil, router)
			util.AssertResponse(response, "can't find template by id: "+templateId, 404)
		})
		convey.Convey("Getting parsed component failed", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(templatePath),
				templateMock.EXPECT().GetRawTemplate(templatePath).Return(rawTemplate, nil),
				templateMock.EXPECT().GetParsedTemplate(rawTemplate, gomock.Any(), planName).Return(&model.Template{}, errors.New("failed")),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/parsed_template/%s?instanceId=%s&planName=%s", templateId, instanceId, planName), nil, router)
			util.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Existing templateId provided", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(templatePath),
				templateMock.EXPECT().GetRawTemplate(templatePath).Return(rawTemplate, nil),
				templateMock.EXPECT().GetParsedTemplate(rawTemplate, gomock.Any(), planName).Return(&model.Template{Id: templateId}, nil),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/parsed_template/%s?instanceId=%s&planName=%s", templateId, instanceId, planName), nil, router)
			var template model.Template
			json.Unmarshal(response.Body.Bytes(), &template)
			convey.So(template.Id, convey.ShouldEqual, templateId)
			convey.So(response.Code, convey.ShouldEqual, 200)
		})

	})
}

func TestCreateCustomTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Post("/api/v1/templates", context.CreateCustomTemplate)

	convey.Convey("Test Create Custom Template", t, func() {
		convey.Convey("Not parsable body", func() {
			response := util.SendRequest("POST", "/api/v1/templates", []byte("not_parsable"), router)
			util.AssertResponse(response, "invalid character", 500)
		})
		convey.Convey("Template without id", func() {
			reqBody, _ := json.Marshal(model.RawTemplate{})
			gomock.InOrder(
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), model.EMPTY_PLAN_NAME).Return(&model.Template{}, nil),
			)
			response := util.SendRequest("POST", "/api/v1/templates", reqBody, router)
			util.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Template with id exists", func() {
			body_bytes, _ := json.Marshal(template)
			gomock.InOrder(
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), model.EMPTY_PLAN_NAME).Return(&template, nil),
				templateMock.EXPECT().GetAvailableTemplates().Return(map[string]string{templateId: "path to existing template"}),
			)
			response := util.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			convey.So(response.Code, convey.ShouldEqual, 409)
		})
		convey.Convey("Adding template fails", func() {
			body_bytes, _ := json.Marshal(template)
			gomock.InOrder(
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), model.EMPTY_PLAN_NAME).Return(&template, nil),
				templateMock.EXPECT().GetAvailableTemplates().Return(make(map[string]string)),
				templateMock.EXPECT().AddCustomTemplate(gomock.Any(), templateId).Return(errors.New("failed")),
			)
			response := util.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			util.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully added", func() {
			body_bytes, _ := json.Marshal(template)
			gomock.InOrder(
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), model.EMPTY_PLAN_NAME).Return(&template, nil),
				templateMock.EXPECT().GetAvailableTemplates().Return(make(map[string]string)),
				templateMock.EXPECT().AddCustomTemplate(gomock.Any(), templateId).Return(nil),
			)
			response := util.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			convey.So(response.Code, convey.ShouldEqual, 201)
		})
	})
}

func TestGetRawTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Get("/api/v1/templates/:templateId", context.GetRawTemplate)

	convey.Convey("Test Get Raw Template", t, func() {
		convey.Convey("No template id provided", func() {
			response := util.SendRequest("GET", "/api/v1/templates//", nil, router)
			util.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Template does not exist", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(""),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/templates/%s", templateId), nil, router)
			util.AssertResponse(response, "template doesn't exist!", 404)
		})
		convey.Convey("Error gettting template", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(templatePath),
				templateMock.EXPECT().GetRawTemplate(templatePath).Return(rawTemplate, errors.New("failed")),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/templates/%s", templateId), nil, router)
			util.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully retrieved template", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplatePath(templateId).Return(templatePath),
				templateMock.EXPECT().GetRawTemplate(templatePath).Return(rawTemplate, nil),
			)
			response := util.SendRequest("GET", fmt.Sprintf("/api/v1/templates/%s", templateId), nil, router)
			util.AssertResponse(response, templateId, 200)
		})
	})
}

func TestDeleteCustomTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Delete("/api/v1/templates/:templateId", context.DeleteCustomTemplate)

	convey.Convey("Test Delete Custom Template", t, func() {
		convey.Convey("No template id provided", func() {
			response := util.SendRequest("DELETE", "/api/v1/templates//", nil, router)
			util.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Deletion failed", func() {
			gomock.InOrder(
				templateMock.EXPECT().RemoveAndUnregisterCustomTemplate(templateId).Return(http.StatusInternalServerError, errors.New("failed")),
			)
			response := util.SendRequest("DELETE", fmt.Sprintf("/api/v1/templates/%s", templateId), nil, router)
			util.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully removed", func() {
			gomock.InOrder(
				templateMock.EXPECT().RemoveAndUnregisterCustomTemplate(templateId).Return(http.StatusNoContent, nil),
			)
			response := util.SendRequest("DELETE", fmt.Sprintf("/api/v1/templates/%s", templateId), nil, router)
			convey.So(response.Code, convey.ShouldEqual, http.StatusNoContent)
		})
	})
}
