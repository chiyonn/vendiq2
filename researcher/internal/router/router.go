package router

import (
	"net/http"

	"github.com/chiyonn/vendiq2/researcher/internal/handler"
	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(LoggingMiddleware)
	r.Get("/health", handler.Health)
	return r
}
