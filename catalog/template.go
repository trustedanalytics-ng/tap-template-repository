package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/trustedanalytics/tapng-go-common/logger"
	"github.com/trustedanalytics/tapng-template-repository/model"
)

var TEMPLATES map[string]*model.TemplateMetadata
var TemplatesPath string = "./catalogData/"
var CustomTemplatesDir string = TemplatesPath + "custom/"

var logger = logger_wrapper.InitLogger("catalog")

type TemplateApi interface {
	GetTemplateMetadataById(id string) *model.TemplateMetadata
	GetAvailableTemplates() map[string]*model.TemplateMetadata
	LoadAvailableTemplates()
	AddAndRegisterCustomTemplate(template model.Template) error
	RemoveAndUnregisterCustomTemplate(templateId string) error
	GetParsedTemplate(templateMetadata *model.TemplateMetadata, catalogPath,
		instanceId, orgId, spaceId string, additionalReplacements map[string]string) (model.Template, error)
	GetRawTemplate(templateMetadata *model.TemplateMetadata, catalogPath string) (model.Template, error)
}

type Template struct{}

func (t *Template) GetTemplateMetadataById(id string) *model.TemplateMetadata {
	if TEMPLATES != nil {
		return TEMPLATES[id]
	} else {
		return nil
	}
}

func (t *Template) GetAvailableTemplates() map[string]*model.TemplateMetadata {
	if TEMPLATES != nil {
		t.LoadAvailableTemplates()
	}
	return TEMPLATES
}

func (t *Template) LoadAvailableTemplates() {
	TEMPLATES = make(map[string]*model.TemplateMetadata)
	logger.Debug("GetAvailableTemplates - need to parse catalog/ directory.")
	template_file_info, err := ioutil.ReadDir(TemplatesPath)
	if err != nil {
		logger.Panic(err)
	}
	for _, templateDir := range template_file_info {
		t.loadTemplateMetadata(templateDir)
	}
}

func (t *Template) loadTemplateMetadata(templateDir os.FileInfo) {
	if templateDir.IsDir() {
		templateDirPath := TemplatesPath + templateDir.Name()
		logger.Debug(" => ", templateDir.Name(), templateDirPath)

		plans_file_info, err := ioutil.ReadDir(templateDirPath)
		if err != nil {
			logger.Panic(err)
		}
		for _, plandir := range plans_file_info {
			t.loadPlans(plandir, templateDirPath, templateDir.Name())
		}
	}
}

func (t *Template) loadPlans(plandir os.FileInfo, templateDirPath, templateDirName string) {
	planDirPath := templateDirPath + "/" + plandir.Name()

	if plandir.IsDir() {
		logger.Debug(" ====> ", plandir.Name(), planDirPath)
		plans_content_file_info, err := ioutil.ReadDir(planDirPath)
		if err != nil {
			logger.Panic(err)
		}

		for _, plan_details := range plans_content_file_info {
			t.loadPlan(plan_details, planDirPath, plandir.Name(), templateDirName)
		}
	} else {
		logger.Debug("Skipping file: ", planDirPath)
	}
}

func (t *Template) loadPlan(plan_details os.FileInfo, planDirPath, planDirName, templateDirName string) {
	var plan_meta model.PlanMetadata
	plan_details_dir_full_name := planDirPath + "/" + plan_details.Name()
	if plan_details.IsDir() {
		logger.Debug("Skipping directory:", plan_details_dir_full_name)
	} else if plan_details.Name() == "plan.json" {
		plan_metadata_file_content, err := ioutil.ReadFile(plan_details_dir_full_name)
		if err != nil {
			logger.Fatal("Error reading file: ", plan_details_dir_full_name, err)
		}
		b := []byte(plan_metadata_file_content)
		err = json.Unmarshal(b, &plan_meta)
		if err != nil {
			logger.Fatal("Error parsing json from file: ", plan_details_dir_full_name, err)
		}

		TEMPLATES[plan_meta.Id] = &model.TemplateMetadata{
			Id:                  plan_meta.Id,
			TemplateDirName:     templateDirName,
			TemplatePlanDirName: planDirName,
		}
	} else {
		logger.Debug(" -----------> ", plan_details.Name(), plan_details_dir_full_name)
	}
}

func (t *Template) AddAndRegisterCustomTemplate(template model.Template) error {
	templateDir := CustomTemplatesDir + template.Id + "/k8s"
	templatePlanDir := CustomTemplatesDir + template.Id

	for i, pvc := range template.Body.PersistentVolumeClaims {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("persistentvolumeclaim_%d.json", i), pvc)
		if err != nil {
			return err
		}
	}
	for i, rc := range template.Body.Deployments {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("deployment_%d.json", i), rc)
		if err != nil {
			return err
		}
	}
	for i, ing := range template.Body.Ingresses {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("ingress_%d.json", i), ing)
		if err != nil {
			return err
		}
	}
	for i, svc := range template.Body.Services {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("service_%d.json", i), svc)
		if err != nil {
			return err
		}
	}
	for i, svcAccount := range template.Body.ServiceAccounts {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("account_%d.json", i), svcAccount)
		if err != nil {
			return err
		}
	}
	for i, secret := range template.Body.Secrets {
		err := save_k8s_file_in_dir(templateDir, fmt.Sprintf("secret_%d.json", i), secret)
		if err != nil {
			return err
		}
	}

	err := save_k8s_file_in_dir(templateDir, "hooks.json", template.Hooks)
	if err != nil {
		return err
	}

	err = save_k8s_file_in_dir(templatePlanDir, "component.json", model.KubernetesComponent{Type: template.Body.Type})
	if err != nil {
		return err
	}

	plan := model.PlanMetadata{Id: template.Id}
	err = save_k8s_file_in_dir(templatePlanDir, "plan.json", plan)
	if err != nil {
		return err
	}

	t.LoadAvailableTemplates()
	return nil
}

func (t *Template) RemoveAndUnregisterCustomTemplate(templateId string) error {
	if strings.Contains(templateId, "..") {
		return errors.New("Illegal templateId")
	}
	templateDir := CustomTemplatesDir + templateId
	err := os.RemoveAll(templateDir)
	if err != nil {
		return err
	}

	t.LoadAvailableTemplates()
	return nil
}

func (t *Template) GetParsedTemplate(templateMetadata *model.TemplateMetadata, catalogPath,
	instanceId, orgId, spaceId string, additionalReplacements map[string]string) (model.Template, error) {

	result, err := GetParsedTemplate(catalogPath, instanceId, orgId, spaceId, templateMetadata, additionalReplacements)
	if err != nil {
		return *result, err
	}
	result.Id = templateMetadata.Id

	return *result, nil
}

func (t *Template) GetRawTemplate(templateMetadata *model.TemplateMetadata, catalogPath string) (model.Template, error) {
	blueprint, err := GetKubernetesBlueprint(catalogPath, templateMetadata.TemplateDirName, templateMetadata.TemplatePlanDirName)
	if err != nil {
		return model.Template{}, err
	}

	result, err := CreateTemplateFromBlueprint(blueprint, true)
	if err != nil {
		return *result, err
	}
	result.Id = templateMetadata.Id

	return *result, nil
}
