package client_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/client"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type mockLogger struct{}

func (l *mockLogger) Printf(format string, v ...interface{}) {
	// no-op
}

func TestSendRequest_Success(t *testing.T) {
	tokenResponse := `{"access_token": "mocked_token", "expires_in": 3600}`
	apiResponse := `{"message": "success"}`

	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Tokenリクエスト
			if req.URL.String() == "https://api.amazon.com/auth/o2/token" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(tokenResponse)),
				}, nil
			}
			// APIリクエスト
			if req.Header.Get("x-amz-access-token") != "mocked_token" {
				t.Fatalf("unexpected token: %s", req.Header.Get("x-amz-access-token"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(apiResponse)),
			}, nil
		},
	}

	cfg := &types.Config{
		BaseURL:      "https://fake.api",
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		RefreshToken: "refresh_token",
	}

	c := client.New(cfg, &mockLogger{})
	c.HTTPClient = mockClient // 差し替え

	endpoint := &types.Endpoint{
		Method: http.MethodGet,
		Path:   "/test",
	}

	data, err := c.SendRequest(context.Background(), endpoint, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `{"message": "success"}`
	if string(data) != expected {
		t.Fatalf("unexpected response: %s", data)
	}
}

