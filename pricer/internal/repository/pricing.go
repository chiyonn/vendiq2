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

func (r *PricingRepository) ReadByASIN(asin string) (*model.Pricing, error) {
	var pricing model.Pricing
	if err := r.db.Where("asin = ?", asin).First(&pricing).Error; err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *PricingRepository) Create(pricing *model.Pricing) error {
	if err := r.db.Create(pricing).Error; err != nil {
		return err
	}
	return nil
}

func (r *PricingRepository) UpdateByASIN(asin string, pricing *model.Pricing) error {
	if err := r.db.Model(&model.Pricing{}).
		Where("asin = ?", asin).
		Updates(pricing).Error; err != nil {
		return err
	}
	return nil
}
