package pack

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"runtime-engine/internal/langs"
)

type Request struct {
	Language string `json:"language" validate:"required"`
	Version  string `json:"version" validate:"required"`
}

type Response struct {
	Error string `json:"error"`
}

func New(log *slog.Logger, m langs.LangManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = m.Pack(req.Language, req.Version)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Error: fmt.Sprintf("failed to package language: %s", err),
			})

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
