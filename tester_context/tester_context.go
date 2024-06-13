package tester_context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/codecrafters-io/tester-utils/tester_definition"
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
	ExecutablePath               string
	IsDebug                      bool
	TestCases                    []TesterContextTestCase
	ShouldSkipAntiCheatTestCases bool
}

type yamlConfig struct {
	Debug bool `yaml:"debug"`
}

func (c TesterContext) Print() {
	fmt.Println("Debug =", c.IsDebug)
}

// GetContext parses flags and returns a Context object
func GetTesterContext(env map[string]string, definition tester_definition.TesterDefinition) (TesterContext, error) {
	submissionDir, ok := env["CODECRAFTERS_SUBMISSION_DIR"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_SUBMISSION_DIR env var not found")
	}

	testCasesJson, ok := env["CODECRAFTERS_TEST_CASES_JSON"]
	if !ok {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES_JSON env var not found")
	}

	testCases := []TesterContextTestCase{}
	if err := json.Unmarshal([]byte(testCasesJson), &testCases); err != nil {
		return TesterContext{}, fmt.Errorf("failed to parse CODECRAFTERS_TEST_CASES_JSON: %s", err)
	}

	var shouldSkipAntiCheatTestCases = false

	skipAntiCheatValue, ok := env["CODECRAFTERS_SKIP_ANTI_CHEAT"]
	if ok && skipAntiCheatValue == "true" {
		shouldSkipAntiCheatTestCases = true
	}

	for _, testCase := range testCases {
		if testCase.Slug == "" {
			return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES_JSON contains a test case with an empty slug")
		}

		if testCase.TesterLogPrefix == "" {
			return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES_JSON contains a test case with an empty tester_log_prefix")
		}

		if testCase.Title == "" {
			return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES_JSON contains a test case with an empty title")
		}
	}

	var executablePath string
	newExecutablePath := path.Join(submissionDir, definition.ExecutableFileName)
	legacyExecutablePath := path.Join(submissionDir, definition.LegacyExecutableFileName)

	_, err := os.Stat(legacyExecutablePath)
	if errors.Is(err, os.ErrNotExist) { // legacy executable does not exist
		executablePath = newExecutablePath
	}
	if err == nil { // legacy executable exists
		_, err := os.Stat(newExecutablePath)
		if errors.Is(err, os.ErrNotExist) { // your_program.sh does not exist
			executablePath = legacyExecutablePath
		}

		if err == nil { // your_program.sh exists
			executablePath = newExecutablePath
		}
	}

	configPath := path.Join(submissionDir, "codecrafters.yml")

	yamlConfig, err := readFromYAML(configPath)
	if err != nil {
		return TesterContext{}, err
	}

	if len(testCases) == 0 {
		return TesterContext{}, fmt.Errorf("CODECRAFTERS_TEST_CASES is empty")
	}

	// TODO: test if executable exists?

	return TesterContext{
		ExecutablePath:               executablePath,
		IsDebug:                      yamlConfig.Debug,
		TestCases:                    testCases,
		ShouldSkipAntiCheatTestCases: shouldSkipAntiCheatTestCases,
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
