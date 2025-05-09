package spapi

import (
	"fmt"
	"net/url"
	"time"
)

func buildQuery(params interface{}) string {
	q := url.Values{}

	// Type assertion + basic field support (example only)
	switch p := params.(type) {
	case *GetInventorySummariesParams:
		if p == nil {
			return ""
		}
		if p.Details != nil {
			q.Set("details", fmt.Sprintf("%t", *p.Details))
		}
		if p.GranularityType != nil {
			q.Set("granularityType", *p.GranularityType)
		}
		if p.GranularityId != nil {
			q.Set("granularityId", *p.GranularityId)
		}
		if p.StartDateTime != nil {
			q.Set("startDateTime", p.StartDateTime.Format(time.RFC3339))
		}
		if p.SellerSkus != nil {
			for _, sku := range *p.SellerSkus {
				q.Add("sellerSkus", sku)
			}
		}
		if p.SellerSku != nil {
			q.Set("sellerSku", *p.SellerSku)
		}
		if p.NextToken != nil {
			q.Set("nextToken", *p.NextToken)
		}
		if p.MarketplaceIds != nil {
			for _, id := range *p.MarketplaceIds {
				q.Add("marketplaceIds", id)
			}
		}
	}
	return q.Encode()
}
