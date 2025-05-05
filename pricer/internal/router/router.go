package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
    "github.com/chiyonn/vendiq2/pricer/internal/handler"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	return r
}
