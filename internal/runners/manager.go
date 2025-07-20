package runners

import (
	"fmt"
	"log/slog"
)

type Language string

const (
	Go Language = "go"
)

type Runner interface {
	Execute(code []byte, log *slog.Logger) (RunnerResult, error)
}

func GetRunner(lang Language) (Runner, error) {
	switch lang {
	case Go:
		return &GoRunner{}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}
