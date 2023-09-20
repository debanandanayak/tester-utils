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

func (t TesterDefinition) StageBySlug(slug string) Stage {
	for _, stage := range t.Stages {
		if stage.Slug == slug {
			return stage
		}
	}

	return Stage{}
}

type stageYAML struct {
	Slug  string `yaml:"slug"`
	Title string `yaml:"name"`
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

	slugsInYaml := []string{}
	for _, stage := range c.Stages {
		slugsInYaml = append(slugsInYaml, stage.Slug)
	}

	slugsInDefinition := []string{}
	for _, stage := range testerDefinition.Stages {
		slugsInDefinition = append(slugsInDefinition, stage.Slug)
	}

	if !assert.Equal(t, slugsInYaml, slugsInDefinition) {
		return
	}

	for stageIndex := range c.Stages {
		assert.Equal(t, stageIndex+1, testerDefinition.Stages[stageIndex].Number)
	}

	for _, stage := range c.Stages {
		stageInDefinition := testerDefinition.StageBySlug(stage.Slug)

		assert.Equal(t, stage.Title, stageInDefinition.Title)
	}
}
