package runners

import (
	"time"
)

type RunnerResult struct {
	CompilationTime time.Duration `json:"compilation_time_ms"`
	ExecutionTime   time.Duration `json:"execution_time_ms"`
	StdoutText      string        `json:"stdout,omitempty"`
	StdoutData      []byte        `json:"stdout_data,omitempty"`
	StderrText      string        `json:"stderr_text,omitempty"`
	StderrData      []byte        `json:"stderr_data,omitempty"`
	ExitCode        int           `json:"exit_code"`
}
