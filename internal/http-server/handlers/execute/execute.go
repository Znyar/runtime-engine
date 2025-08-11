package execute

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"runtime-engine/internal/config"
	"runtime-engine/internal/runner"
)

type Request struct {
	Language             string `json:"language" validate:"required"`
	Version              string `json:"version" validate:"required"`
	Code                 string `json:"code" validate:"required"`
	Filename             string `json:"filename" validate:"required"`
	CompileTimeout       int64  `json:"compile_timeout"`
	RunTimeout           int64  `json:"run_timeout"`
	RunCPUTimeout        int64  `json:"run_cpu_timeout"`
	CompileCPUTimeout    int64  `json:"compile_cpu_timeout"`
	CompileMemoryLimitKB int64  `json:"compile_memory_limit_KB"`
	RunMemoryLimitKB     int64  `json:"run_memory_limit_KB"`
}

type Response struct {
	Error        string        `json:"error,omitempty"`
	RunnerResult runner.Result `json:"runner_result,omitempty"`
}

func New(log *slog.Logger, cfg *config.HttpServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{
			CompileTimeout:       cfg.CompileTimeout,
			RunTimeout:           cfg.RunTimeout,
			RunCPUTimeout:        cfg.RunCPUTimeout,
			CompileCPUTimeout:    cfg.CompileCPUTimeout,
			CompileMemoryLimitKB: cfg.CompileMemoryLimitKB,
			RunMemoryLimitKB:     cfg.RunMemoryLimitKB,
		}

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Error: "request body is empty",
			})

			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Error: fmt.Sprintf("failed to decode request: %s", err),
			})

			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{
				Error: fmt.Sprintf("invalid request: %s", err),
			})

			return
		}

		res, err := runner.Execute(runner.Request{
			Language:             req.Language,
			Version:              req.Version,
			Code:                 req.Code,
			Filename:             req.Filename,
			CompileTimeout:       req.CompileTimeout,
			RunTimeout:           req.RunTimeout,
			RunCPUTimeout:        req.RunCPUTimeout,
			CompileCPUTimeout:    req.CompileCPUTimeout,
			CompileMemoryLimitKB: req.CompileMemoryLimitKB,
			RunMemoryLimitKB:     req.RunMemoryLimitKB,
		}, log)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Error:        fmt.Sprintf("failed to process code: %s", err),
				RunnerResult: res,
			})

			return
		}

		render.JSON(w, r, Response{
			RunnerResult: res,
		})
	}
}
