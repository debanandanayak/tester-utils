package tester_utils

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type testRunnerStep struct {
	// testCase is the test case that'll be run against the user's code.
	testCase TestCase

	// testerLogPrefix is the prefix that'll be used for all logs emitted by the tester. Example: "stage-1"
	testerLogPrefix string

	// title is the title of the test case. Example: "Stage #1: Bind to a port"
	title string
}

// testRunner is used to run multiple tests
type testRunner struct {
	isQuiet bool // Used for anti-cheat tests, where we only want Critical logs to be emitted
	steps   []testRunnerStep
}

func newTestRunner(steps []testRunnerStep) testRunner {
	return testRunner{
		steps: steps,
	}
}

func newQuietStageRunner(steps []testRunnerStep) testRunner {
	return testRunner{isQuiet: true, steps: steps}
}

func (r testRunner) getLoggerForStep(isDebug bool, step testRunnerStep) *logger.Logger {
	if r.isQuiet {
		return logger.GetQuietLogger("")
	} else {
		return logger.GetLogger(isDebug, fmt.Sprintf("[%s] ", step.testerLogPrefix))
	}
}

// Run runs all tests in a stageRunner
func (r testRunner) Run(isDebug bool, executable *executable.Executable) bool {
	for index, step := range r.steps {
		if index != 0 {
			fmt.Println("")
		}

		stageHarness := StageHarness{
			Logger:     r.getLoggerForStep(isDebug, step),
			Executable: executable,
		}

		logger := stageHarness.Logger
		logger.Infof("Running tests for %s:", step.title)

		stepResultChannel := make(chan error, 1)
		go func() {
			err := step.testCase.TestFunc(&stageHarness)
			stepResultChannel <- err
		}()

		timeout := step.testCase.CustomOrDefaultTimeout()

		var err error
		select {
		case stageErr := <-stepResultChannel:
			err = stageErr
		case <-time.After(timeout):
			err = fmt.Errorf("timed out, test exceeded %d seconds", int64(timeout.Seconds()))
		}

		if err != nil {
			r.reportTestError(err, isDebug, logger)
		} else {
			logger.Successf("Test passed.")
		}

		stageHarness.RunTeardownFuncs()

		if err != nil {
			return false
		}
	}

	return true
}

// Fuck you, go
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (r testRunner) reportTestError(err error, isDebug bool, logger *logger.Logger) {
	logger.Errorf("%s", err)

	if isDebug {
		logger.Errorf("Test failed")
	} else {
		logger.Errorf("Test failed " +
			"(try setting 'debug: true' in your codecrafters.yml to see more details)")
	}
}
