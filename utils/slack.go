package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type SlackPayload struct {
	Channel     string            `json:"channel"`
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Text           string `json:"text"`
	Color          string `json:"color"`
	AttachmentType string `json:"attachment_type"`
}

const (
	CSVPath      = "diff_updated.csv"
	SlackWebhook = "https://slack.com/api/chat.postMessage"
	SlackChannel = "C05JHLCMVK8"
)

func SendSlackMesage(edge, server, token string, logger *zap.Logger) error {

	records, err := ReadFromCSV(CSVPath, logger)
	if err != nil {
		logger.Error("Failed to read CSV file", zap.Error(err))
		return err
	}

	edge = strings.ToUpper(edge)
	mainText := fmt.Sprintf("*Virtual Servers Órfãos Detectados - %s*", edge)

	payload := SlackPayload{
		Channel: SlackChannel,
		Text:    mainText,
		Attachments: []SlackAttachment{
			{
				Text: fmt.Sprintf(
					"Existe um total de %d virtual servers órfãos no %s\n Clique nesse endereço para verificar: %s",
					len(records)-1, edge, server,
				),
				Color:          "#D00000",
				AttachmentType: "default",
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to convert payload to JSON", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", SlackWebhook, bytes.NewBuffer(data))
	if err != nil {
		logger.Error("Failed to create request", zap.Error(err))
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to send message to Slack", zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to send message to Slack", zap.Error(err))
		return fmt.Errorf("failed to send message to Slack: %s", resp.Status)
	}

	logger.Info("Message sent to Slack successfully")

	return nil
}
