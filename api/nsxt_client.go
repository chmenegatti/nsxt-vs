package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chmenegatti/nsxt-vs/config"
)

type NSXtClient struct {
	baseURL   string
	sessionID string
	authToken string
}

func NewNSXtClient(cfg config.NSXtConfig) *NSXtClient {
	return &NSXtClient{
		baseURL:   cfg.URL,
		sessionID: cfg.SessionID,
		authToken: cfg.Auth,
	}
}

func (c *NSXtClient) FetchData(endpoint string) ([]byte, error) {
	// Disable SSL verification (for testing purposes only)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("JSESSIONID=%s", c.sessionID))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.authToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

type VirtualServer struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Path        string `json:"path"`
}

type APIResponse struct {
	Results []VirtualServer `json:"results"`
}

func (c *NSXtClient) GetVirtualServers() ([]VirtualServer, error) {
	data, err := c.FetchData("/policy/api/v1/infra/lb-virtual-servers/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch virtual servers: %w", err)
	}

	var response APIResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Results, nil
}
