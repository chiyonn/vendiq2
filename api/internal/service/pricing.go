package service

import (
	"github.com/chiyonn/vendiq2/api/internal/model"
	"github.com/chiyonn/vendiq2/api/internal/repository"
)

type PricingService struct {
	repo *repository.PricingRepository
}

func NewPricingService(repo *repository.PricingRepository) *PricingService {
	return &PricingService{
        repo: repo,
    }
}

func (h *PricingService) GetAll() ([]*model.Pricing, error) {
	pricings, err := h.repo.ReadAll()
	if err != nil {
		return nil, err
	}
	return pricings, nil
}
