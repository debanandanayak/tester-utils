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
			{Slug: "init", Title: "Bind to a port"},
			{Slug: "ping-pong", Title: "Respond to PING"},
			{Slug: "ping-pong-multiple", Title: "Respond to multiple PINGs"},
			{Slug: "concurrent-clients", Title: "Handle concurrent clients"},
			{Slug: "echo", Title: "Implement the ECHO command"},
			{Slug: "set_get", Title: "Implement the SET & GET commands"},
			{Slug: "expiry", Title: "Expiry"},
		},
	}

	runtimeT := &testingInterface.RuntimeT{}

	yamlPath := "test_helpers/tester_definition_test/course_definition.yml"
	definition.TestAgainstYAML(runtimeT, yamlPath)

	assert.False(t, runtimeT.Failed())
}
