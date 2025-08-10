package executor

import (
	"log/slog"
	"runtime-engine/internal/runner"
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

func (e *DefaultExecutor) Run(lang string, version string, code []byte, log *slog.Logger) (runner.Result, error) {
	op := "executor.Run"

	e.semaphore.Acquire()
	defer e.semaphore.Release()

	result, err := runner.Execute(code, lang, version, log)
	if err != nil {
		log.Error(op, "runner execute error: %s", err)
		return result, err
	}

	return result, nil
}
