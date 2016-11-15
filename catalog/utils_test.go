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
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"

	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-template-repository/model"
)

func TestAdjustParams(t *testing.T) {
	instanceId := "test-instance-id"
	properShortDnsName := util.UuidToShortDnsName(instanceId)

	replacements := map[string]string{
		model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_INSTANCE_ID): instanceId,
	}

	convey.Convey("Test adjustParams", t, func() {
		convey.Convey("Test PLACEHOLDER_SHORT_INSTANCE_ID", func() {
			content := model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_SHORT_INSTANCE_ID)
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldEqual, properShortDnsName)
		})

		convey.Convey("Test PLACEHOLDER_IDX_AND_SHORT_INSTANCE_ID", func() {
			content := model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_IDX_AND_SHORT_INSTANCE_ID)
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldEqual, properShortDnsName)
		})

		convey.Convey("Test PLACEHOLDER_RANDOM missing index", func() {
			content := model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_RANDOM)
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldEqual, content)
		})

		convey.Convey("Test PLACEHOLDER_RANDOM with index", func() {
			content := model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_RANDOM) + "1"
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldNotEqual, content)
		})

		convey.Convey("Test PLACEHOLDER_EXTRA_ENVS", func() {
			sampleValue := "sample value"
			replacements[model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_EXTRA_ENVS)] = sampleValue

			content := model.GetPlaceholderWithDollarPrefix(model.PLACEHOLDER_EXTRA_ENVS)
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldEqual, sampleValue)
		})

		convey.Convey("Test base64 encoding", func() {
			content := `{"pass": "$base64-password"}`
			response := adjustParams(content, replacements)
			convey.So(response, convey.ShouldEqual, `{"pass": "cGFzc3dvcmQ="}`)
		})
	})
}

func TestFilterByPlanName(t *testing.T) {
	const planA = "plan-a"
	const planB = "plan-b"
	const notExistingPlan = "not-exist"

	template := model.Template{
		Body: model.KubernetesComponent{
			Deployments: []*extensions.Deployment{
				{ObjectMeta: getObjectMetaWithPlanNamesInAnnotation([]string{planA, planB})},
				{ObjectMeta: getObjectMetaWithPlanNamesInAnnotation([]string{planA})},
				{ObjectMeta: getObjectMetaWithPlanNamesInAnnotation([]string{model.EMPTY_PLAN_NAME})},
				{},
			},
		},
	}

	convey.Convey("Test filterByPlanName method", t, func() {
		convey.Convey("Should returns all plans with specific name", func() {
			response := filterByPlanName(template, planA)
			convey.So(len(response.Body.Deployments), convey.ShouldEqual, 4)

			response = filterByPlanName(template, planB)
			convey.So(len(response.Body.Deployments), convey.ShouldEqual, 3)

			response = filterByPlanName(template, model.EMPTY_PLAN_NAME)
			convey.So(len(response.Body.Deployments), convey.ShouldEqual, 4)

			response = filterByPlanName(template, notExistingPlan)
			convey.So(len(response.Body.Deployments), convey.ShouldEqual, 2)
		})
	})
}

func getObjectMetaWithPlanNamesInAnnotation(planNames []string) api.ObjectMeta {
	meta := api.ObjectMeta{
		Annotations: make(map[string]string),
	}
	meta.Annotations[model.PLAN_NAMES_ANNOTATION] = strings.Join(planNames, ",")
	return meta
}
