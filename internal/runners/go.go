package runners

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type GoRunner struct{}

func (r *GoRunner) Execute(code []byte) (RunnerResult, error) {
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		return RunnerResult{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("failed to remove temp file: %s\n", name)
		}
	}(tmpFile.Name())

	if _, err := tmpFile.Write(code); err != nil {
		return RunnerResult{}, fmt.Errorf("failed to write code: %w", err)
	}
	err = tmpFile.Close()
	if err != nil {
		return RunnerResult{}, fmt.Errorf("failed to close file: %w", err)
	}

	compileStart := time.Now()
	execFile := tmpFile.Name() + ".exe"
	compileCmd := exec.Command("go", "build", "-o", execFile, tmpFile.Name())
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err = compileCmd.Run(); err != nil {
		return RunnerResult{
			Stderr:          compileStderr.Bytes(),
			ExitCode:        compileCmd.ProcessState.ExitCode(),
			CompilationTime: time.Since(compileStart),
		}, nil
	}
	compilationTime := time.Since(compileStart)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("failed to remove temp file: %s\n", name)
		}
	}(execFile)

	execStart := time.Now()
	execCmd := exec.Command(execFile)
	var execStdout, execStderr bytes.Buffer
	execCmd.Stdout = &execStdout
	execCmd.Stderr = &execStderr

	if err = execCmd.Run(); err != nil {
		return RunnerResult{
			Stderr:          execStderr.Bytes(),
			ExitCode:        execCmd.ProcessState.ExitCode(),
			ExecutionTime:   time.Since(execStart),
			CompilationTime: compilationTime,
		}, nil
	}

	return RunnerResult{
		Stdout:          execStdout.Bytes(),
		CompilationTime: compilationTime,
		ExecutionTime:   time.Since(execStart),
		ExitCode:        execCmd.ProcessState.ExitCode(),
	}, nil
}
