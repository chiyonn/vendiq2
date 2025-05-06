package client

import (
    "bytes"
    "log"
    "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/chiyonn/vendiq2/pricer/internal/core"
)

const tokenEndpoint = "https://api.amazon.com/auth/o2/token"

type SPAPIClient struct {
	baseURL        string
	httpClient     *http.Client
    auth           AuthManager
    inventory      InventoryOperator
	accessToken    string
    expiresAt time.Time

	refreshToken   string
	clientID    string
	clientSecret string
    l *log.Logger
    mu sync.Mutex
}

type lwaRefreshTokenRequest struct {
    GrantType    string `json:"grant_type"`
    RefreshToken string `json:"refresh_token"`
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func NewSPAPIClient(l *log.Logger) (*SPAPIClient, error) {

    refreshToken, err := core.ReadSecret("SPAPI_REFRESH_TOKEN")
    if err != nil {
        return nil, err
    }

    clientID, err := core.ReadSecret("SPAPI_CLIENT_ID")
    if err != nil {
        return nil, err
    }

    clientSecret, err := core.ReadSecret("SPAPI_CLIENT_SECRET")
    if err != nil {
        return nil, err
    }

    client := &SPAPIClient{
        baseURL:      "https://sellingpartnerapi-fe.amazon.com",
        httpClient:   &http.Client{Timeout: 10 * time.Second},
        refreshToken: refreshToken,
        clientID:     clientID,
        clientSecret: clientSecret,
        l:            l,
        limiters:     make(map[string]*rate.Limiter),
    }

	return client, nil
}

func (c *SPAPIClient) getAccessToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

    if c.IsTokenValid() {
        return nil
    }

	reqBody, err := json.Marshal(lwaRefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: c.refreshToken,
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch access token: %d, response: %s", resp.StatusCode, string(body))
	}

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.accessToken = data.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(data.ExpiresIn-60) * time.Second)
	c.l.Printf("Access token refreshed, expires in %d seconds", data.ExpiresIn)

	return nil
}

func (c *SPAPIClient) setHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-amz-access-token", c.accessToken)
}

func (c *SPAPIClient) buildRequest(endpoint Endpoint, params *GetInventorySummariesParams) (*http.Request, error) {
	query := buildQuery(params)
	fullURL := fmt.Sprintf("%s%s?%s", c.baseURL, endpoint.Path, query)

	req, err := http.NewRequest(endpoint.Method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeader(req)

    return req, nil
}

func (c *SPAPIClient) getLimiter(key string, rateLimit float64, burst int) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	if l, ok := c.limiters[key]; ok {
		return l
	}

	limiter := rate.NewLimiter(rate.Limit(rateLimit), burst)
	c.limiters[key] = limiter
	return limiter
}


func (c *SPAPIClient) sendRequest(ctx context.Context, endpoint Endpoint, params *GetInventorySummariesParams) ([]byte, error) {
	if err := c.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	limiter := c.getLimiter(endpoint.Path, endpoint.Rate, endpoint.Burst)
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	req, err := c.buildRequest(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("SPAPI %s request failed to %s: %w", endpoint.Method, endpoint.Path, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SPAPI %s request failed to %s: %w", endpoint.Method, endpoint.Path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		c.l.Printf("Token request failed: status=%d body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("SPAPI error (%d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *SPAPIClient) IsTokenValid() bool {
    const tokenExpiryBuffer = 2 * time.Second
    return time.Until(c.expiresAt) > tokenExpiryBuffer
}

func (c *SPAPIClient) GetInventorySummaries(ctx context.Context, params *GetInventorySummariesParams) (*GetInventorySummariesResponse, error) {
	endpoint := Endpoint{
		Method: "GET",
		Path:   "/fba/inventory/v1/summaries",
		Rate:   2, // use later
		Burst:  2, // use later
	}

	body, err := c.sendRequest(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	var response GetInventorySummariesResponse
    var apiErr APIError
    _ = json.Unmarshal(body, &apiErr)
    c.l.Printf("SPAPI error: code=%s, message=%s", apiErr.Code, apiErr.Message)

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func buildQuery(params *GetInventorySummariesParams) string {
	q := url.Values{}

	if params == nil {
		return ""
	}

	if params.Details != nil {
		q.Set("details", fmt.Sprintf("%t", *params.Details))
	}
	if params.GranularityType != nil {
		q.Set("granularityType", *params.GranularityType)
	}
	if params.GranularityId != nil {
		q.Set("granularityId", *params.GranularityId)
	}
	if params.StartDateTime != nil {
		q.Set("startDateTime", params.StartDateTime.Format(time.RFC3339))
	}
	if params.SellerSkus != nil {
		for _, sku := range *params.SellerSkus {
			q.Add("sellerSkus", sku)
		}
	}
	if params.SellerSku != nil {
		q.Set("sellerSku", *params.SellerSku)
	}
	if params.NextToken != nil {
		q.Set("nextToken", *params.NextToken)
	}
	if params.MarketplaceIds != nil {
		for _, id := range *params.MarketplaceIds {
			q.Add("marketplaceIds", id)
		}
	}
	return q.Encode()
}
