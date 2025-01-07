package nsxt

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type LbServices struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Path        string `json:"path"`
}

func (c *Client) GetLbServices(id string) (string, error) {

	url := c.BaseURL + "/policy/api/v1/infra/lb-services/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Logger.Error("Failed to create request", zap.Error(err))
		return "", err
	}

	req.Header.Add("cookie", "JSESSIONID="+c.SessionID)
	req.Header.Add("Authorization", "Basic "+c.Auth)

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.Logger.Error("Request failed", zap.Error(err))
		return "", err
	}

	body, err := io.ReadAll(res.Body)

	err = res.Body.Close()

	if err != nil {
		c.Logger.Error("Failed to read response body", zap.Error(err))
		return "", err
	}

	var response LbServices
	if err := json.Unmarshal(body, &response); err != nil {
		c.Logger.Error("Failed to parse JSON response", zap.Error(err))
		return "", err
	}

	return response.DisplayName, nil
}

func (c *Client) DeleteLbVs(id string) error {

	url := c.BaseURL + "/policy/api/v1/infra/lb-virtual-servers/" + id

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		c.Logger.Error("Failed to create request", zap.Error(err))
		return err
	}

	req.Header.Add("cookie", "JSESSIONID="+c.SessionID)
	req.Header.Add("Authorization", "Basic "+c.Auth)

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.Logger.Error("Request failed", zap.Error(err))
		return err
	}

	body, err := io.ReadAll(res.Body)

	err = res.Body.Close()

	if err != nil {
		c.Logger.Error("Failed to read response body", zap.Error(err))
		return err
	}

	var response LbServices
	if err := json.Unmarshal(body, &response); err != nil {
		c.Logger.Error("Failed to parse JSON response", zap.Error(err))
		return err
	}

	return nil
}
