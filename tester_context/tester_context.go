package tester_context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// TesterContextTestCase represents one element in the CODECRAFTERS_TEST_CASES environment variable
type TesterContextTestCase struct {
	// Slug is the slug of the test case. Example: "bind-to-port"
	Slug string `json:"slug"`

	// TesterLogPrefix is the prefix that'll be used for all logs emitted by the tester. Example: "stage-1"
	TesterLogPrefix string `json:"tester_log_prefix"`

	// Title is the title of the test case. Example: "Stage #1: Bind to a port"
	Title string `json:"title"`
}

// TesterContext holds all flags passed in via environment variables, or from the codecrafters.yml file
type TesterContext struct {
	ExecutablePath string
	IsDebug        bool
	TestCases      []TesterContextTestCase
}

type yamlConfig struct {
	Debug bool `yaml:"debug"`
}

func (c TesterContext) Print() {
	fmt.Println("Debug =", c.IsDebug)
}

// GetContext parses flags and returns a Context object
func GetTesterContext(env map[string]string, executableFileName string) (TesterContext, error) {
	submissionDir, ok := env["CODECRAFTERS_SUBMISSION_DIR"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_SUBMISSION_DIR env var not found")
	}

	testCasesJson, ok := env["CODECRAFTERS_TEST_CASES_JSON"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES env var not found")
	}

	testCases := []TesterContextTestCase{}
	if err := json.Unmarshal([]byte(testCasesJson), &testCases); err != nil {
		return TesterContext{}, fmt.Errorf("failed to parse CODECRAFTERS_TEST_CASES: %s", err)
	}

	configPath := path.Join(submissionDir, "codecrafters.yml")
	executablePath := path.Join(submissionDir, executableFileName)

	yamlConfig, err := readFromYAML(configPath)
	if err != nil {
		return TesterContext{}, err
	}

	if len(testCases) == 0 {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES is empty")
	}

	// TODO: test if executable exists?

	return TesterContext{
		ExecutablePath: executablePath,
		IsDebug:        yamlConfig.Debug,
		TestCases:      testCases,
	}, nil
}

func readFromYAML(configPath string) (yamlConfig, error) {
	c := &yamlConfig{}

	fileContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return yamlConfig{}, err
	}

	if err := yaml.Unmarshal(fileContents, c); err != nil {
		return yamlConfig{}, err
	}

	return *c, nil
}
