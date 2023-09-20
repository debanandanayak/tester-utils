package tester_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testingInterface "github.com/mitchellh/go-testing-interface"
)

func TestTestAgainstYAMLFailure(t *testing.T) {
	definition := TesterDefinition{
		TestCases: []TestCase{
			{Slug: "test-1"},
			{Slug: "test-2"},
		},
	}

	runtimeT := &testingInterface.RuntimeT{}

	yamlPath := "test_helpers/tester_definition_test/course_definition.yml"
	ValidateTesterDefinitionAgainstYAML(runtimeT, definition, yamlPath)

	assert.True(t, runtimeT.Failed())
}

func TestTestAgainstYAMLSuccess(t *testing.T) {
	definition := TesterDefinition{
		Stages: []Stage{
			{Slug: "init", Title: "Bind to a port", Number: 1},
			{Slug: "ping-pong", Title: "Respond to PING", Number: 2},
			{Slug: "ping-pong-multiple", Title: "Respond to multiple PINGs", Number: 3},
			{Slug: "concurrent-clients", Title: "Handle concurrent clients", Number: 4},
			{Slug: "echo", Title: "Implement the ECHO command", Number: 5},
			{Slug: "set_get", Title: "Implement the SET & GET commands", Number: 6},
			{Slug: "expiry", Title: "Expiry", Number: 7},
		},
	}

	runtimeT := &testingInterface.RuntimeT{}

	yamlPath := "test_helpers/tester_definition_test/course_definition.yml"
	definition.TestAgainstYAML(runtimeT, yamlPath)

	assert.False(t, runtimeT.Failed())
}
