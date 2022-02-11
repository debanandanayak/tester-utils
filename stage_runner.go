package tester_utils

import (
	"fmt"
	"math/rand"
	"time"
)

// stageRunner is used to run multiple stages
type stageRunner struct {
	isQuiet         bool // Used for anti-cheat tests, where we only want Critical logs to be emitted
	isForFirstStage bool // We emit friendlier logs on test failures for the first stage
	stages          []Stage
}

// Stage represents a stage in a challenge.
//
// The Slug in a Stage should match that in the course's YAML definition.
type Stage struct {
	Slug                    string
	Title                   string
	TestFunc                func(stageHarness *StageHarness) error
	Timeout                 time.Duration
	ShouldRunPreviousStages bool
}

func (s Stage) CustomOrDefaultTimeout() time.Duration {
	if (s.Timeout == 0) || (s.Timeout == time.Duration(0)) {
		return 10 * time.Second
	} else {
		return s.Timeout
	}
}

func newStageRunner(stages []Stage, isForFirstStage bool) stageRunner {
	return stageRunner{
		isForFirstStage: isForFirstStage,
		stages:          stages,
	}
}

func newQuietStageRunner(stages []Stage) stageRunner {
	return stageRunner{isQuiet: true, stages: stages}
}

func (r stageRunner) getLoggerForStage(isDebug bool, stageNumber int) *Logger {
	if r.isQuiet {
		return getQuietLogger("")
	} else {
		return getLogger(isDebug, fmt.Sprintf("[stage-%d] ", stageNumber))
	}
}

func (r stageRunner) FirstStageSlug() string {
	return r.stages[0].Slug
}

func (r stageRunner) LastStageSlug() string {
	return r.stages[len(r.stages)-1].Slug
}

// Run runs all tests in a stageRunner
func (r stageRunner) Run(isDebug bool, executable *Executable) bool {
	for index, stage := range r.stages {
		stageNumber := index + 1

		stageHarness := StageHarness{
			Logger:     r.getLoggerForStage(isDebug, stageNumber),
			Executable: executable,
		}

		logger := stageHarness.Logger
		logger.Infof("Running tests for Stage #%d: %s", stageNumber, stage.Title)

		stageResultChannel := make(chan error, 1)
		go func() {
			err := stage.TestFunc(&stageHarness)
			stageResultChannel <- err
		}()

		var err error
		select {
		case stageErr := <-stageResultChannel:
			err = stageErr
		case <-time.After(stage.CustomOrDefaultTimeout()):
			err = fmt.Errorf("timed out, test exceeded %d seconds", int64(stage.CustomOrDefaultTimeout().Seconds()))
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

func (r stageRunner) StageBySlug(stageSlug string) Stage {
	for _, stage := range r.stages {
		if stage.Slug == stageSlug {
			return stage
		}
	}

	panic("Didn't find stage by slug " + stageSlug)
}

// ForStage returns a stageRunner with fewer stages
func (r stageRunner) ForStage(stageSlug string) stageRunner {
	currentStage := r.StageBySlug(stageSlug)

	if !currentStage.ShouldRunPreviousStages {
		return stageRunner{
			isForFirstStage: r.isForFirstStage,
			stages:          []Stage{currentStage},
		}
	}

	return r.Truncated(stageSlug)
}

// Truncated returns a stageRunner with fewer stages
func (r stageRunner) Truncated(stageSlug string) stageRunner {
	newStages := make([]Stage, 0)
	for _, stage := range r.stages {
		newStages = append(newStages, stage)
		if stage.Slug == stageSlug {
			return stageRunner{isForFirstStage: r.isForFirstStage, stages: newStages}
		}
	}

	panic(fmt.Sprintf("Stage slug %v not found. Stages: %v", stageSlug, r.stages))
}

// Randomized returns a stage runner that has stages randomized
func (r stageRunner) Randomized() stageRunner {
	return stageRunner{
		isForFirstStage: r.isForFirstStage,
		stages:          shuffleStages(r.stages),
	}
}

// Fuck you, go
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func shuffleStages(stages []Stage) []Stage {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Stage, len(stages))
	perm := r.Perm(len(stages))
	for i, randIndex := range perm {
		ret[i] = stages[randIndex]
	}
	return ret
}

func (r stageRunner) reportTestError(err error, isDebug bool, logger *Logger) {
	logger.Errorf("%s", err)

	if r.isForFirstStage {
		logger.Errorf("Test failed " +
			"(see README.md for instructions on how to pass this stage)")
	} else if isDebug {
		logger.Errorf("Test failed")
	} else {
		logger.Errorf("Test failed " +
			"(try setting 'debug: true' in your codecrafters.yml to see more details)")
	}
}
