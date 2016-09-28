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
package model

const (
	PLACEHOLDER_ORG   = "org"
	PLACEHOLDER_SPACE = "space"

	PLACEHOLDER_DOMAIN_NAME = "domain_name"
	PLACEHOLDER_IMAGE       = "image"
	PLACEHOLDER_HOSTNAME    = "hostname"

	PLACEHOLDER_INSTANCE_ID               = "instance_id"
	PLACEHOLDER_IDX_AND_SHORT_INSTANCE_ID = "idx_and_short_instance_id"
	PLACEHOLDER_SHORT_INSTANCE_ID         = "short_instance_id"
	PLACEHOLDER_BOUND_INSTANCE_ID         = "bound_instance_id"

	PLACEHOLDER_BROKER_SHORT_INSTANCE_ID = "broker_short_instance_id"
	PLACEHOLDER_BROKER_INSTANCE_ID       = "broker_instance_id"

	PLACEHOLDER_RANDOM     = "random"
	PLACEHOLDER_RANDOM_DNS = "random_dns"

	PLACEHOLDER_OFFERING_ID           = "offering_id"
	PLACEHOLDER_PLAN_ID               = "plan_id"
	PLACEHOLDER_SOURCE_OFFERING_ID    = "source_offering_id"
	PLACEHOLDER_SOURCE_PLAN_ID_PREFIX = "source_plan_id-"
)

func GetPlaceholderWithDollarPrefix(placeholder string) string {
	return "$" + placeholder
}

func GetPrefixedSourcePlanName(planName string) string {
	return PLACEHOLDER_SOURCE_PLAN_ID_PREFIX + planName
}