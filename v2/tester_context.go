package tester_utils

import (
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// TesterContext holds all flags that a user has passed in
type TesterContext struct {
	executablePath   string
	isDebug          bool
	currentStageSlug string
	apiKey           string // Not used right now, but will be once our endpoint has auth
}

type YAMLConfig struct {
	Debug bool `yaml:"debug"`
}

func (c TesterContext) print() {
	fmt.Println("Debug =", c.isDebug)
	fmt.Println("Stage =", c.currentStageSlug)
}

// GetContext parses flags and returns a Context object
func GetTesterContext(env map[string]string, executableFileName string) (TesterContext, error) {
	submissionDir, ok := env["CODECRAFTERS_SUBMISSION_DIR"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_SUBMISSION_DIR env var not found")
	}

	currentStageSlug, ok := env["CODECRAFTERS_CURRENT_STAGE_SLUG"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_CURRENT_STAGE_SLUG env var not found")
	}

	configPath := path.Join(submissionDir, "codecrafters.yml")
	executablePath := path.Join(submissionDir, executableFileName)

	yamlConfig, err := ReadFromYAML(configPath)
	if err != nil {
		return TesterContext{}, err
	}

	// TODO: test if executable exists?

	return TesterContext{
		executablePath:   executablePath,
		isDebug:          yamlConfig.Debug,
		currentStageSlug: currentStageSlug,
		apiKey:           "dummy",
	}, nil
}

func ReadFromYAML(configPath string) (YAMLConfig, error) {
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
