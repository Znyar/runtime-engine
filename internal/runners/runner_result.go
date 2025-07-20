package runners

type RunnerResult struct {
	CompilationTimeMs float64 `json:"compilation_time_ms,omitempty"`
	ExecutionTimeMs   float64 `json:"execution_time_ms,omitempty"`
	StdoutText        string  `json:"stdout,omitempty"`
	StdoutData        []byte  `json:"stdout_data,omitempty"`
	StderrText        string  `json:"stderr_text,omitempty"`
	StderrData        []byte  `json:"stderr_data,omitempty"`
	ExitCode          int     `json:"exit_code,omitempty"`
}
