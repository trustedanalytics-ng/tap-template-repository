package catalog

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

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
