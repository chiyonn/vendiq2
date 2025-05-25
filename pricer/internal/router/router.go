package router

import (
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/handler"
	"github.com/chiyonn/vendiq2/pricer/internal/repository"
	"github.com/chiyonn/vendiq2/pricer/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	database := db.GetDB()
	repo := repository.NewPricingRepository(database)
	srv := service.NewPricingService(repo)

	h := handler.NewPricingHandler(srv)

	r.Get("/health", handler.Health)
	r.Get("/pricings", h.GetAll)

	return r
}

