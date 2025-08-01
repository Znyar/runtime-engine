package runners

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"runtime-engine/pkg/logger/sl"
	"time"
)

type GoRunner struct{}

func (r *GoRunner) Execute(code []byte, log *slog.Logger) (RunnerResult, error) {
	const op = "runners.go.Execute"

	log = log.With(
		slog.String("op", op),
	)

	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		log.Error("failed to create temp file", sl.Err(err))
		return RunnerResult{}, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Error("failed to remove temp file", sl.Err(err))
		}
	}(tmpFile.Name())

	if _, err := tmpFile.Write(code); err != nil {
		log.Error("failed to create temp file", sl.Err(err))
		return RunnerResult{}, err
	}
	err = tmpFile.Close()
	if err != nil {
		log.Error("failed to close file:", err)
		return RunnerResult{}, err
	}

	log.Debug("compiling")
	compileStart := time.Now()
	execFileName := tmpFile.Name() + ".exe"
	compileCmd := exec.Command("/data/go/1.24.5/bin/go", "build", "-o", execFileName, tmpFile.Name())
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err = compileCmd.Run(); err != nil {
		return RunnerResult{
			StderrData:        compileStderr.Bytes(),
			StderrText:        string(compileStderr.Bytes()),
			ExitCode:          compileCmd.ProcessState.ExitCode(),
			CompilationTimeMs: time.Since(compileStart).Seconds() * 1000,
		}, err
	}
	compilationTime := time.Since(compileStart)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Error("failed to remove temp file", sl.Err(err))
		}
	}(execFileName)

	execStart := time.Now()
	execCmd := exec.Command(execFileName)
	var execStdout, execStderr bytes.Buffer
	execCmd.Stdout = &execStdout
	execCmd.Stderr = &execStderr

	log.Debug("running")
	if err = execCmd.Run(); err != nil {
		return RunnerResult{
			StderrData:        execStderr.Bytes(),
			StderrText:        string(execStderr.Bytes()),
			ExitCode:          execCmd.ProcessState.ExitCode(),
			ExecutionTimeMs:   time.Since(execStart).Seconds() * 1000,
			CompilationTimeMs: compilationTime.Seconds() * 1000,
		}, err
	}

	return RunnerResult{
		StdoutData:        execStdout.Bytes(),
		StdoutText:        string(execStdout.Bytes()),
		CompilationTimeMs: compilationTime.Seconds() * 1000,
		ExecutionTimeMs:   time.Since(execStart).Seconds() * 1000,
		ExitCode:          execCmd.ProcessState.ExitCode(),
	}, nil
}
