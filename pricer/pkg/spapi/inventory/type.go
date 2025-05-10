package inventory

import (
	"time"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/internal/queryutil"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

type GetInventorySummariesParams struct {
	Details         *bool      `query:"details"`
	GranularityType *string    `query:"granularityType"`
	GranularityId   *string    `query:"granularityId"`
	StartDateTime   *time.Time `query:"startDateTime"`
	SellerSkus      *[]string  `query:"sellerSkus"`
	SellerSku       *string    `query:"sellerSku"`
	NextToken       *string    `query:"nextToken"`
	MarketplaceIds  *[]string  `query:"marketplaceIds"`
}

func (p *GetInventorySummariesParams) Stringfy() string {
	return queryutil.StructToQuery(p).Encode()
}

type GetInventorySummariesResponse struct {
	Payload    *GetInventorySummariesResult `json:"payload"`
	Pagination *Pagination                  `json:"pagination"`
	Errors     *types.ErrorList             `json:"errors"`
}

type GetInventorySummariesResult struct {
	granularity        Granularity
	InventorySummaries []InventorySummary
}

type Granularity struct {
	granularityType string
	granularityId   string
}

type InventorySummary struct {
	asin             *string
	fnSku            *string
	sellerSku        *string
	condition        *string
	inventoryDetails *InventoryDetails
	lastUpdatedTime  *time.Time
	productName      *string
	totalQuantity    *string
	stores           *[]string
}

type Pagination struct {
	NextToken *string
}

type InventoryDetails struct {
	fulfillableQuantity      *int
	inboundWorkingQuantity   *int
	inboundShippedQuantity   *int
	inboundReceivingQuantity *int
	reservedQuantity         *ReservedQuantity
	researchingQuantity      *ResearchingQuantity
	unfulfillableQuantity    *UnfulfillableQuantity
}

type ReservedQuantity struct {
	totalReservedQuantity        *int
	pendingCustomerOrderQuantity *int
	pendingTransshipmentQuantity *int
	fcProcessingQuantity         *int
}

type ResearchingQuantity struct {
	totalResearchingQuantity     *int
	researchingQuantityBreakdown []ResearchingQuantityEntry
}

type UnfulfillableQuantity struct {
	totalUnfulfillableQuantity *int
	customerDamagedQuantity    *int
	warehouseDamagedQuantity   *int
	distributorDamagedQuantity *int
	carrierDamagedQuantity     *int
	defectiveQuantity          *int
	expiredQuantity            *int
}

type ResearchingQuantityEntry struct {
	name     string
	quantity int
}
