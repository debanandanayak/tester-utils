package tester_utils

import "fmt"

type Tester struct {
	antiCheatStageRunner stageRunner
	stageRunner          stageRunner
	context              testerContext
}

// NewTester creates a Tester based on the TesterDefinition provided
func NewTester(env map[string]string, definition TesterDefinition) (Tester, error) {
	context, err := getTesterContext(env, definition.ExecutableFileName)
	if err != nil {
		fmt.Println(err.Error())
		return Tester{}, err
	}

	isForFirstStage := context.currentStageSlug == definition.Stages[0].Slug

	return Tester{
		context:              context,
		stageRunner:          newStageRunner(definition.Stages, isForFirstStage),
		antiCheatStageRunner: newQuietStageRunner(definition.AntiCheatStages),
	}, nil
}

func (tester Tester) getQuietExecutable() *Executable {
	return NewExecutable(tester.context.executablePath)
}

func (tester Tester) getExecutable() *Executable {
	return NewVerboseExecutable(tester.context.executablePath, getLogger(true, "[your_program] ").Plainln)
}

// PrintDebugContext is to be run as early as possible after creating a Tester
func (tester Tester) PrintDebugContext() {
	if !tester.context.isDebug {
		return
	}

	tester.context.print()
	fmt.Println("")
}

// RunAntiCheatStages runs any anti-cheat stages specified in the TesterDefinition. Only critical logs are emitted. If
// the stages pass, the user won't see any visible output.
func (tester Tester) RunAntiCheatStages() bool {
	stageRunner := tester.antiCheatStageRunner
	return stageRunner.Run(false, tester.getQuietExecutable())
}

// RunStages runs all the stages upto the current stage the user is attempting. Returns true if all stages pass.
func (tester Tester) RunStages() bool {
	stageRunner := tester.stageRunner.ForStage(tester.context.currentStageSlug)
	return stageRunner.Run(tester.context.isDebug, tester.getExecutable())
}

func (tester Tester) IsFirstStage() bool {
	return tester.stageRunner.FirstStageSlug() == tester.context.currentStageSlug
}

func (tester Tester) IsLastStage() bool {
	return tester.stageRunner.LastStageSlug() == tester.context.currentStageSlug
}

// PrintFailureMessage is to be executed if RunStages fails. Don't execute if RunAntiCheatStages fails.
func (tester Tester) PrintFailureMessage() {
	fmt.Println("")
	fmt.Printf("View stage instructions at: %s\n", tester.context.coursePageUrl)
	fmt.Println("")
}

// PrintSuccessMessage is to be executed after RunStages and RunAntiCheatStages
func (tester Tester) PrintSuccessMessage() {
	fmt.Println("")
	fmt.Println("All tests ran successfully. Congrats!")

	if tester.IsLastStage() {
		fmt.Println("")
		fmt.Printf("Want to try another language or approach? Visit the course page: %s\n", tester.context.coursePageUrl)
		fmt.Println("")
	} else {
		fmt.Println("")
		fmt.Printf("View instructions for the next stage at: %s\n", tester.context.coursePageUrl)
		fmt.Println("")
	}
}
