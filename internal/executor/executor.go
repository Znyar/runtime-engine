package executor

import (
	"log/slog"
	"runtime-engine/internal/runner"
)

type Executor interface {
	Run(lang string, version string, code []byte, log *slog.Logger) (runner.Result, error)
}
