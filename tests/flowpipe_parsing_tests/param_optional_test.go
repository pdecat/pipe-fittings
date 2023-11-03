package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turbot/pipe-fittings/misc"
)

func TestParamOptional(t *testing.T) {
	assert := assert.New(t)

	pipelines, _, err := misc.LoadPipelines(context.TODO(), "./pipelines/param_optional.fp")
	assert.Nil(err, "error found")

	validateMyParam := pipelines["local.pipeline.test_param_optional"]
	if validateMyParam == nil {
		assert.Fail("test_param_optional pipeline not found")
		return
	}

	stringValid := map[string]interface{}{}

	assert.Equal(0, len(validateMyParam.ValidatePipelineParam(stringValid)))
}
