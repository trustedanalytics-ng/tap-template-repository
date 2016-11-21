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

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPlaceholder(t *testing.T) {
	convey.Convey("Test GetPlaceholderWithDollarPrefix", t, func() {
		convey.Convey("Placeholder provided returns $placeholder", func() {
			result := GetPlaceholderWithDollarPrefix("placeholder")
			convey.So(result, convey.ShouldEqual, "$placeholder")
		})
	})
}

func TestDefaultReplacements(t *testing.T) {
	convey.Convey("Test GetMapWithDefaultReplacementsIfKeyNotExists", t, func() {
		convey.Convey("Empty map provided will return map with defaults", func() {
			result := make(map[string]string)
			result = GetMapWithDefaultReplacementsIfKeyNotExists(result)

			convey.Convey("Test PLACEHOLDER_ORG exists", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_ORG)], convey.ShouldEqual, defaultOrg)
			})
			convey.Convey("Test PLACEHOLDER_SPACE exists", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_SPACE)], convey.ShouldEqual, defaultSpace)
			})
			convey.Convey("Test PLACEHOLDER_CEPH_MONITORS_LIST exists", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_CEPH_MONITORS_LIST)], convey.ShouldEqual, "")
			})
			convey.Convey("Test PLACEHOLDER_MEMORY_LIMIT exists", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_MEMORY_LIMIT)], convey.ShouldEqual, defaultMemoryLimit)
			})
		})
		convey.Convey("Map with org defined will return map with defaults and defined org", func() {
			result := make(map[string]string)
			result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_ORG)] = "myOrg"
			result = GetMapWithDefaultReplacementsIfKeyNotExists(result)
			convey.Convey("Test PLACEHOLDER_ORG is myOrg", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_ORG)], convey.ShouldEqual, "myOrg")
			})
			convey.Convey("Test PLACEHOLDER_SPACE exists", func() {
				convey.So(result[GetPlaceholderWithDollarPrefix(PLACEHOLDER_SPACE)], convey.ShouldEqual, defaultSpace)
			})
		})
	})
}
