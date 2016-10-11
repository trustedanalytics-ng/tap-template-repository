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
	"net/http"
	"testing"

	"github.com/gocraft/web"

	"encoding/json"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/trustedanalytics/tap-template-repository/catalog"
	"github.com/trustedanalytics/tap-template-repository/model"
	TestUtils "github.com/trustedanalytics/tap-template-repository/test"
)

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
				templateMock.EXPECT().GetAvailableTemplates().Return(make(map[string]*model.TemplateMetadata)),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates", nil, router)
			TestUtils.AssertResponse(response, "[]", 200)
		})
		convey.Convey("Template retrieval failed", func() {
			availableTemplates := make(map[string]*model.TemplateMetadata)
			availableTemplates["test"] = &model.TemplateMetadata{Id: "templateId"}
			gomock.InOrder(
				templateMock.EXPECT().GetAvailableTemplates().Return(availableTemplates),
				templateMock.EXPECT().GetRawTemplate(availableTemplates["test"], gomock.Any()).Return(model.Template{}, errors.New("failed")),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates", nil, router)
			TestUtils.AssertResponse(response, "failed", 500)
		})
		convey.Convey("All avialable templates successfully retrieved", func() {
			availableTemplates := make(map[string]*model.TemplateMetadata)
			availableTemplates["test"] = &model.TemplateMetadata{Id: "templateId"}
			gomock.InOrder(
				templateMock.EXPECT().GetAvailableTemplates().Return(availableTemplates),
				templateMock.EXPECT().GetRawTemplate(availableTemplates["test"], gomock.Any()).Return(model.Template{Id: "templateId"}, nil),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates", nil, router)
			var templates []model.Template
			json.Unmarshal(response.Body.Bytes(), &templates)
			convey.So(len(templates), convey.ShouldEqual, 1)
			convey.So(templates[0].Id, convey.ShouldEqual, "templateId")
			convey.So(response.Code, convey.ShouldEqual, 200)
		})
	})
}

func TestGenerateParsedTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Get("/api/v1/parsed_template/:templateId/", context.GenerateParsedTemplate)

	convey.Convey("Test Generate Parsed Template", t, func() {
		convey.Convey("No templateId provided", func() {
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template//?instanceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			TestUtils.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("No instanceId provided", func() {
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/1", nil, router)
			TestUtils.AssertResponse(response, "uuid can't be empty!", 400)
		})
		convey.Convey("Template with templateId not found", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(nil),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/templateId?instanceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			TestUtils.AssertResponse(response, "Can't find template by id: templateId", 404)
		})
		convey.Convey("Getting parsed component failed", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(&model.TemplateMetadata{
					Id:                  "templateId",
					TemplateDirName:     "dir",
					TemplatePlanDirName: "planDir",
				}),
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.Template{Id: "templateId"}, errors.New("failed")),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/templateId?instanceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			TestUtils.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Existing templateId provided", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(&model.TemplateMetadata{
					Id:                  "templateId",
					TemplateDirName:     "dir",
					TemplatePlanDirName: "planDir",
				}),
				templateMock.EXPECT().GetParsedTemplate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.Template{Id: "templateId"}, nil),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/parsed_template/templateId?instanceId=a5740d8a-9f4b-4711-a1a0-eae62db54474", nil, router)
			var template model.Template
			json.Unmarshal(response.Body.Bytes(), &template)
			convey.So(template.Id, convey.ShouldEqual, "templateId")
			convey.So(response.Code, convey.ShouldEqual, 200)
		})

	})
}

func TestCreateCustomTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Post("/api/v1/templates", context.CreateCustomTemplate)

	convey.Convey("Test Create Custom Template", t, func() {
		convey.Convey("Not parsable body", func() {
			response := TestUtils.SendRequest("POST", "/api/v1/templates", []byte("not_parsable"), router)
			TestUtils.AssertResponse(response, "invalid character", 500)
		})
		convey.Convey("Template without id", func() {
			body := model.Template{}
			body_bytes, _ := json.Marshal(body)
			response := TestUtils.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			TestUtils.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Template with id exists", func() {
			body := model.Template{
				Id: "templateId",
			}
			body_bytes, _ := json.Marshal(body)
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(&model.TemplateMetadata{
					Id: "templateId",
				}),
			)
			response := TestUtils.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			convey.So(response.Code, convey.ShouldEqual, 409)
		})
		convey.Convey("Adding template fails", func() {
			body := model.Template{
				Id: "templateId",
			}
			body_bytes, _ := json.Marshal(body)
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(nil),
				templateMock.EXPECT().AddAndRegisterCustomTemplate(body).Return(errors.New("failed")),
			)
			response := TestUtils.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			TestUtils.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully added", func() {
			body := model.Template{
				Id: "templateId",
			}
			body_bytes, _ := json.Marshal(body)
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(nil),
				templateMock.EXPECT().AddAndRegisterCustomTemplate(body).Return(nil),
			)
			response := TestUtils.SendRequest("POST", "/api/v1/templates", body_bytes, router)
			convey.So(response.Code, convey.ShouldEqual, 201)
		})
	})
}

func TestGetCustomTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Get("/api/v1/templates/:templateId", context.GetCustomTemplate)

	convey.Convey("Test Get Custom Template", t, func() {
		convey.Convey("No template id provided", func() {
			response := TestUtils.SendRequest("GET", "/api/v1/templates//", nil, router)
			TestUtils.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Template does not exist", func() {
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(nil),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates/templateId", nil, router)
			TestUtils.AssertResponse(response, "Template doesn't exist!", 404)
		})
		convey.Convey("Error gettting template", func() {
			templateMeta := model.TemplateMetadata{
				Id: "templateId",
			}
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(&templateMeta),
				templateMock.EXPECT().GetRawTemplate(&templateMeta, gomock.Any()).Return(model.Template{}, errors.New("failed")),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates/templateId", nil, router)
			TestUtils.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully retrieved template", func() {
			templateMeta := model.TemplateMetadata{
				Id: "templateId",
			}
			gomock.InOrder(
				templateMock.EXPECT().GetTemplateMetadataById("templateId").Return(&templateMeta),
				templateMock.EXPECT().GetRawTemplate(&templateMeta, gomock.Any()).Return(model.Template{
					Id: "templateId",
				}, nil),
			)
			response := TestUtils.SendRequest("GET", "/api/v1/templates/templateId", nil, router)
			template := model.Template{}
			json.Unmarshal(response.Body.Bytes(), &template)
			convey.So(template.Id, convey.ShouldEqual, "templateId")
			convey.So(response.Code, convey.ShouldEqual, 200)
		})
	})
}

func TestDeleteCustomTemplate(t *testing.T) {
	router, context, templateMock := prepareMocksAndRouter(t)
	router.Delete("/api/v1/templates/:templateId", context.DeleteCustomTemplate)

	convey.Convey("Test Delete Custom Template", t, func() {
		convey.Convey("No template id provided", func() {
			response := TestUtils.SendRequest("DELETE", "/api/v1/templates//", nil, router)
			TestUtils.AssertResponse(response, "templateId can't be empty!", 400)
		})
		convey.Convey("Deletion failed", func() {
			gomock.InOrder(
				templateMock.EXPECT().RemoveAndUnregisterCustomTemplate("templateId").Return(http.StatusInternalServerError, errors.New("failed")),
			)
			response := TestUtils.SendRequest("DELETE", "/api/v1/templates/templateId", nil, router)
			TestUtils.AssertResponse(response, "failed", 500)
		})
		convey.Convey("Successfully removed", func() {
			gomock.InOrder(
				templateMock.EXPECT().RemoveAndUnregisterCustomTemplate("templateId").Return(http.StatusNoContent, nil),
			)
			response := TestUtils.SendRequest("DELETE", "/api/v1/templates/templateId", nil, router)
			convey.So(response.Code, convey.ShouldEqual, http.StatusNoContent)
		})
	})
}
