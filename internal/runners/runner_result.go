package runners

import (
	"time"
)

type RunnerResult struct {
	CompilationTime time.Duration
	ExecutionTime   time.Duration
	Stdout          []byte
	Stderr          []byte
	ExitCode        int
}
