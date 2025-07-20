package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime-engine/internal/executor"
	"runtime-engine/internal/runners"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Enable JSON output")
	flag.Parse()

	var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: runner [flags] <language> <file>")
		os.Exit(1)
	}

	lang := runners.Language(flag.Arg(0))
	filePath := flag.Arg(1)

	code, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	exec := executor.NewDefaultExecutor(10)

	result, err := exec.Run(lang, code, log)
	if err != nil {
		fmt.Printf("Execution error: %v", err)
		if *jsonOutput {
			err := json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
				"error": err.Error(),
			})
			if err != nil {
				fmt.Printf("Error while encoding JSON: %v\n", err)
			}
		}
		os.Exit(1)
	}

	if *jsonOutput {
		err := json.NewEncoder(os.Stdout).Encode(result)
		if err != nil {
			fmt.Printf("Error while encoding JSON: %v\n", err)
		}
	} else {
		_, err = fmt.Fprintf(os.Stdout, "Result:\nStdout: %s\nStderr: %s\nExit code: %d\nExecution time: %s\nCompilation time: %s", result.StdoutText, result.StderrText, result.ExitCode, result.ExecutionTimeMs, result.CompilationTimeMs)
		if err != nil {
			fmt.Printf("Error while printing result: %v\n", err)
		}
	}
}
