package tester_utils

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type TesterDefinition struct {
	ExecutableFileName string
	Stages             []Stage
	AntiCheatStages    []Stage
}

type StageYAML struct {
	Slug string `yaml:"slug"`
}

type CourseYAML struct {
	Stages []StageYAML `yaml:"stages"`
}

// TestAgainstYaml tests whether the stage slugs in TesterDefintion match those in the course YAML at yamlPath.
func (testerDefinition TesterDefinition) TestAgainstYAML(t *testing.T, yamlPath string) {
	bytes, err := ioutil.ReadFile("test_helpers/course_definition.yml")
	if err != nil {
		t.Fatal(err)
	}

	c := CourseYAML{}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		t.Fatal(err)
	}

	stagesInYaml := []string{}
	for _, stage := range c.Stages {
		stagesInYaml = append(stagesInYaml, stage.Slug)
	}

	stagesInDefinition := []string{}
	for _, stage := range testerDefinition.Stages {
		stagesInDefinition = append(stagesInDefinition, stage.Slug)
	}

	assert.Equal(t, stagesInYaml, stagesInDefinition)
}
