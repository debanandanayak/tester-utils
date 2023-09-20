package testing

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/codecrafters-io/tester-utils/stdio_mocker"
	"github.com/stretchr/testify/assert"
)

type TesterOutputTestCase struct {
	// CodePath is the path to the code that'll be tested.
	CodePath string

	// ExpectedExitCode is the exit code that we expect the tester to return.
	ExpectedExitCode int

	// UntilStageSlug is the slug of the stage that we want to test until. Either this or StageSlug must be provided.
	UntilStageSlug string

	// StageSlugs is the list of stages that we want to test. Either this or UntilStageSlug must be provided.
	StageSlugs []string

	// StdoutFixturePath is the path to the fixture file that contains the expected stdout output.
	StdoutFixturePath string

	// NormalizeOutputFunc is a function that normalizes the tester's output. This is useful for removing things like timestamps.
	NormalizeOutputFunc func([]byte) []byte
}

func buildTestCasesJson(slugs []string) string {
	testCases := []map[string]string{}

	for index, slug := range slugs {
		testCases = append(testCases, map[string]string{
			"slug":              slug,
			"tester_log_prefix": fmt.Sprintf("stage-%d", index+1),
			"title":             fmt.Sprintf("Stage #%d: %s", index+1, slug),
		})
	}

	testCasesJson, _ := json.Marshal(testCases)
	return string(testCasesJson)
}

func buildTestCasesJsonUntilStageSlug(untilStageSlug string, testerDefinition tester_utils.TesterDefinition) string {
	stageSlugs := []string{}

	for _, testCase := range testerDefinition.TestCases {
		stageSlugs = append(stageSlugs, testCase.Slug)

		if testCase.Slug == untilStageSlug {
			break
		}
	}

	return buildTestCasesJson(stageSlugs)
}

func TestTesterOutput(t *testing.T, testerDefinition tester_utils.TesterDefinition, testCases map[string]TesterOutputTestCase) {
	m := stdio_mocker.NewStdIOMocker()
	defer m.End()

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			m.Start()

			if testCase.UntilStageSlug != "" && testCase.StageSlugs != nil && len(testCase.StageSlugs) > 0 {
				t.Fatal("Either UntilStageSlug or StageSlugs must be provided, not both")
			}

			if testCase.UntilStageSlug == "" && (testCase.StageSlugs == nil || len(testCase.StageSlugs) == 0) {
				t.Fatal("Either UntilStageSlug or StageSlugs must be provided")
			}

			var testCasesJson string

			if testCase.UntilStageSlug != "" {
				testCasesJson = buildTestCasesJsonUntilStageSlug(testCase.UntilStageSlug, testerDefinition)
			} else {
				testCasesJson = buildTestCasesJson(testCase.StageSlugs)
			}

			exitCode := runCLIStage(testerDefinition, testCasesJson, testCase.CodePath)
			if !assert.Equal(t, testCase.ExpectedExitCode, exitCode) {
				failWithMockerOutput(t, m)
			}

			m.End()
			CompareOutputWithFixture(t, m.ReadStdout(), testCase.NormalizeOutputFunc, testCase.StdoutFixturePath)
		})
	}
}

//func normalizeTesterOutput(testerOutput []byte) []byte {
//	re, _ := regexp.Compile("read tcp 127.0.0.1:\\d+->127.0.0.1:6379: read: connection reset by peer")
//	return re.ReplaceAll(testerOutput, []byte("read tcp 127.0.0.1:xxxxx+->127.0.0.1:6379: read: connection reset by peer"))
//}

func runCLIStage(testerDefinition tester_utils.TesterDefinition, testCasesJson string, relativePath string) (exitCode int) {
	// When a command is run with a different working directory, a relative path can cause problems.
	path, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	tester, err := tester_utils.NewTester(map[string]string{
		"CODECRAFTERS_TEST_CASES_JSON": testCasesJson,
		"CODECRAFTERS_SUBMISSION_DIR":  path,
	}, testerDefinition)

	if err != nil {
		fmt.Printf("%s", err)
		return 1
	}

	return tester.RunCLI()
}

func failWithMockerOutput(t *testing.T, m *stdio_mocker.IOMocker) {
	m.End()
	t.Error(fmt.Sprintf("stdout: \n%s\n\nstderr: \n%s", m.ReadStdout(), m.ReadStderr()))
	t.FailNow()
}
