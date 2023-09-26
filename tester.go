package tester_utils

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/tester_context"
)

type Tester struct {
	context    tester_context.TesterContext
	definition TesterDefinition
}

// NewTester creates a Tester based on the TesterDefinition provided
func NewTester(env map[string]string, definition TesterDefinition) (Tester, error) {
	context, err := tester_context.GetTesterContext(env, definition.ExecutableFileName)
	if err != nil {
		fmt.Println(err.Error())
		return Tester{}, err
	}

	tester := Tester{
		context:    context,
		definition: definition,
	}

	if err := tester.validateContext(); err != nil {
		return Tester{}, err
	}

	return tester, nil
}

func (tester Tester) PrintDebugContext() {
	// PrintDebugContext is to be run as early as possible after creating a Tester func (tester Tester) PrintDebugContext() {
	if !tester.context.IsDebug {
		return
	}

	tester.context.Print()
	fmt.Println("")
}

// RunAntiCheatStages runs any anti-cheat stages specified in the TesterDefinition. Only critical logs are emitted. If
// the stages pass, the user won't see any visible output.
func (tester Tester) RunAntiCheatStages() bool {
	return tester.getAntiCheatRunner().Run(false, tester.getQuietExecutable())
}

// RunStages runs all the stages upto the current stage the user is attempting. Returns true if all stages pass.
func (tester Tester) RunStages() bool {
	return tester.getRunner().Run(tester.context.IsDebug, tester.getExecutable())
}

// RunCLI executes the tester based on user-provided env vars
func (tester Tester) RunCLI() int {
	tester.PrintDebugContext()

	// Validate context?

	if !tester.RunStages() {
		return 1
	}

	if os.Getenv("CODECRAFTERS_DISABLE_ANTI_CHEAT_TEST_CASES") != "true" {
		if !tester.RunAntiCheatStages() {
			return 1
		}
	}

	return 0
}

func (tester Tester) getRunner() testRunner {
	steps := []testRunnerStep{}

	for _, testerContextTestCase := range tester.context.TestCases {
		definitionTestCase := tester.definition.TestCaseBySlug(testerContextTestCase.Slug)

		steps = append(steps, testRunnerStep{
			testCase:        definitionTestCase,
			testerLogPrefix: testerContextTestCase.TesterLogPrefix,
			title:           testerContextTestCase.Title,
		})
	}

	return testRunner{
		isQuiet: false,
		steps:   steps,
	}
}

func (tester Tester) getAntiCheatRunner() testRunner {
	steps := []testRunnerStep{}

	for index, testCase := range tester.definition.AntiCheatTestCases {
		steps = append(steps, testRunnerStep{
			testCase:        testCase,
			testerLogPrefix: fmt.Sprintf("ac-%d", index+1),
			title:           fmt.Sprintf("AC%d", index+1),
		})
	}
	return testRunner{
		isQuiet: true, // We only want Critical logs to be emitted for anti-cheat tests
		steps:   steps,
	}
}

func (tester Tester) getQuietExecutable() *executable.Executable {
	return executable.NewExecutable(tester.context.ExecutablePath)
}

func (tester Tester) getExecutable() *executable.Executable {
	return executable.NewVerboseExecutable(tester.context.ExecutablePath, logger.GetLogger(true, "[your_program] ").Plainln)
}

func (tester Tester) validateContext() error {
	for _, testerContextTestCase := range tester.context.TestCases {
		testerDefinitionTestCase := tester.definition.TestCaseBySlug(testerContextTestCase.Slug)

		if testerDefinitionTestCase.Slug != testerContextTestCase.Slug {
			return fmt.Errorf("tester context does not have test case with slug %s", testerContextTestCase.Slug)
		}
	}

	return nil
}
