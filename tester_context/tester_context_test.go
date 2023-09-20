package tester_context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresAppDir(t *testing.T) {
	_, err := GetTesterContext(map[string]string{
		"CODECRAFTERS_TEST_CASES_JSON": "[{ \"slug\": \"test\", \"tester_log_prefix\": \"test\", \"title\": \"Test\"}]",
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
		"CODECRAFTERS_TEST_CASES_JSON": "[{ \"slug\": \"test\", \"tester_log_prefix\": \"test\", \"title\": \"Test\"}]",
		"CODECRAFTERS_SUBMISSION_DIR":  "./test_helpers/valid_app_dir",
	}, "script.sh")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, context.ExecutablePath, "test_helpers/valid_app_dir/script.sh")
	assert.Equal(t, len(context.TestCases), 1)
	assert.Equal(t, context.TestCases[0].Slug, "test")
	assert.Equal(t, context.TestCases[0].TesterLogPrefix, "test")
	assert.Equal(t, context.TestCases[0].Title, "Test")
}
