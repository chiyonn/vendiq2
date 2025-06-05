package bot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chiyonn/spapi/auth"
	"github.com/chiyonn/spapi/client"
	"github.com/chiyonn/spapi/endpoint/inventory"
	"github.com/chiyonn/spapi/endpoint/listingsitem"
	"github.com/chiyonn/spapi/endpoint/productpricing"
	"github.com/chiyonn/vendiq2/pricer/internal/db"
	"github.com/chiyonn/vendiq2/pricer/internal/repository"
	"github.com/chiyonn/vendiq2/pricer/internal/service"
)

const MarketplaceIdJP = "A1VC38T7YXB528"

type PricerBot interface {
	Run(ctx context.Context) error
}

type inventoryAPI interface {
	GetInventorySummaries(ctx context.Context, params *inventory.GetInventorySummariesParams) (*inventory.GetInventorySummariesResponse, error)
}

type DefaultPricerBot struct {
	client         *client.Client
	srv            *service.PricingService
	inventory      inventoryAPI
	listingsitem   *listingsitem.ListingsItemsAPI
	productpricing *productpricing.ProductPricingAPI
}

func NewPricerBot(cfg *auth.AuthConfig, httpClient *http.Client) (PricerBot, error) {
	c, err := client.NewClient(httpClient, "JP", cfg, client.NewRateLimitManager())
	if err != nil {
		return nil, err
	}

	invAPI := inventory.NewInventoryAPI(c)
	listAPI := listingsitem.NewListingsItemsAPI(c)
	pricingAPI := productpricing.NewProductPricingAPI(c)
	database := db.GetDB()
	repo := repository.NewPricingRepository(database)
	srv, err := service.NewPricingService(cfg, httpClient, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing service: %w", err)
	}

	return &DefaultPricerBot{
		client:         c,
		srv:            srv,
		inventory:      invAPI,
		listingsitem:   listAPI,
		productpricing: pricingAPI,
	}, nil
}

func (b *DefaultPricerBot) Run(ctx context.Context) error {
	products, err := b.srv.FetchAllProductsOnSale(ctx)
	if err != nil {
		return err
	}

	adjustmentRequired := b.RetrieveASINsAdjustmentRequired(products)

	// Extract ASINs from filtered products
	var asins []string
	for _, p := range adjustmentRequired {
		if p.ASIN != nil {
			asins = append(asins, *p.ASIN)
		}
	}

	//var pricings []productpricing.GetPricingResult
	const chunkSize = 20
	for i := 0; i < len(asins); i += chunkSize {
		end := i + chunkSize
		if end > len(asins) {
			end = len(asins)
		}
		chunk := asins[i:end]

		params := &productpricing.GetPricingParams{
			MarketplaceIds: b.client.MarketplaceID,
			ASINs:          &chunk,
			ItemType:       "Asin",
		}

		resp, err := b.productpricing.GetPricing(ctx, params)
		if err != nil {
			fmt.Printf("GetPricing error for chunk %v: %v\n", chunk, err)
			continue
		}

		fmt.Printf("Got %d pricing results\n", len(*resp.Payload))
	}

	//var patches []listingsitem.PatchOperation
	//for _, p := range pricings {

	//	val := map[string]any{
	//		"marketplace_id": b.client.MarketplaceID,
	//		"currency":       "JPY",
	//		"our_price": []any{
	//			map[string]any{
	//				"schedule": []any{
	//					map[string]any{
	//						"value_with_tax": 0,
	//					},
	//				},
	//			},
	//		},
	//	}

	//	op := listingsitem.PatchOperation{
	//		OP:    "replace",
	//		Path:  "/attributes/purchasable_offer",
	//		Value: &[]any{val},
	//	}

	//	patches = append(patches, op)

	//}

	//params := listingsitem.PatchListingsItemQuery{
	//	MarketplaceIds: []string{b.client.MarketplaceID},
	//}

	return nil
}

func (b *DefaultPricerBot) FetchAllProductsOnSale(ctx context.Context) ([]inventory.InventorySummary, error) {
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

		res, err := b.inventory.GetInventorySummaries(ctx, params)
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

func (b *DefaultPricerBot) RetrieveASINsAdjustmentRequired(products []inventory.InventorySummary) []inventory.InventorySummary {
	var requireds []inventory.InventorySummary

	for _, p := range products {
		if !b.IsPriceAdjustmentRequired(p) {
			continue
		}
		requireds = append(requireds, p)
	}

	return requireds
}

func (b *DefaultPricerBot) IsPriceAdjustmentRequired(product inventory.InventorySummary) bool {
	return true
}
