package tester_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	err := NewExecutable("/blah").Start()
	assertErrorContains(t, err, "no such file")
	assertErrorContains(t, err, "/blah")

	err = NewExecutable("./test_helpers/executable_test/stdout_echo.sh").Start()
	assert.NoError(t, err)
}

func assertErrorContains(t *testing.T, err error, expectedMsg string) {
	assert.Contains(t, err.Error(), expectedMsg)
}

func TestRun(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/stdout_echo.sh")
	result, err := e.Run("hey")
	assert.NoError(t, err)
	assert.Equal(t, "hey\n", string(result.Stdout))
}

func TestOutputCapture(t *testing.T) {
	// Stdout capture
	e := NewExecutable("./test_helpers/executable_test/stdout_echo.sh")
	result, err := e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, "hey\n", string(result.Stdout))
	assert.Equal(t, "", string(result.Stderr))

	// Stderr capture
	e = NewExecutable("./test_helpers/executable_test/stderr_echo.sh")
	result, err = e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, "", string(result.Stdout))
	assert.Equal(t, "hey\n", string(result.Stderr))
}

func TestLargeOutputCapture(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/large_echo.sh")
	result, err := e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, 1024*1024, len(result.Stdout))
	assert.Equal(t, "blah\n", string(result.Stderr))
}

func TestExitCode(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/exit_with.sh")

	result, _ := e.Run("0")
	assert.Equal(t, 0, result.ExitCode)

	result, _ = e.Run("1")
	assert.Equal(t, 1, result.ExitCode)

	result, _ = e.Run("2")
	assert.Equal(t, 2, result.ExitCode)
}

func TestExecutableStartNotAllowedIfInProgress(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/sleep_for.sh")

	// Run once
	err := e.Start("0.01")
	assert.NoError(t, err)

	// Starting again when in progress should throw an error
	err = e.Start("0.01")
	assertErrorContains(t, err, "process already in progress")

	// Running again when in progress should throw an error
	_, err = e.Run("0.01")
	assertErrorContains(t, err, "process already in progress")

	e.Wait()

	// Running again once finished should be fine
	err = e.Start("0.01")
	assert.NoError(t, err)
}

func TestSuccessiveExecutions(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/stdout_echo.sh")

	result, _ := e.Run("1")
	assert.Equal(t, "1\n", string(result.Stdout))

	result, _ = e.Run("2")
	assert.Equal(t, "2\n", string(result.Stdout))
}

func TestHasExited(t *testing.T) {
	e := NewExecutable("./test_helpers/executable_test/sleep_for.sh")

	e.Start("0.1")
	assert.False(t, e.HasExited(), "Expected to not have exited")

	time.Sleep(150 * time.Millisecond)
	assert.True(t, e.HasExited(), "Expected to have exited")
}

func TestStdin(t *testing.T) {
	e := NewExecutable("grep")

	e.Start("cat")
	assert.False(t, e.HasExited(), "Expected to not have exited")

	e.StdinPipe.Write([]byte("has cat"))
	assert.False(t, e.HasExited(), "Expected to not have exited")

	e.StdinPipe.Close()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, e.HasExited(), "Expected to have exited")
}

func TestRunWithStdin(t *testing.T) {
	e := NewExecutable("grep")

	result, err := e.RunWithStdin([]byte("has cat"), "cat")
	assert.NoError(t, err)

	assert.Equal(t, result.ExitCode, 0)

	result, err = e.RunWithStdin([]byte("only dog"), "cat")
	assert.NoError(t, err)

	assert.Equal(t, result.ExitCode, 1)
}

// Rogue == doesn't respond to SIGTERM
func TestTerminatesRoguePrograms(t *testing.T) {
	e := NewExecutable("bash")

	err := e.Start("-c", "trap '' SIGTERM SIGINT; sleep 60")
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = e.Kill()
	assert.EqualError(t, err, "program failed to exit in 2 seconds after receiving sigterm")
}
