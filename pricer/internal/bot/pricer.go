package bot

import (
	"github.com/chiyonn/vendiq2/pricer/internal/client"
)

type PricerBot struct {
	BaseBot
	spapi *client.SPAPIClient
}

func NewPricerBot(spapi *client.SPAPIClient) *PricerBot {
	return &PricerBot{
		spapi: spapi,
	}
}

func (b *PricerBot) Run() {
	for {
        b.fetchAllProductsOnSale()
	}
}

func (b *PricerBot) fetchAllProductsOnSale() {
    params := client.GetInventorySummariesParams{}

    res, err := b.spapi.GetInventorySummaries(&params)
    if err != nil {
        return nil, err
    }
}
