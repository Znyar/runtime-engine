package executor

import (
	"runtime-engine/internal/runners"
	"runtime-engine/pkg/logger"
)

type Executor interface {
	Run(lang runners.Language, code []byte, log *logger.Logger) (runners.RunnerResult, error)
}
