package router

import (
	"fmt"
	"net/http"

	"github.com/chiyonn/spapi/auth"
	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/handler"
	"github.com/chiyonn/vendiq2/pricer/internal/repository"
	"github.com/chiyonn/vendiq2/pricer/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg *auth.AuthConfig, httpClient *http.Client) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(LoggingMiddleware)
	database := db.GetDB()
	repo := repository.NewPricingRepository(database)
	qsrv := service.NewQueueService()
	psrv, err := service.NewPricingService(cfg, httpClient, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}

	pricingHandler := handler.NewPricingHandler(psrv)
	queueHandler := handler.NewQueueHandler(qsrv)

	r.Get("/health", handler.Health)

	r.Get("/pricings", pricingHandler.GetAll)
	r.Get("/pricings/sync", pricingHandler.SyncAll)

	r.Get("/queues", queueHandler.GetQueues)
	r.Post("/queue", queueHandler.PostQueue)

	return r, nil
}
