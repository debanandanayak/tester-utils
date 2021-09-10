package tester_utils

import "fmt"

// StageHarness is passed to your Stage's TestFunc.
//
// If the program is a long-lived program that must be alive during the duration of the test (like a Redis server),
// do something like this at the start of your test function:
//
//  if err := stageHarness.Executable.Run(); err != nil {
//     return err
//  }
//  defer stageHarness.Executable.Kill()
//
// If the program is a script that must be executed and then checked for output (like a Git command), use it like this:
//
//  result, err := executable.Run("cat-file", "-p", "sha")
//  if err != nil {
//      return err
//   }
type StageHarness struct {
	// Logger is to be used for all logs generated from the test function.
	Logger *Logger

	// Executable is the program to be tested.
	Executable *Executable

	// teardownFuncs are run once the error has been reported to the user
	teardownFuncs []func()
}

func (s StageHarness) RegisterTeardownFunc(teardownFunc func()) {
	fmt.Println("registering teardown func")
	s.teardownFuncs = append(s.teardownFuncs, teardownFunc)
}

func (s StageHarness) RunTeardownFuncs() {
	fmt.Println("running teardown funcs...")
	fmt.Println(len(s.teardownFuncs))
	for _, teardownFunc := range s.teardownFuncs {
		fmt.Println("running a teardown func")
		teardownFunc()
		fmt.Println("ran a teardown func")
	}
}
