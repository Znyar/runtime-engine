package executor

import (
	"runtime-engine/internal/runners"
	"runtime-engine/pkg/logger"
	"runtime-engine/pkg/semaphore"
)

type DefaultExecutor struct {
	semaphore *semaphore.Semaphore
}

func NewDefaultExecutor(maxParallel int) *DefaultExecutor {
	return &DefaultExecutor{
		semaphore: semaphore.New(maxParallel),
	}
}

func (e *DefaultExecutor) Run(lang runners.Language, code []byte, log *logger.Logger) (runners.RunnerResult, error) {
	op := "executor.Run"

	e.semaphore.Acquire()
	defer e.semaphore.Release()

	r, err := runners.GetRunner(lang)
	if err != nil {
		log.Error(op, "get runner error: %s", err)
		return runners.RunnerResult{}, err
	}

	result, err := r.Execute(code)
	if err != nil {
		return result, err
	}

	return result, nil
}
