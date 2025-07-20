package executor

import (
	"log/slog"
	"runtime-engine/internal/runners"
)

type Executor interface {
	Run(lang runners.Language, code []byte, log *slog.Logger) (runners.RunnerResult, error)
}
