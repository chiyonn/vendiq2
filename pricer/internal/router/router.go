package router

import (
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/di"
	"github.com/chiyonn/vendiq2/pricer/internal/handler"
	"github.com/go-chi/chi/v5"
)

func NewRouter(c *di.Container) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)
	r.Get("/pricings", c.PricingHandler.GetAll)

	return r
}
