package execute

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"runtime-engine/internal/executor"
	"runtime-engine/internal/runner"
)

type Request struct {
	Language string `json:"language" validate:"required"`
	Version  string `json:"version" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type Response struct {
	Error        string        `json:"error,omitempty"`
	RunnerResult runner.Result `json:"runner_result,omitempty"`
}

func New(log *slog.Logger, e executor.Executor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.execute.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

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

		res, err := e.Run(req.Language, req.Version, []byte(req.Code), log)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Error: fmt.Sprintf("failed to process code: %s", err),
			})

			return
		}

		render.JSON(w, r, Response{
			RunnerResult: res,
		})

	}
}
