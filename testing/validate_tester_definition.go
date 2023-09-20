package testing

import (
	"os"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type stageYAML struct {
	Slug  string `yaml:"slug"`
	Title string `yaml:"name"`
}

type courseYAML struct {
	Stages []stageYAML `yaml:"stages"`
}

// ValidateTesterDefinitionAgainstYAML tests whether the stage slugs in TesterDefintion match those in the course YAML at yamlPath.
func ValidateTesterDefinitionAgainstYAML(t *testing.RuntimeT, testerDefinition tester_utils.TesterDefinition, yamlPath string) {
	bytes, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}

	c := courseYAML{}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		t.Fatal(err)
	}

	slugsInYaml := []string{}
	for _, stage := range c.Stages {
		slugsInYaml = append(slugsInYaml, stage.Slug)
	}

	slugsInDefinition := []string{}
	for _, stage := range testerDefinition.TestCases {
		slugsInDefinition = append(slugsInDefinition, stage.Slug)
	}

	if !assert.Equal(t, slugsInYaml, slugsInDefinition) {
		return
	}
}
