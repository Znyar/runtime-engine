package executor

import (
	"log/slog"
	"runtime-engine/internal/runners"
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

func (e *DefaultExecutor) Run(lang runners.Language, code []byte, log *slog.Logger) (runners.RunnerResult, error) {
	op := "executor.Run"

	e.semaphore.Acquire()
	defer e.semaphore.Release()

	r, err := runners.GetRunner(lang)
	if err != nil {
		log.Error(op, "get runner error: %s", err)
		return runners.RunnerResult{}, err
	}

	result, err := r.Execute(code, log)
	if err != nil {
		log.Error(op, "runner execute error: %s", err)
		return result, err
	}

	return result, nil
}
