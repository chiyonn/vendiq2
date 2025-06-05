package bot

import (
	"context"
	"fmt"
	"testing"

	"github.com/chiyonn/spapi/endpoint/inventory"
	"github.com/chiyonn/spapi/endpoint/model"
	"github.com/stretchr/testify/assert"
)

type filteringMockInventoryAPI struct{}

func (m *filteringMockInventoryAPI) GetInventorySummaries(ctx context.Context, params *inventory.GetInventorySummariesParams) (*inventory.GetInventorySummariesResponse, error) {
	return &inventory.GetInventorySummariesResponse{
		Payload: &inventory.GetInventorySummariesResult{
			InventorySummaries: []inventory.InventorySummary{
				{ASIN: strPtr("IN_STOCK"), TotalQuantity: intPtr(10)},
				{ASIN: strPtr("OUT_OF_STOCK"), TotalQuantity: intPtr(0)},
			},
		},
		Pagination: nil,
	}, nil
}

type paginatedMockInventoryAPI struct {
	callCount int
}

func (m *paginatedMockInventoryAPI) GetInventorySummaries(ctx context.Context, params *inventory.GetInventorySummariesParams) (*inventory.GetInventorySummariesResponse, error) {
	m.callCount++

	switch m.callCount {
	case 1:
		next := "token-2"
		return &inventory.GetInventorySummariesResponse{
			Payload: &inventory.GetInventorySummariesResult{
				InventorySummaries: []inventory.InventorySummary{
					{ASIN: strPtr("PAGE1_ITEM1"), TotalQuantity: intPtr(5)},
				},
			},
			Pagination: &model.Pagination{
				NextToken: &next,
			},
		}, nil
	case 2:
		return &inventory.GetInventorySummariesResponse{
			Payload: &inventory.GetInventorySummariesResult{
				InventorySummaries: []inventory.InventorySummary{
					{ASIN: strPtr("PAGE2_ITEM1"), TotalQuantity: intPtr(7)},
				},
			},
			Pagination: nil,
		}, nil
	default:
		return nil, nil
	}
}

type errorMockInventoryAPI struct{}

func (m *errorMockInventoryAPI) GetInventorySummaries(ctx context.Context, params *inventory.GetInventorySummariesParams) (*inventory.GetInventorySummariesResponse, error) {
	return nil, fmt.Errorf("inventory fetch error")
}

type errorAfterFirstPageMockInventoryAPI struct {
	callCount int
}

func (m *errorAfterFirstPageMockInventoryAPI) GetInventorySummaries(ctx context.Context, params *inventory.GetInventorySummariesParams) (*inventory.GetInventorySummariesResponse, error) {
	if m.callCount == 0 {
		m.callCount++
		next := "token-next"
		return &inventory.GetInventorySummariesResponse{
			Payload: &inventory.GetInventorySummariesResult{
				InventorySummaries: []inventory.InventorySummary{
					{ASIN: strPtr("PAGE1_ITEM"), TotalQuantity: intPtr(5)},
				},
			},
			Pagination: &model.Pagination{NextToken: &next},
		}, nil
	}
	return nil, fmt.Errorf("inventory fetch error")
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestFetchAllProductsOnSale_FiltersOutOfStock(t *testing.T) {
	b := DefaultPricerBot{
		client:    nil,
		inventory: &filteringMockInventoryAPI{},
	}

	result, err := b.FetchAllProductsOnSale(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "IN_STOCK", *result[0].ASIN)
}

func TestFetchAllProductsOnSale_Pagination(t *testing.T) {
	b := &DefaultPricerBot{
		client:    nil,
		inventory: &paginatedMockInventoryAPI{},
	}

	result, err := b.FetchAllProductsOnSale(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "PAGE1_ITEM1", *result[0].ASIN)
	assert.Equal(t, "PAGE2_ITEM1", *result[1].ASIN)
}

func TestFetchAllProductsOnSale_ErrorOnInitialFetch(t *testing.T) {
	b := &DefaultPricerBot{
		client:    nil,
		inventory: &errorMockInventoryAPI{},
	}

	_, err := b.FetchAllProductsOnSale(context.Background())
	assert.Error(t, err)
}

func TestFetchAllProductsOnSale_ErrorAfterFirstPage(t *testing.T) {
	b := &DefaultPricerBot{
		client:    nil,
		inventory: &errorAfterFirstPageMockInventoryAPI{},
	}

	_, err := b.FetchAllProductsOnSale(context.Background())
	assert.Error(t, err)
}

type customPricerBot struct {
	adjust map[string]bool
}

func (c *customPricerBot) IsPriceAdjustmentRequired(p inventory.InventorySummary) bool {
	if p.ASIN == nil {
		return false
	}
	return c.adjust[*p.ASIN]
}

func (c *customPricerBot) RetrieveASINsAdjustmentRequired(products []inventory.InventorySummary) []inventory.InventorySummary {
	var requireds []inventory.InventorySummary
	for _, p := range products {
		if !c.IsPriceAdjustmentRequired(p) {
			continue
		}
		requireds = append(requireds, p)
	}
	return requireds
}

func TestRetrieveASINsAdjustmentRequired(t *testing.T) {
	b := &customPricerBot{adjust: map[string]bool{"A": true, "B": false}}
	products := []inventory.InventorySummary{
		{ASIN: strPtr("A")},
		{ASIN: strPtr("B")},
	}

	result := b.RetrieveASINsAdjustmentRequired(products)
	assert.Len(t, result, 1)
	assert.Equal(t, "A", *result[0].ASIN)
}
