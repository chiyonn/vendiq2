package di

import (
	"github.com/chiyonn/vendiq2/pricer/internal/handler"
	"github.com/chiyonn/vendiq2/pricer/internal/repository"
	"github.com/chiyonn/vendiq2/pricer/internal/service"
	"gorm.io/gorm"
)

type Container struct {
	PricingHandler *handler.PricingHandler
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		PricingHandler: handler.NewPricingHandler(service.NewPricingService(repository.NewPricingRepository(db))),
	}
}
