package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime-engine/pkg/logger/sl"
	"time"
)

type Request struct {
	Language             string
	Version              string
	Code                 string
	Filename             string
	CompileTimeout       int64
	RunTimeout           int64
	RunCPUTimeout        int64
	CompileCPUTimeout    int64
	CompileMemoryLimitKB int64
	RunMemoryLimitKB     int64
}

type Result struct {
	CompilationTimeMs float64   `json:"compilation_time_ms,omitempty"`
	ExecutionTimeMs   float64   `json:"execution_time_ms,omitempty"`
	StdoutText        string    `json:"stdout,omitempty"`
	StdoutData        []byte    `json:"stdout_data,omitempty"`
	StderrText        string    `json:"stderr_text,omitempty"`
	StderrData        []byte    `json:"stderr_data,omitempty"`
	ExitCode          int       `json:"exit_code"`
	Timestamp         time.Time `json:"timestamp"`
}

func Execute(req Request, log *slog.Logger) (Result, error) {
	tmpDir, err := os.MkdirTemp("", "exec-*")
	if err != nil {
		log.Error("failed to create temp dir", sl.Err(err))
		return Result{Timestamp: time.Now()}, err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Error("failed to remove temp dir", sl.Err(err))
			return
		}
	}(tmpDir)

	mainFile := filepath.Join(tmpDir, req.Filename)
	if err := os.WriteFile(mainFile, []byte(req.Code), 0644); err != nil {
		log.Error("failed to write temp file", sl.Err(err))
		return Result{Timestamp: time.Now()}, err
	}

	compileCtx, compileCancel := context.WithTimeout(context.Background(), time.Duration(req.CompileTimeout)*time.Second)
	defer compileCancel()

	log.Debug("compiling")
	compileStart := time.Now()
	execFileName := "main.exe"
	execFullPath := filepath.Join(tmpDir, execFileName)
	compileCmd := exec.CommandContext(compileCtx, "/bin/sh", "-c",
		fmt.Sprintf(`ulimit -v %d -t %d; exec bwrap --unshare-all --dev /dev --proc /proc \
			--ro-bind /data/%s/%s /data/%s/%s \
			--bind %s /job \
			--tmpfs /tmp --die-with-parent --chdir /job /data/%s/%s/bin/go build -o %s %s`,
			req.CompileMemoryLimitKB, req.CompileCPUTimeout,
			req.Language, req.Version, req.Language, req.Version,
			tmpDir,
			req.Language, req.Version,
			execFileName, req.Filename))
	compileCmd.Dir = tmpDir
	compileCmd.Env = append(os.Environ(),
		"GOROOT=/data/"+req.Language+"/"+req.Version,
		"GOCACHE=/tmp",
		"GO111MODULE=auto",
		"GOTMPDIR=/tmp",
	)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err = compileCmd.Run(); err != nil {
		if errors.Is(compileCtx.Err(), context.DeadlineExceeded) {
			log.Error("compilation timed out (wall time)")
			return Result{
				StderrText:        fmt.Sprintf("compilation timed out after %v s", req.CompileTimeout),
				ExitCode:          -1,
				CompilationTimeMs: time.Since(compileStart).Seconds() * 1000,
				Timestamp:         time.Now(),
			}, err
		}
		log.Error("compilation failed",
			sl.Err(err),
			slog.String("stderr", compileStderr.String()))
		return Result{
			StderrData:        compileStderr.Bytes(),
			StderrText:        compileStderr.String(),
			ExitCode:          compileCmd.ProcessState.ExitCode(),
			CompilationTimeMs: time.Since(compileStart).Seconds() * 1000,
			Timestamp:         time.Now(),
		}, err
	}
	compilationTime := time.Since(compileStart)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Error("failed to remove temp file", sl.Err(err))
			return
		}
	}(execFullPath)

	execCtx, execCancel := context.WithTimeout(context.Background(), time.Duration(req.RunTimeout)*time.Second)
	defer execCancel()

	log.Debug("running")
	execStart := time.Now()
	execCmd := exec.CommandContext(execCtx, "/bin/sh", "-c",
		fmt.Sprintf(`ulimit -v %d -t %d; exec bwrap --unshare-all --dev /dev --proc /proc \
			--bind %s /job \
			--tmpfs /tmp --die-with-parent --chdir /job /job/%s`,
			req.RunMemoryLimitKB, req.RunCPUTimeout,
			tmpDir,
			execFileName))
	execCmd.Dir = tmpDir
	execCmd.Env = append(os.Environ(),
		"GOCACHE=/tmp",
		"GO111MODULE=auto",
		"GOTMPDIR=/tmp",
	)
	var execStdout, execStderr bytes.Buffer
	execCmd.Stdout = &execStdout
	execCmd.Stderr = &execStderr

	if err = execCmd.Run(); err != nil {
		if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
			log.Error("execution timed out (wall time)")
			return Result{
				StderrText:        fmt.Sprintf("execution timed out after %v s", req.RunTimeout),
				ExitCode:          -1,
				ExecutionTimeMs:   time.Since(execStart).Seconds() * 1000,
				CompilationTimeMs: compilationTime.Seconds() * 1000,
				Timestamp:         time.Now(),
			}, err
		}
		log.Error("execution failed",
			sl.Err(err),
			slog.String("stderr", execStderr.String()))
		return Result{
			StderrData:        execStderr.Bytes(),
			StderrText:        execStderr.String(),
			ExitCode:          execCmd.ProcessState.ExitCode(),
			ExecutionTimeMs:   time.Since(execStart).Seconds() * 1000,
			CompilationTimeMs: compilationTime.Seconds() * 1000,
			Timestamp:         time.Now(),
		}, err
	}

	return Result{
		StdoutData:        execStdout.Bytes(),
		StdoutText:        execStdout.String(),
		CompilationTimeMs: compilationTime.Seconds() * 1000,
		ExecutionTimeMs:   time.Since(execStart).Seconds() * 1000,
		ExitCode:          execCmd.ProcessState.ExitCode(),
		Timestamp:         time.Now(),
	}, nil
}
