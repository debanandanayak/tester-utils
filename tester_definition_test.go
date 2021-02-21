package tester_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testingInterface "github.com/mitchellh/go-testing-interface"
)

func TestTestAgainstYAMLFailure(t *testing.T) {
	definition := TesterDefinition{
		Stages: []Stage{
			{Slug: "test-1"},
			{Slug: "test-2"},
		},
	}

	runtimeT := &testingInterface.RuntimeT{}

	yamlPath := "test_helpers/tester_definition_test/course_definition.yml"
	definition.TestAgainstYAML(runtimeT, yamlPath)

	assert.True(t, runtimeT.Failed())
}

func TestTestAgainstYAMLSuccess(t *testing.T) {
	definition := TesterDefinition{
		Stages: []Stage{
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
	definition.TestAgainstYAML(runtimeT, yamlPath)

	assert.False(t, runtimeT.Failed())
}
