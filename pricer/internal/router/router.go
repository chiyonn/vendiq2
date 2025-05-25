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
	psrv := service.NewPricingService(repo)
	qsrv := service.NewQueueService()

	pricingHandler := handler.NewPricingHandler(psrv)
	queueHandler := handler.NewQueueHandler(qsrv)

	r.Get("/health", handler.Health)

	r.Get("/pricings", pricingHandler.GetAll)

	r.Post("/queue", queueHandler.PostQueue)

	return r
}
