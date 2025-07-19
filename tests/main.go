package main

import (
	"os"
	"runtime-engine/internal/executor"
	"runtime-engine/internal/runners"
	"runtime-engine/pkg/logger"
	"sync"
	"time"
)

func main() {
	var log = logger.New(os.Stdout, logger.LevelDebug, true)

	if len(os.Args) < 3 {
		log.Info("Usage: runner <language> <file>")
		os.Exit(1)
	}

	lang := os.Args[1]
	filePath := os.Args[2]

	code, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("Error reading file: %v", err)
		os.Exit(1)
	}

	exec := executor.NewCachedExecutor(5*time.Second, 10)

	var wg sync.WaitGroup
	jobsCount := 50

	for i := 0; i < jobsCount; i++ {
		wg.Add(1)
		go func(jobId int) {
			defer wg.Done()
			startNewJob(exec, runners.Language(lang), code, jobId, log)
		}(i + 1)
		time.Sleep(1000 * time.Millisecond)
	}

	wg.Wait()
	log.Info("All jobs completed")
}

func startNewJob(exec *executor.CachedExecutor, lang runners.Language, code []byte, jobId int, log *logger.Logger) {
	log.Info("Job %d started", jobId)

	res, err := exec.Run(lang, code, log)

	if err != nil {
		log.Error("Execution error: %v", err)
		return
	}

	log.Info("Job %d result:\nStdout: %s\nStderr: %s\nExit code: %d\nExecution time: %s\nCompilation time: %s", jobId, res.StdoutText, res.StderrText, res.ExitCode, res.ExecutionTime, res.CompilationTime)
}
