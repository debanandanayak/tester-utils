package tester_utils

import (
	"fmt"

	"math/rand"
	"time"
)

// stageRunner is used to run multiple stages
type stageRunner struct {
	isQuiet bool // Used for anti-cheat tests, where we only want Critical logs to be emitted
	stages  []Stage
}

// Stage represents a stage in a challenge.
//
// The Slug in a Stage should match that in the course's YAML definition.
type Stage struct {
	Slug     string
	Title    string
	TestFunc func(stageHarness StageHarness) error
}

func newStageRunner(stages []Stage) stageRunner {
	return stageRunner{stages: stages}
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
		logger.Infof("Running test: %s", stage.Title)

		stageResultChannel := make(chan error, 1)
		go func() {
			err := stage.TestFunc(stageHarness)
			stageResultChannel <- err
		}()

		var err error
		select {
		case stageErr := <-stageResultChannel:
			err = stageErr
		case <-time.After(10 * time.Second):
			err = fmt.Errorf("timed out, test exceeded 10 seconds")
		}

		if err != nil {
			reportTestError(err, isDebug, logger)
			return false
		}

		logger.Successf("Test passed.")
	}

	return true
}

// Truncated returns a stageRunner with fewer stages
func (r stageRunner) Truncated(stageSlug string) stageRunner {
	newStages := make([]Stage, 0)
	for _, stage := range r.stages {
		newStages = append(newStages, stage)
		if stage.Slug == stageSlug {
			return stageRunner{stages: newStages}
		}
	}

	panic(fmt.Sprintf("Stage slug %v not found. Stages: %v", stageSlug, r.stages))
}

// Randomized returns a stage runner that has stages randomized
func (r stageRunner) Randomized() stageRunner {
	return stageRunner{
		stages: shuffleStages(r.stages),
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

func reportTestError(err error, isDebug bool, logger *Logger) {
	logger.Errorf("%s", err)
	if isDebug {
		logger.Errorf("Test failed")
	} else {
		logger.Errorf("Test failed " +
			"(try setting 'debug: true' in your codecrafters.yml to see more details)")
	}
}
