package catalog

import (
	"encoding/json"
	"os"

	"github.com/trustedanalytics/tapng-image-factory/logger"
	"github.com/trustedanalytics/tapng-template-repository/model"
)

var TEMPLATES map[string]*model.TemplateMetadata
var TemplatesPath string = "./catalogData/"
var CustomTemplatesDir string = TemplatesPath + "custom/"

var logger = logger_wrapper.InitLogger("api")

func GetTemplateMetadataByIdTemplateMock(id string) *model.TemplateMetadata {
	template := model.TemplateMetadata{
		Id:                  id,
		TemplateDirName:     "dir",
		TemplatePlanDirName: "planDir",
	}
	return &template
}

func GetTemplateMetadataByIdNilMock(id string) *model.TemplateMetadata {
	return nil
}

func GetAvailableTemplates() map[string]*model.TemplateMetadata {
	return TEMPLATES
}

func LoadAvailableTemplates() {

}

func loadTemplateMetadata(templateDir os.FileInfo) {

}

func loadPlans(plandir os.FileInfo, templateDirPath, templateDirName string) {

}

func loadPlan(plan_details os.FileInfo, planDirPath, planDirName, templateDirName string) {

}

func AddAndRegisterCustomTemplate(template model.Template) error {
	return nil
}

func RemoveAndUnregisterCustomTemplate(templateId string) error {
	return nil
}

func GetParsedTemplate(templateMetadata *model.TemplateMetadata, catalogPath, instanceId, orgId, spaceId string) (model.Template, error) {
	result := model.Template{Id: templateMetadata.Id}
	return result, nil
}

func GetRawTemplate(templateMetadata *model.TemplateMetadata, catalogPath string) (model.Template, error) {
	result := model.Template{Id: templateMetadata.Id}
	return result, nil
}

func GetParsedJobHooks(jobs []string, instanceId, svcMetaId, planMetaId, org, space string) ([]*model.JobHook, error) {
	parsedJobs := []string{}
	return unmarshallJobs(parsedJobs)
}

func unmarshallJobs(jobs []string) ([]*model.JobHook, error) {
	result := []*model.JobHook{}
	for _, job := range jobs {
		jobHook := &model.JobHook{}
		err := json.Unmarshal([]byte(job), jobHook)
		if err != nil {
			logger.Error("Unmarshalling JobHook error:", err)
			return result, err
		}
		result = append(result, jobHook)
	}
	return result, nil
}

func GetJobHooks(catalogPath string, temp *model.TemplateMetadata) ([]string, error) {
	jobHooks := []string{}
	return jobHooks, nil
}
