package tester_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresAppDir(t *testing.T) {
	_, err := GetTesterContext(map[string]string{
		"CODECRAFTERS_CURRENT_STAGE_SLUG": "init",
	}, "script.sh")
	if !assert.Error(t, err) {
		t.FailNow()
	}
}

func TestRequiresCurrentStageSlug(t *testing.T) {
	_, err := GetTesterContext(map[string]string{
		"CODECRAFTERS_SUBMISSION_DIR": "./test_helpers/valid_app_dir",
	}, "script.sh")
	if !assert.Error(t, err) {
		t.FailNow()
	}
}

func TestSuccessParse(t *testing.T) {
	context, err := GetTesterContext(map[string]string{
		"CODECRAFTERS_CURRENT_STAGE_SLUG": "init",
		"CODECRAFTERS_SUBMISSION_DIR":     "./test_helpers/valid_app_dir",
	}, "script.sh")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, context.executablePath, "test_helpers/valid_app_dir/script.sh")
	assert.Equal(t, context.currentStageSlug, "init")
}
