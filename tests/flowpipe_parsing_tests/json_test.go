package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turbot/pipe-fittings/load_mod"
)

func TestJsonSimple(t *testing.T) {
	assert := assert.New(t)

	pipelines, _, err := load_mod.LoadPipelines(context.TODO(), "./pipelines/json.fp")
	assert.Nil(err, "error found")

	assert.GreaterOrEqual(len(pipelines), 1, "wrong number of pipelines")

	if pipelines["local.pipeline.json"] == nil {
		assert.Fail("json pipeline not found")
		return
	}

	step := pipelines["local.pipeline.json"].GetStep("echo.json")
	if step == nil {
		assert.Fail("echo.json step not found")
		return
	}
}

func TestJsonFor(t *testing.T) {
	assert := assert.New(t)

	pipelines, _, err := load_mod.LoadPipelines(context.TODO(), "./pipelines/json.fp")
	assert.Nil(err, "error found")

	assert.GreaterOrEqual(len(pipelines), 1, "wrong number of pipelines")

	if pipelines["local.pipeline.json_for"] == nil {
		assert.Fail("json_for pipeline not found")
		return
	}

	step := pipelines["local.pipeline.json_for"].GetStep("echo.json")
	if step == nil {
		assert.Fail("echo.json step not found")
		return
	}
}
