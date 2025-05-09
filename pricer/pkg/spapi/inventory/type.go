package inventory

import (
	"time"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/internal/queryutil"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

type GetInventorySummariesParams struct {
	Details         *bool
	GranularityType *string
	GranularityId   *string
	StartDateTime   *time.Time
	SellerSkus      *[]string
	SellerSku       *string
	NextToken       *string
	MarketplaceIds  *[]string
}
func (p *GetInventorySummariesParams) Stringfy() string {
    return queryutil.StructToQuery(p)
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
