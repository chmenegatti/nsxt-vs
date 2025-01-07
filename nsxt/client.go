package nsxt

import (
	"crypto/tls"
	"net/http"

	"go.uber.org/zap"
)

type Client struct {
	BaseURL    string
	SessionID  string
	Auth       string
	Logger     *zap.Logger
	httpClient *http.Client
}

func NewClient(baseURL, sessionID, auth string, logger *zap.Logger) *Client {
	// Create a custom HTTP client with SSL validation disabled
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &Client{
		BaseURL:    baseURL,
		SessionID:  sessionID,
		Auth:       auth,
		Logger:     logger,
		httpClient: httpClient,
	}
}
