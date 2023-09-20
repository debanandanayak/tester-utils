package tester_utils

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func passFunc(stageHarness *StageHarness) error {
	return nil
}

func failFunc(stageHarness *StageHarness) error {
	return errors.New("fail")
}

func TestAllStagesPass(t *testing.T) {
	definition := TesterDefinition{
		TestCases: []TestCase{
			{Slug: "test-1", TestFunc: passFunc},
			{Slug: "test-2", TestFunc: passFunc},
		},
	}

	testCasesJson, _ := json.Marshal([]map[string]string{
		{"slug": "test-1", "tester_log_prefix": "test-1", "title": "Stage #1: The First Stage"},
		{"slug": "test-2", "tester_log_prefix": "test-2", "title": "Stage #2: The Second Stage"},
	})

	tester, err := NewTester(map[string]string{
		"CODECRAFTERS_SUBMISSION_DIR":  "./test_helpers/valid_app_dir",
		"CODECRAFTERS_TEST_CASES_JSON": string(testCasesJson),
	}, definition)

	if err != nil {
		t.Error(err)
	}

	exitCode := tester.RunCLI()
	assert.Equal(t, exitCode, 0)
}

func TestOneStageFails(t *testing.T) {
	definition := TesterDefinition{
		TestCases: []TestCase{
			{Slug: "test-1", TestFunc: passFunc},
			{Slug: "test-2", TestFunc: failFunc},
		},
	}

	testCasesJson, _ := json.Marshal([]map[string]string{
		{"slug": "test-1", "tester_log_prefix": "test-1", "title": "Stage #1: The First Stage"},
		{"slug": "test-2", "tester_log_prefix": "test-2", "title": "Stage #2: The Second Stage"},
	})

	tester, err := NewTester(map[string]string{
		"CODECRAFTERS_SUBMISSION_DIR":  "./test_helpers/valid_app_dir",
		"CODECRAFTERS_TEST_CASES_JSON": string(testCasesJson),
	}, definition)

	if err != nil {
		t.Error(err)
	}

	exitCode := tester.RunCLI()
	assert.Equal(t, exitCode, 1)
}
