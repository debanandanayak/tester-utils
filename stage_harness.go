package tester_utils

import (
	"github.com/codecrafters-io/tester-utils/logger"
)

// StageHarness is passed to your Stage's TestFunc.
//
// If the program is a long-lived program that must be alive during the duration of the test (like a Redis server),
// do something like this at the start of your test function:
//
//	if err := stageHarness.Executable.Start(); err != nil {
//	   return err
//	}
//	stageHarness.RegisterTeardownFunc(func() { stageHarness.Executable.Kill() })
//
// If the program is a script that must be executed and then checked for output (like a Git command), use it like this:
//
//	result, err := stageHarness.Executable.Run("cat-file", "-p", "sha")
//	if err != nil {
//	    return err
//	 }
type StageHarness struct {
	// Logger is to be used for all logs generated from the test function.
	Logger *logger.Logger

	// Executable is the program to be tested.
	Executable *Executable

	// teardownFuncs are run once the error has been reported to the user
	teardownFuncs []func()
}

func (s *StageHarness) RegisterTeardownFunc(teardownFunc func()) {
	s.teardownFuncs = append(s.teardownFuncs, teardownFunc)
}

func (s StageHarness) RunTeardownFuncs() {
	for _, teardownFunc := range s.teardownFuncs {
		teardownFunc()
	}
}
