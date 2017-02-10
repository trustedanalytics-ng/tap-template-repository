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
// Automatically generated by MockGen. DO NOT EDIT!
// Source: catalog/template.go

package catalog

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/trustedanalytics-ng/tap-template-repository/model"
)

// Mock of TemplateApi interface
type MockTemplateApi struct {
	ctrl     *gomock.Controller
	recorder *_MockTemplateApiRecorder
}

// Recorder for MockTemplateApi (not exported)
type _MockTemplateApiRecorder struct {
	mock *MockTemplateApi
}

func NewMockTemplateApi(ctrl *gomock.Controller) *MockTemplateApi {
	mock := &MockTemplateApi{ctrl: ctrl}
	mock.recorder = &_MockTemplateApiRecorder{mock}
	return mock
}

func (_m *MockTemplateApi) EXPECT() *_MockTemplateApiRecorder {
	return _m.recorder
}

func (_m *MockTemplateApi) GetTemplatePath(id string) string {
	ret := _m.ctrl.Call(_m, "GetTemplatePath", id)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockTemplateApiRecorder) GetTemplatePath(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTemplatePath", arg0)
}

func (_m *MockTemplateApi) GetAvailableTemplates() map[string]string {
	ret := _m.ctrl.Call(_m, "GetAvailableTemplates")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

func (_mr *_MockTemplateApiRecorder) GetAvailableTemplates() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAvailableTemplates")
}

func (_m *MockTemplateApi) GetTemplatesPaths() {
	_m.ctrl.Call(_m, "GetTemplatesPaths")
}

func (_mr *_MockTemplateApiRecorder) GetTemplatesPaths() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTemplatesPaths")
}

func (_m *MockTemplateApi) AddCustomTemplate(rawTemplate model.RawTemplate, templateId string) error {
	ret := _m.ctrl.Call(_m, "AddCustomTemplate", rawTemplate, templateId)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockTemplateApiRecorder) AddCustomTemplate(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddCustomTemplate", arg0, arg1)
}

func (_m *MockTemplateApi) RemoveAndUnregisterCustomTemplate(templateId string) (int, error) {
	ret := _m.ctrl.Call(_m, "RemoveAndUnregisterCustomTemplate", templateId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTemplateApiRecorder) RemoveAndUnregisterCustomTemplate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveAndUnregisterCustomTemplate", arg0)
}

func (_m *MockTemplateApi) GetParsedTemplate(rawTemplate model.RawTemplate, additionalReplacements map[string]string, planName string) (*model.Template, error) {
	ret := _m.ctrl.Call(_m, "GetParsedTemplate", rawTemplate, additionalReplacements, planName)
	ret0, _ := ret[0].(*model.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTemplateApiRecorder) GetParsedTemplate(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetParsedTemplate", arg0, arg1, arg2)
}

func (_m *MockTemplateApi) GetRawTemplate(templatePath string) (model.RawTemplate, error) {
	ret := _m.ctrl.Call(_m, "GetRawTemplate", templatePath)
	ret0, _ := ret[0].(model.RawTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTemplateApiRecorder) GetRawTemplate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRawTemplate", arg0)
}
