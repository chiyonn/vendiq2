package model

import (
	"time"
)

type Pricing struct {
	ASIN           string   `gorm:"column:asin;not null"              json:"asin"`
	MainImageURL   string   `gorm:"column:main_image_url;not null"    json:"mainImageUrl"`
	MinPrice       *float64 `gorm:"column:min_price"                  json:"minPrice,omitempty"`
	MaxPrice       *float64 `gorm:"column:max_price"                  json:"maxPrice,omitempty"`
	NumOfSellers   int      `gorm:"column:num_of_sellers;not null"    json:"numOfSellers"`
	BuyboxPrice    float64  `gorm:"column:buybox_price;not null"      json:"buyboxPrice"`
	BuyboxSellerID string   `gorm:"column:buybox_seller_id;not null"  json:"buyboxSellerId"`
	AutoPricing    bool     `gorm:"column:auto_pricing;not null;default:false" json:"autoPricing"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"    json:"deletedAt,omitempty"`
}
