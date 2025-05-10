package bot

import (
	"context"
	"log"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/client"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/inventory"
)

type PricerBot interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type pricerBot struct {
	client *client.Client
}

func NewPricerBot(c *client.Client) PricerBot {
	return &pricerBot{
		client: c,
	}
}

func (b *pricerBot) Run(ctx context.Context) error {
	return b.fetchAllProductsOnSale(ctx)
}

func (b *pricerBot) Stop(ctx context.Context) error {
	return nil
}

func (b *pricerBot) fetchAllProductsOnSale(ctx context.Context) error {
	var allSummaries []inventory.InventorySummary
	var nextToken *string
	details := true

	for {
		params := &inventory.GetInventorySummariesParams{
			Details:   &details,
			NextToken: nextToken,
		}

		res, err := inventory.GetInventorySummaries(ctx, b.client, params)
		if err != nil {
			return err
		}

		allSummaries = append(allSummaries, res.Payload.InventorySummaries...)

		if res.Pagination == nil || res.Pagination.NextToken == nil || *res.Pagination.NextToken == "" {
			break
		}
		nextToken = res.Pagination.NextToken
	}

	log.Printf("在庫取得できた件数: %d", len(allSummaries))
	return nil
}
