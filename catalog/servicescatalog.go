/**
 * Copyright (c) 2015 Intel Corporation
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
	"io/ioutil"
	"log"
	"strings"
	//"k8s.io/kubernetes/pkg/api"
	"os"
)

type TapServices struct {
	Services []ServiceTemplate `json:"services"`
}

type ServiceTemplate struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Bindable    bool           `json:"bindable"`
	Tags        []string       `json:"tags"`
	Plans       []PlanMetadata `json:"plans"`
	Metadata    ServiceTemplateMetadata`json:"metadata"`
}

type ServiceTemplateMetadata struct {
	DisplayName		string `json:"displayName"`
	ImageUrl		string `json:"imageUrl"`
	LongDescription		string `json:"longDescription"`
	ProviderDisplayName	string `json:"providerDisplayName"`
	DocumentationUrl	string `json:"documentationUrl"`
	SupportUrl		string `json:"supportUrl"`
}



type PlanMetadata struct {
	Id                              string `json:"id"`
	Name                            string `json:"name"`
	Description                     string `json:"description"`
	Free                            bool   `json:"free"`
	SecretsDependencies             []string `json:"secretsDependencies"`
	CreatesSecrets                  []string `json:"createsSecrets"`
	ReplicationControllersTemplates []string `json:"-"`
	K8SServicesTemplates            []string `json:"-"`
}

func FindServiceByIdInServices(serviceId string, services []ServiceTemplate) (ServiceTemplate, error) {
	for _, service := range services {
		if service.Id == serviceId {
			return service, nil
		}
	}
	return ServiceTemplate{}, errors.New("No such service by ID: " + serviceId)
}
func FindPlanByIdInService(planId string, service ServiceTemplate) (PlanMetadata, error){
	for _, plan := range service.Plans {
		if plan.Id == planId {
			return plan, nil
		}
	}
	return PlanMetadata{}, errors.New("No such plan by ID: " + planId)
}

func GetTapServices() (TapServices, error){
	tapServices, _ := getServicesFromDir("./service-templates")

	return tapServices, nil
}

func GetTapSecrets() ([]string, error){
	tapSecrets, _ := getSecretsFromDir("./secrets-templates")
	return tapSecrets, nil
}

func getSecretsFromDir(catalogPath string) ([]string, error) {
	tapSecrests := []string{}
	catalog_file_info := readDirOrFail(catalogPath)
	for _, f := range catalog_file_info {
		if ! isFileAJsonAndType(f, "") { //any json
			continue
		}
		secret := readFileOrFail(catalogPath+"/"+f.Name())//unmarsahelJsonOrFail(f.Name(), &api.Secret{})
		tapSecrests = append(tapSecrests, secret)
	}
	return tapSecrests, nil
}

func getServicesFromDir(catalogPath string) (TapServices, error) {
	tapServices := TapServices{}
	tapServices.Services = []ServiceTemplate{}

	catalog_file_info := readDirOrFail(catalogPath)

	for _, serviceDir := range catalog_file_info {
		serviceMetadata, _ := getServiceMetadataNoPlansFromDir(catalogPath+"/"+serviceDir.Name())

		serviceMetadata.Plans, _ = getServicePlansFromDir(catalogPath+"/"+serviceDir.Name())

		tapServices.AppendServices(serviceMetadata)


	}
	return tapServices, nil
}

func (s *TapServices) AppendServices(service ServiceTemplate) []ServiceTemplate {
	s.Services = append(s.Services, service)
	return s.Services
}

func (s *PlanMetadata) AppendReplicationControllers(replicationControllers []string) []string {
	for _, r := range replicationControllers {
		s.ReplicationControllersTemplates = append(s.ReplicationControllersTemplates, r)
	}
	return s.ReplicationControllersTemplates
}

func (s *PlanMetadata) AppendK8SServices(services []string) []string {
	for _, svc := range services {
		s.K8SServicesTemplates = append(s.K8SServicesTemplates, svc)
	}
	return s.K8SServicesTemplates
}


func getServiceMetadataNoPlansFromDir(catalogPath string) (ServiceTemplate, error){
	serviceMetadata := ServiceTemplate{}
	fcontent := readFileOrFail(catalogPath + "/service.json")
	err := json.Unmarshal([]byte(fcontent), &serviceMetadata)
	serviceMetadata.Plans = []PlanMetadata{}
	return serviceMetadata, err
}

func getServicePlansFromDir(catalogPath string) ([]PlanMetadata, error){
	servicePlans := []PlanMetadata{}
	catalog_file_info := readDirOrFail(catalogPath)

	for _, planDir := range catalog_file_info {
		if ! planDir.IsDir() {
			continue
		}
		servicePlan, _ := getServicePlanFromFile(catalogPath+"/"+planDir.Name() + "/plan.json")
		replicationControllers, _ := getReplicationControllersFromDir(catalogPath+"/"+planDir.Name())
		k8sServices, _ := getK8SServicesFromDir(catalogPath+"/"+planDir.Name())

		servicePlan.AppendReplicationControllers(replicationControllers)
		servicePlan.AppendK8SServices(k8sServices)

		servicePlans = append(servicePlans, servicePlan)
	}

	return servicePlans, nil
}

func getServicePlanFromFile(planFile string) (PlanMetadata, error) {
	servicePlan := PlanMetadata{}
	fcontent := readFileOrFail(planFile)
	err := json.Unmarshal([]byte(fcontent), &servicePlan)
	return servicePlan, err
}

func getReplicationControllersFromDir(catalogPath string) ([]string, error){//api.ReplicationController, error) {
	replicationControllers := []string{} //api.ReplicationController{}
	catalogPath = catalogPath + "/k8s"
	files_in_path := readDirOrFail(catalogPath)
	for _, f := range files_in_path {
		if ! isFileAJsonAndType(f, "replicationcontroller") {
			continue
		}

		replicationController := readFileOrFail(catalogPath+"/"+f.Name())//unmarsahelJsonOrFail(f.Name(), &api.ReplicationController{})
		replicationControllers = append(replicationControllers, replicationController)
	}
	return replicationControllers, nil
}

func getK8SServicesFromDir(catalogPath string) ([]string, error){ //([]api.Service, error){
	services := []string{} //[]api.Service{}
	catalogPath = catalogPath + "/k8s"
	files_in_path := readDirOrFail(catalogPath)
	for _, f := range files_in_path {
		if ! isFileAJsonAndType(f, "service") {
			continue
		}

		service := readFileOrFail(catalogPath+"/"+f.Name())//unmarsahelJsonOrFail(f.Name(),&api.Service{})
		services = append(services, service)
	}
	return services, nil
}

func isFileAJsonAndType(f os.FileInfo, fileType string) bool {
	return strings.HasPrefix(f.Name(), fileType) && strings.HasSuffix(f.Name(), ".json")
}

func unmarsahelJsonOrFail(file string, structType interface{}) interface{}{
	fcontent := readFileOrFail(file)

	err := json.Unmarshal([]byte(fcontent), structType)
	if err != nil {
		log.Panicln(err)
	}
	return structType
}

func readDirOrFail(catalogPath string) []os.FileInfo {
	log.Println("Read catalog:", catalogPath)
	catalog_file_info, err := ioutil.ReadDir(catalogPath)
	if err != nil {
		log.Panicln(err)
	}
	return catalog_file_info
}

func readFileOrFail(file string) string {
	log.Println("Read file:", file)
	fcontent, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic("Error reading file: ", fcontent, err)
	}
	return string(fcontent)
}