package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/chiyonn/vendiq2/pricer/internal/model"
	"github.com/chiyonn/vendiq2/pricer/internal/repository"

	"gorm.io/gorm"
	"github.com/chiyonn/spapi/auth"
	"github.com/chiyonn/spapi/client"
	"github.com/chiyonn/spapi/endpoint/inventory"
)

const MarketplaceIdJP = "A1VC38T7YXB528"

type PricingService struct {
	repo      *repository.PricingRepository
	inventory *inventory.InventoryAPI
}

func NewPricingService(
	cfg *auth.AuthConfig,
	httpClient *http.Client,
	repo *repository.PricingRepository,
) (*PricingService, error) {
	c, err := client.NewClient(httpClient, "JP", cfg, client.NewRateLimitManager())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pricing service: %w", err)
	}

	return &PricingService{
		repo:      repo,
		inventory: inventory.NewInventoryAPI(c),
	}, nil
}

func (h *PricingService) GetAll() ([]*model.Pricing, error) {
	pricings, err := h.repo.ReadAll()
	if err != nil {
		return nil, err
	}
	return pricings, nil
}

func (h *PricingService) GetByASIN(asin string) (*model.Pricing, error) {
	return h.repo.ReadByASIN(asin)
}

func (h *PricingService) SyncAll() error {
	ctx := context.Background()

	inventory, err := h.FetchAllProductsOnSale(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch inventory summaries: %w", err)
	}

	for _, i := range inventory {
		if i.ASIN == nil {
			continue
		}
		asin := *i.ASIN

		p := &model.Pricing{
			ASIN: asin,
			AutoPricing: false,
		}

		existing, err := h.repo.ReadByASIN(asin)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to read pricing for ASIN %s: %w", asin, err)
		}

		if existing != nil {
			err = h.repo.UpdateByASIN(asin, p)
			if err != nil {
				return fmt.Errorf("failed to update pricing for ASIN %s: %w", asin, err)
			}
		} else {
			err = h.repo.Create(p)
			if err != nil {
				return fmt.Errorf("failed to create pricing for ASIN %s: %w", asin, err)
			}
		}
	}

	return nil
}

func (h *PricingService) FetchAllProductsOnSale(ctx context.Context) ([]inventory.InventorySummary, error) {
	var allSummaries []inventory.InventorySummary
	var nextToken *string
	details := true

	for {
		params := &inventory.GetInventorySummariesParams{
			GranularityType: "Marketplace",
			GranularityId:   MarketplaceIdJP,
			Details:         &details,
			NextToken:       nextToken,
			MarketplaceIds:  []string{MarketplaceIdJP},
		}

		res, err := h.inventory.GetInventorySummaries(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, summary := range res.Payload.InventorySummaries {
			if summary.TotalQuantity != nil && *summary.TotalQuantity > 0 {
				allSummaries = append(allSummaries, summary)
			}
		}

		if res.Pagination == nil || res.Pagination.NextToken == nil || *res.Pagination.NextToken == "" {
			break
		}
		nextToken = res.Pagination.NextToken
	}

	return allSummaries, nil
}

