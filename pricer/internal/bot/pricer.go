package bot

import (
	"context"
	"net/http"

	"github.com/chiyonn/spapi/auth"
	"github.com/chiyonn/spapi/client"
	"github.com/chiyonn/spapi/endpoint/inventory"
)

const MarketplaceIdJP = "A1VC38T7YXB528"

type PricerBot interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type pricerBot struct {
	client    *client.Client
	inventory *inventory.InventoryAPI
}

func NewPricerBot(cfg *auth.AuthConfig, httpClient *http.Client) (PricerBot, error) {
	c, err := client.NewClient(httpClient, "JP", cfg, client.NewRateLimitManager())
	if err != nil {
		return nil, err
	}

	invAPI := inventory.NewInventoryAPI(c)

	return &pricerBot{
		client:    c,
		inventory: invAPI,
	}, nil
}

func (b *pricerBot) Run(ctx context.Context) error {
	_, err := b.fetchAllProductsOnSale(ctx)
	return err
}

func (b *pricerBot) Stop(ctx context.Context) error {
	panic("not implemented")
}

func (b *pricerBot) fetchAllProductsOnSale(ctx context.Context) ([]inventory.InventorySummary, error) {
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

		res, err := b.inventory.GetInventorySummaries(params)
		if err != nil {
			return nil, err
		}

		// TODO: append only items in stock.
		allSummaries = append(allSummaries, res.Payload.InventorySummaries...)

		if res.Pagination == nil || res.Pagination.NextToken == nil || *res.Pagination.NextToken == "" {
			break
		}
		nextToken = res.Pagination.NextToken
	}

	return allSummaries, nil
}
