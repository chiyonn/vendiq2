package repository

import (
	"github.com/chiyonn/vendiq2/pricer/internal/model"
	"gorm.io/gorm"
)

type PricingRepository struct {
	db *gorm.DB
}

func NewPricingRepository(db *gorm.DB) *PricingRepository {
	return &PricingRepository{db: db}
}

func (r *PricingRepository) ReadAll() ([]*model.Pricing, error) {
	var pricings []*model.Pricing
	if err := r.db.Find(&pricings).Error; err != nil {
		return nil, err
	}
	return pricings, nil
}

