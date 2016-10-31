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
package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-template-repository/model"
)

var templatesPaths map[string]string
var templatesMutex sync.RWMutex

const (
	TemplatesPath       = "./catalogData/"
	CustomTemplatesPath = TemplatesPath + "custom/"
	templateFileName    = "template.json"
)

var logger, _ = commonLogger.InitLogger("catalog")

type TemplateApi interface {
	GetTemplatePath(id string) string
	GetAvailableTemplates() map[string]string
	GetTemplatesPaths()
	AddCustomTemplate(rawTemplate model.RawTemplate, templateId string) error
	RemoveAndUnregisterCustomTemplate(templateId string) (int, error)
	GetParsedTemplate(rawTemplate model.RawTemplate, additionalReplacements map[string]string) (*model.Template, error)
	GetRawTemplate(templatePath string) (model.RawTemplate, error)
}

type TemplateApiConnector struct{}

func (t *TemplateApiConnector) GetTemplatePath(id string) string {
	templatesMutex.RLock()
	defer templatesMutex.RUnlock()

	if templatesPaths != nil {
		return templatesPaths[id]
	} else {
		return ""
	}
}

func (t *TemplateApiConnector) GetAvailableTemplates() map[string]string {
	templatesMutex.RLock()
	defer templatesMutex.RUnlock()

	result := make(map[string]string)
	for k, v := range templatesPaths {
		result[k] = v
	}
	return result
}

func (t *TemplateApiConnector) GetTemplatesPaths() {
	templatesMutex.Lock()
	defer templatesMutex.Unlock()

	templatesPaths = make(map[string]string)

	rootDir, err := ioutil.ReadDir(TemplatesPath)
	if err != nil {
		logger.Panic(err)
	}
	for _, templateTypeDir := range rootDir {
		logger.Debugf("Loading Templates of type %s", templateTypeDir.Name())
		t.loadTemplatesByType(templateTypeDir)
	}
}

func (t *TemplateApiConnector) loadTemplatesByType(templateTypeDir os.FileInfo) {
	if templateTypeDir.IsDir() {
		templateTypeDirPath := TemplatesPath + templateTypeDir.Name()
		templates, err := ioutil.ReadDir(templateTypeDirPath)
		if err != nil {
			logger.Panic(err)
		}

		for _, templateDir := range templates {
			logger.Debugf("Loading Template id: %s", templateDir.Name())
			t.loadTemplate(templateDir, templateTypeDirPath, templateTypeDir.Name())
		}
	}
}

func (t *TemplateApiConnector) loadTemplate(templateDir os.FileInfo, templateTypeDirPath, templateTypeDirName string) {
	if templateDir.IsDir() {
		templateDirPath := templateTypeDirPath + "/" + templateDir.Name()
		files, err := ioutil.ReadDir(templateDirPath)
		if err != nil {
			logger.Panic(err)
		}

		isTemplateFound := false
		for _, file := range files {
			if !file.IsDir() && file.Name() == templateFileName {
				templatesPaths[templateDir.Name()] = templateDirPath + "/" + templateFileName
				logger.Debugf("LOADED - Template id: %s with path: %s", templateDir.Name(), templatesPaths[templateDir.Name()])

				isTemplateFound = true
				break
			}
		}

		if !isTemplateFound {
			logger.Errorf("Can't find required %s file for template: %s", templateFileName, templateDir.Name())
		}
	} else {
		logger.Debug("NOT A DIR - skipping Template: ", templateDir.Name())
	}
}

func (t *TemplateApiConnector) AddCustomTemplate(rawTemplate model.RawTemplate, templateId string) error {
	rawTemplateByte, err := json.Marshal(rawTemplate)
	if err != nil {
		return err
	}

	templatePath := CustomTemplatesPath + templateId
	if err := saveTemplateInFile(templatePath, templateFileName, rawTemplateByte); err != nil {
		return err
	}

	t.GetTemplatesPaths()
	return nil
}

func (t *TemplateApiConnector) RemoveAndUnregisterCustomTemplate(templateId string) (int, error) {
	if strings.Contains(templateId, "..") {
		return http.StatusBadRequest, errors.New("illegal templateId")
	}

	if templatePath := t.GetTemplatePath(templateId); templatePath != "" {
		if !strings.HasPrefix(templatePath, CustomTemplatesPath) {
			return http.StatusForbidden, fmt.Errorf("removing template %s is forbidden", templateId)
		}
	} else {
		return http.StatusNotFound, fmt.Errorf("there is no template %q", templateId)
	}

	if err := os.RemoveAll(CustomTemplatesPath + templateId); err != nil {
		return http.StatusInternalServerError, err
	}

	t.GetTemplatesPaths()
	return http.StatusNoContent, nil
}

func (t *TemplateApiConnector) GetParsedTemplate(rawTemplate model.RawTemplate, additionalReplacements map[string]string) (*model.Template, error) {
	result := &model.Template{}

	rawTemplateByte, err := json.Marshal(&rawTemplate)
	if err != nil {
		return result, err
	}

	parsedRawTemplate := adjustParams(string(rawTemplateByte), additionalReplacements)

	if err := json.Unmarshal([]byte(parsedRawTemplate), result); err != nil {
		logger.Error("Unmarshalling Template error:", err)
		return result, err
	}

	if result.Body.Type == "" {
		result.Body.Type = model.ComponentTypeInstance
	}

	return result, nil
}

func (t *TemplateApiConnector) GetRawTemplate(templatePath string) (model.RawTemplate, error) {
	result := model.RawTemplate{}

	rawFile, err := ioutil.ReadFile(templatePath)
	if err != nil {
		logger.Errorf("Error reading file: %s", templatePath, err)
		return result, err
	}

	err = json.Unmarshal(rawFile, &result)
	return result, err
}
