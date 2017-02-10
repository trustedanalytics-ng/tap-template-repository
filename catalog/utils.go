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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"

	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
	"github.com/trustedanalytics-ng/tap-template-repository/model"
)

var possibleRandChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var possibleRandDnsChars = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func adjustParams(content string, replacements map[string]string) string {
	for key, value := range replacements {
		if key == model.GetPlaceholderWithDollarPrefix(model.PlaceholderExtraEnvs) {
			rawEscapedValue, _ := json.Marshal(value)
			value = strings.Trim(string(rawEscapedValue), `"`)
		}
		content = strings.Replace(content, key, value, -1)
	}

	instanceId := replacements[model.GetPlaceholderWithDollarPrefix(model.PlaceholderInstanceID)]

	properShortDnsName := commonHttp.UuidToShortDnsName(instanceId)
	content = strings.Replace(content, model.GetPlaceholderWithDollarPrefix(model.PlaceholderShortInstanceID), properShortDnsName, -1)
	content = strings.Replace(content, model.GetPlaceholderWithDollarPrefix(model.PlaceholderIdxAndShortInstanceID), properShortDnsName, -1)

	for i := 0; i < 9; i++ {
		content = strings.Replace(content, model.GetPlaceholderWithDollarPrefix(model.PlaceholderRandom)+strconv.Itoa(i), getRandomString(10, possibleRandChars), -1)
		content = strings.Replace(content, model.GetPlaceholderWithDollarPrefix(model.PlaceholderRandomDNS)+strconv.Itoa(i), getRandomString(6, possibleRandDnsChars), -1)
	}
	return encodeByte64ToString(content)
}

func encodeByte64ToString(content string) string {
	rp := regexp.MustCompile(`\$base64\-(.*?)\"`)
	fs := rp.FindAllString(content, -1)
	for _, sub := range fs {
		sub = strings.Replace(sub, "$base64-", "", -1)
		sub = strings.Replace(sub, "\"", "", -1)
		content = strings.Replace(content, "$base64-"+sub, base64.StdEncoding.EncodeToString([]byte(sub)), -1)
	}
	return content
}

func getRandomString(length int, possibleChars []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = possibleChars[rand.Intn(len(possibleChars))]
	}
	return string(b)
}

func saveTemplateInFile(path, fileName string, file []byte) error {
	templateFilePath := path + "/" + fileName
	logger.Debugf("Saving template in file: %s", templateFilePath)

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error("Crate dir error:", err)
		return err
	}

	if err := ioutil.WriteFile(templateFilePath, file, 0666); err != nil {
		logger.Error("Save template in file error:", err)
		return err
	}
	return nil
}

func filterByPlanName(template model.Template, planName string) *model.Template {
	for i, component := range template.Body {
		template.Body[i].Deployments = getDeploymentsForPlan(component, planName)
		template.Body[i].Ingresses = getIngressesForPlan(component, planName)
		template.Body[i].Services = getServicesForPlan(component, planName)
		template.Body[i].ServiceAccounts = getServiceAccountsForPlan(component, planName)
		template.Body[i].Secrets = getSecretsAccountsForPlan(component, planName)
		template.Body[i].ConfigMaps = getConfigMapsAccountsForPlan(component, planName)
		template.Body[i].PersistentVolumeClaims = getPersistentVolumeClaimsForPlan(component, planName)
	}
	return &template
}

func getDeploymentsForPlan(k8sComponent model.KubernetesComponent, planName string) []*extensions.Deployment {
	deployments := []*extensions.Deployment{}
	for _, deployment := range k8sComponent.Deployments {
		if shouldComponentBeAttached(deployment.ObjectMeta, planName) {
			deployments = append(deployments, deployment)
		}
	}
	return deployments
}
func getIngressesForPlan(k8sComponent model.KubernetesComponent, planName string) []*extensions.Ingress {
	ingresses := []*extensions.Ingress{}
	for _, ingress := range k8sComponent.Ingresses {
		if shouldComponentBeAttached(ingress.ObjectMeta, planName) {
			ingresses = append(ingresses, ingress)
		}
	}
	return ingresses
}

func getServicesForPlan(k8sComponent model.KubernetesComponent, planName string) []*api.Service {
	services := []*api.Service{}
	for _, service := range k8sComponent.Services {
		if shouldComponentBeAttached(service.ObjectMeta, planName) {
			services = append(services, service)
		}
	}
	return services
}

func getServiceAccountsForPlan(k8sComponent model.KubernetesComponent, planName string) []*api.ServiceAccount {
	serviceAccounts := []*api.ServiceAccount{}
	for _, serviceAccount := range k8sComponent.ServiceAccounts {
		if shouldComponentBeAttached(serviceAccount.ObjectMeta, planName) {
			serviceAccounts = append(serviceAccounts, serviceAccount)
		}
	}
	return serviceAccounts
}

func getSecretsAccountsForPlan(k8sComponent model.KubernetesComponent, planName string) []*api.Secret {
	secrets := []*api.Secret{}
	for _, secret := range k8sComponent.Secrets {
		if shouldComponentBeAttached(secret.ObjectMeta, planName) {
			secrets = append(secrets, secret)
		}
	}
	return secrets
}

func getConfigMapsAccountsForPlan(k8sComponent model.KubernetesComponent, planName string) []*api.ConfigMap {
	configMaps := []*api.ConfigMap{}
	for _, configMap := range k8sComponent.ConfigMaps {
		if shouldComponentBeAttached(configMap.ObjectMeta, planName) {
			configMaps = append(configMaps, configMap)
		}
	}
	return configMaps
}

func getPersistentVolumeClaimsForPlan(k8sComponent model.KubernetesComponent, planName string) []*api.PersistentVolumeClaim {
	claims := []*api.PersistentVolumeClaim{}
	for _, pvc := range k8sComponent.PersistentVolumeClaims {
		if shouldComponentBeAttached(pvc.ObjectMeta, planName) {
			claims = append(claims, pvc)
		}
	}
	return claims
}

func shouldComponentBeAttached(meta api.ObjectMeta, planName string) bool {
	if planName == "" {
		return true
	}

	if value, exist := meta.Annotations[model.PLAN_NAMES_ANNOTATION]; exist {
		if value == "" {
			return true
		}

		for _, plan := range strings.Split(value, ",") {
			if planName == plan {
				return true
			}
		}
	} else {
		return true
	}
	return false
}
