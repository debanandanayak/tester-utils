package tester_utils

import (
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// testerContext holds all flags that a user has passed in
type testerContext struct {
	executablePath   string
	isDebug          bool
	currentStageSlug string
}

type YAMLConfig struct {
	Debug bool `yaml:"debug"`
}

func (c testerContext) print() {
	fmt.Println("Debug =", c.isDebug)
	fmt.Println("Stage =", c.currentStageSlug)
}

// GetContext parses flags and returns a Context object
func getTesterContext(env map[string]string, executableFileName string) (testerContext, error) {
	submissionDir, ok := env["CODECRAFTERS_SUBMISSION_DIR"]
	if !ok {
		return testerContext{}, fmt.Errorf("CODECRAFTERS_SUBMISSION_DIR env var not found")
	}

	currentStageSlug, ok := env["CODECRAFTERS_CURRENT_STAGE_SLUG"]
	if !ok {
		return testerContext{}, fmt.Errorf("CODECRAFTERS_CURRENT_STAGE_SLUG env var not found")
	}

	configPath := path.Join(submissionDir, "codecrafters.yml")
	executablePath := path.Join(submissionDir, executableFileName)

	yamlConfig, err := readFromYAML(configPath)
	if err != nil {
		return testerContext{}, err
	}

	// TODO: test if executable exists?

	return testerContext{
		executablePath:   executablePath,
		isDebug:          yamlConfig.Debug,
		currentStageSlug: currentStageSlug,
	}, nil
}

func readFromYAML(configPath string) (YAMLConfig, error) {
	c := &YAMLConfig{}

	fileContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return YAMLConfig{}, err
	}

	if err := yaml.Unmarshal(fileContents, c); err != nil {
		return YAMLConfig{}, err
	}

	return *c, nil
}
