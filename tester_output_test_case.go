package tester_utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TesterOutputTestCase struct {
	CodePath            string
	ExpectedExitCode    int
	StageName           string
	StdoutFixturePath   string
	NormalizeOutputFunc func([]byte) []byte
}

func TestTesterOutput(t *testing.T, testerDefinition TesterDefinition, testCases map[string]TesterOutputTestCase) {
	m := NewStdIOMocker()
	defer m.End()

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			m.Start()

			exitCode := runCLIStage(testerDefinition, testCase.StageName, testCase.CodePath)
			if !assert.Equal(t, testCase.ExpectedExitCode, exitCode) {
				failWithMockerOutput(t, m)
			}

			m.End()
			CompareOutputWithFixture(t, m.ReadStdout(), testCase.NormalizeOutputFunc, testCase.StdoutFixturePath)
		})
	}
}

func CompareOutputWithFixture(t *testing.T, testerOutput []byte, normalizeOutputFunc func([]byte) []byte, fixturePath string) {
	shouldRecordFixture := os.Getenv("CODECRAFTERS_RECORD_FIXTURES")

	if shouldRecordFixture == "true" {
		if err := os.MkdirAll(filepath.Dir(fixturePath), os.ModePerm); err != nil {
			panic(err)
		}

		if err := os.WriteFile(fixturePath, testerOutput, 0644); err != nil {
			panic(err)
		}

		return
	}

	fixtureContents, err := os.ReadFile(fixturePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Fixture file %s does not exist. To create a new one, use CODECRAFTERS_RECORD_FIXTURES=true", fixturePath)
			return
		}

		panic(err)
	}

	testerOutput = normalizeOutputFunc(testerOutput)
	fixtureContents = normalizeOutputFunc(fixtureContents)

	if bytes.Compare(testerOutput, fixtureContents) != 0 {
		diffExecutablePath, err := exec.LookPath("diff")
		if err != nil {
			panic(err)
		}

		diffExecutable := NewExecutable(diffExecutablePath)

		tmpFile, err := ioutil.TempFile("", "")
		if err != nil {
			panic(err)
		}

		if _, err = tmpFile.Write(testerOutput); err != nil {
			panic(err)
		}

		result, err := diffExecutable.Run(fixturePath, tmpFile.Name())
		if err != nil {
			panic(err)
		}

		os.Stdout.Write([]byte("\n\nDifferences detected:\n\n"))
		os.Stdout.Write(result.Stdout)
		os.Stdout.Write([]byte("\n\nRe-run this test with CODECRAFTERS_RECORD_FIXTURES=true to update fixtures\n\n"))
		t.FailNow()
	}
}

//func normalizeTesterOutput(testerOutput []byte) []byte {
//	re, _ := regexp.Compile("read tcp 127.0.0.1:\\d+->127.0.0.1:6379: read: connection reset by peer")
//	return re.ReplaceAll(testerOutput, []byte("read tcp 127.0.0.1:xxxxx+->127.0.0.1:6379: read: connection reset by peer"))
//}

func runCLIStage(testerDefinition TesterDefinition, slug string, relativePath string) (exitCode int) {
	// When a command is run with a different working directory, a relative path can cause problems.
	path, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	tester, err := NewTester(map[string]string{
		"CODECRAFTERS_CURRENT_STAGE_SLUG": slug,
		"CODECRAFTERS_SUBMISSION_DIR":     path,
	}, testerDefinition)

	if err != nil {
		fmt.Printf("%s", err)
		return 1
	}

	return tester.RunCLI()
}

func failWithMockerOutput(t *testing.T, m *IOMocker) {
	m.End()
	t.Error(fmt.Sprintf("stdout: \n%s\n\nstderr: \n%s", m.ReadStdout(), m.ReadStderr()))
	t.FailNow()
}
