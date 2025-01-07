package nsxt

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type VirtualServer struct {
	ID            string `json:"id"`
	DisplayName   string `json:"display_name"`
	LBServicePath string `json:"lb_service_path"`
}

type VirtualServerResponse struct {
	Results     []VirtualServer `json:"results"`
	ResultCount int             `json:"result_count"`
	Cursor      string          `json:"cursor"`
}

func (c *Client) GetVirtualServers() ([]VirtualServer, error) {
	var allResults []VirtualServer
	cursor := ""

	for {
		// Construir URL com cursor se não for a primeira página
		url := c.BaseURL + "/policy/api/v1/infra/lb-virtual-servers"
		if cursor != "" {
			url += "?cursor=" + cursor
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.Logger.Error("Failed to create request", zap.Error(err))
			return nil, err
		}

		req.Header.Add("cookie", "JSESSIONID="+c.SessionID)
		req.Header.Add("Authorization", "Basic "+c.Auth)

		res, err := c.httpClient.Do(req)
		if err != nil {
			c.Logger.Error("Request failed", zap.Error(err))
			return nil, err
		}

		body, err := io.ReadAll(res.Body)

		err = res.Body.Close()

		if err != nil {
			c.Logger.Error("Failed to read response body", zap.Error(err))
			return nil, err
		}

		var response VirtualServerResponse
		if err := json.Unmarshal(body, &response); err != nil {
			c.Logger.Error("Failed to parse JSON response", zap.Error(err))
			return nil, err
		}

		// Adiciona os resultados desta página ao slice final
		allResults = append(allResults, response.Results...)

		// Log para acompanhamento do progresso
		c.Logger.Info(
			"Fetched page of virtual servers",
			zap.Int("count", len(response.Results)),
			zap.Int("total_so_far", len(allResults)),
		)

		// Se não há mais cursor, terminamos a paginação
		if response.Cursor == "" {
			break
		}

		// Atualiza o cursor para a próxima página
		cursor = response.Cursor

		// Opcional: adicionar um pequeno delay para não sobrecarregar a API
		time.Sleep(100 * time.Millisecond)
	}

	c.Logger.Info(
		"Completed fetching all virtual servers",
		zap.Int("total_records", len(allResults)),
	)

	return allResults, nil
}
