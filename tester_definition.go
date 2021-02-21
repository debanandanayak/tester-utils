package tester_utils

import (
	"io/ioutil"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type TesterDefinition struct {
	// Example: spawn_redis_server.sh
	ExecutableFileName string

	Stages          []Stage
	AntiCheatStages []Stage
}

type stageYAML struct {
	Slug string `yaml:"slug"`
}

type courseYAML struct {
	Stages []stageYAML `yaml:"stages"`
}

// TestAgainstYaml tests whether the stage slugs in TesterDefintion match those in the course YAML at yamlPath.
func (testerDefinition TesterDefinition) TestAgainstYAML(t testing.T, yamlPath string) {
	bytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}

	c := courseYAML{}
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
