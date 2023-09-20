package testing

import (
	"testing"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/stretchr/testify/assert"

	testingInterface "github.com/mitchellh/go-testing-interface"
)

func TestTestAgainstYAMLFailure(t *testing.T) {
	definition := tester_utils.TesterDefinition{
		TestCases: []tester_utils.TestCase{
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
	definition := tester_utils.TesterDefinition{
		TestCases: []tester_utils.TestCase{
			{Slug: "init"},
			{Slug: "ping-pong"},
			{Slug: "ping-pong-multiple"},
			{Slug: "concurrent-clients"},
			{Slug: "echo"},
			{Slug: "set_get"},
			{Slug: "expiry"},
		},
	}

	runtimeT := &testingInterface.RuntimeT{}

	yamlPath := "test_helpers/tester_definition_test/course_definition.yml"
	ValidateTesterDefinitionAgainstYAML(runtimeT, definition, yamlPath)

	assert.False(t, runtimeT.Failed())
}
