package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chmenegatti/nsxt-vs/config"
)

type NSXtAPIClient struct {
	config config.NSXtConfig
}

func NewNSXtAPIClient(cfg config.NSXtConfig) *NSXtAPIClient {
	return &NSXtAPIClient{config: cfg}
}

func (c *NSXtAPIClient) FetchData(path string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := c.config.URL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("JSESSIONID=%s", c.config.SessionID))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.config.Auth))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return io.ReadAll(resp.Body)
}

type VirtualServer struct {
	ID            string `json:"id"`
	DisplayName   string `json:"display_name"`
	Path          string `json:"path"`
	LbServicePath string `json:"lb_service_path"`
}

type VsResponse struct {
	Results []VirtualServer `json:"results"`
	Cursor  string          `json:"cursor"`
}

func (c *NSXtAPIClient) GetVirtualServers() ([]VirtualServer, error) {
	var (
		apiResponse VsResponse
		results     VsResponse
	)

	cursor := "00040000"
	for {
		rawData, err := c.FetchData(fmt.Sprintf("/policy/api/v1/infra/lb-virtual-servers/?cursor=%s", cursor))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(rawData, &apiResponse); err != nil {
			return nil, err
		}

		results.Results = append(results.Results, apiResponse.Results...)

		if apiResponse.Cursor == "" {
			break
		}
		cursor = apiResponse.Cursor
		apiResponse.Cursor = ""
	}

	return results.Results, nil
}
