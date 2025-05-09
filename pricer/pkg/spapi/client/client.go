package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

const tokenEndpoint = "https://api.amazon.com/auth/o2/token"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *types.Logger

	mu           sync.Mutex
	accessToken  string
	expiresAt    time.Time
	refreshToken string
	clientID     string
	clientSecret string
}

func New(cfg *types.Config, logger *types.Logger) *Client {
	return &Client{
		BaseURL:      cfg.BaseURL,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		Logger:       logger,
		refreshToken: cfg.RefreshToken,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
	}
}

func (c *Client) getAccessToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Until(c.expiresAt) > 2*time.Minute {
		return nil // still valid
	}

	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": c.refreshToken,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token fetch failed: %d, %s", resp.StatusCode, string(body))
	}

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("token decode failed: %w", err)
	}

	c.accessToken = data.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(data.ExpiresIn-60) * time.Second)
	return nil
}

func (c *Client) SendRequest(ctx context.Context, endpoint *types.Endpoint, params types.Queryable) ([]byte, error) {
	if err := c.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("token error: %w", err)
	}

	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, endpoint.Path, params.Stringfy())
	req, err := http.NewRequestWithContext(ctx, endpoint.Method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-amz-access-token", c.accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("SPAPI error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
