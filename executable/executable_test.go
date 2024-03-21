package executable

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	err := NewExecutable("/blah").Start()
	assertErrorContains(t, err, "not found")
	assertErrorContains(t, err, "/blah")

	// Permissions are not preserved across remote git repos.
	_ = removeFileExecutablePermission("./test_helpers/not_executable.sh")

	err = NewExecutable("./test_helpers/not_executable.sh").Start()
	assertErrorContains(t, err, "not an executable file")
	assertErrorContains(t, err, "not_executable.sh")

	err = NewExecutable("./test_helpers/haskell").Start()
	assertErrorContains(t, err, "not an executable file")
	assertErrorContains(t, err, "haskell")

	err = NewExecutable("./test_helpers/stdout_echo.sh").Start()
	assert.NoError(t, err)
}

func assertErrorContains(t *testing.T, err error, expectedMsg string) {
	assert.Contains(t, err.Error(), expectedMsg)
}

func removeFileExecutablePermission(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	currentMode := fileInfo.Mode()

	// Clear the executable bits for user, group, and others
	newMode := currentMode &^ (0111)

	// Update the file mode
	err = os.Chmod(filePath, newMode)
	if err != nil {
		return err
	}
	return nil
}

func TestRun(t *testing.T) {
	e := NewExecutable("./test_helpers/stdout_echo.sh")
	result, err := e.Run("hey")
	assert.NoError(t, err)
	assert.Equal(t, "hey\n", string(result.Stdout))
}

func TestOutputCapture(t *testing.T) {
	// Stdout capture
	e := NewExecutable("./test_helpers/stdout_echo.sh")
	result, err := e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, "hey\n", string(result.Stdout))
	assert.Equal(t, "", string(result.Stderr))

	// Stderr capture
	e = NewExecutable("./test_helpers/stderr_echo.sh")
	result, err = e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, "", string(result.Stdout))
	assert.Equal(t, "hey\n", string(result.Stderr))
}

func TestLargeOutputCapture(t *testing.T) {
	e := NewExecutable("./test_helpers/large_echo.sh")
	result, err := e.Run("hey")

	assert.NoError(t, err)
	assert.Equal(t, 1024*1024, len(result.Stdout))
	assert.Equal(t, "blah\n", string(result.Stderr))
}

func TestExitCode(t *testing.T) {
	e := NewExecutable("./test_helpers/exit_with.sh")

	result, _ := e.Run("0")
	assert.Equal(t, 0, result.ExitCode)

	result, _ = e.Run("1")
	assert.Equal(t, 1, result.ExitCode)

	result, _ = e.Run("2")
	assert.Equal(t, 2, result.ExitCode)
}

func TestExecutableStartNotAllowedIfInProgress(t *testing.T) {
	e := NewExecutable("./test_helpers/sleep_for.sh")

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
	e := NewExecutable("./test_helpers/stdout_echo.sh")

	result, _ := e.Run("1")
	assert.Equal(t, "1\n", string(result.Stdout))

	result, _ = e.Run("2")
	assert.Equal(t, "2\n", string(result.Stdout))
}

func TestHasExited(t *testing.T) {
	e := NewExecutable("./test_helpers/sleep_for.sh")

	e.Start("0.1")
	assert.False(t, e.HasExited(), "Expected to not have exited")

	time.Sleep(150 * time.Millisecond)
	assert.True(t, e.HasExited(), "Expected to have exited")
}

func TestStdin(t *testing.T) {
	e := NewExecutable("/usr/bin/grep")

	err := e.Start("cat")
	if err != nil {
		return
	}
	assert.False(t, e.HasExited(), "Expected to not have exited")

	e.StdinPipe.Write([]byte("has cat"))
	assert.False(t, e.HasExited(), "Expected to not have exited")

	e.StdinPipe.Close()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, e.HasExited(), "Expected to have exited")
}

func TestRunWithStdin(t *testing.T) {
	e := NewExecutable("/usr/bin/grep")

	result, err := e.RunWithStdin([]byte("has cat"), "cat")
	assert.NoError(t, err)

	assert.Equal(t, result.ExitCode, 0)

	result, err = e.RunWithStdin([]byte("only dog"), "cat")
	assert.NoError(t, err)

	assert.Equal(t, result.ExitCode, 1)
}

// Rogue == doesn't respond to SIGTERM
func TestTerminatesRoguePrograms(t *testing.T) {
	e := NewExecutable("/bin/bash")

	err := e.Start("-c", "trap '' SIGTERM SIGINT; sleep 60")
	if err != nil {
		return
	}
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = e.Kill()
	assert.EqualError(t, err, "program failed to exit in 2 seconds after receiving sigterm")

	// Starting again shouldn't throw an error
	err = e.Start("-c", "trap '' SIGTERM SIGINT; sleep 60")
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = e.Kill()
	assert.EqualError(t, err, "program failed to exit in 2 seconds after receiving sigterm")
}
