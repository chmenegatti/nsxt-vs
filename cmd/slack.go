package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	CSVPath      = "diff_enriched.csv"                      // Caminho do arquivo CSV
	SlackWebhook = "https://slack.com/api/chat.postMessage" // URL do webhook do Slack
	SlackChannel = "C05JHLCMVK8"
)

// Row representa uma linha do arquivo CSV
type Row struct {
	ID          string
	DisplayName string
	Service     string
	ClientCode  string
}

// Map para armazenar o estado anterior do CSV
var previousState = make(map[string]Row)

func VerifyAndSendSlackMessage(edge, token string) {
	log.Println("Iniciando a leitura do arquivo CSV...")

	file, err := os.Open(CSVPath)
	if err != nil {
		log.Printf("Erro ao abrir o arquivo CSV: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Erro ao fechar o arquivo CSV: %v\n", err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler o arquivo CSV: %v\n", err)
		return
	}

	// Processa as linhas do CSV
	currentState := make(map[string]Row)
	for i, record := range records {
		// Ignora o cabeçalho
		if i == 0 {
			continue
		}

		row := Row{
			ID:          record[0],
			DisplayName: record[1],
			ClientCode:  record[3],
		}
		currentState[row.ID] = row
	}

	// Detecta atualizações e envia para o Slack
	updates := detectUpdates(currentState, previousState)
	if len(updates) > 0 {
		batchedUpdates := batchRows(updates, 30) // Agrupa em lotes de 30 linhas
		for _, batch := range batchedUpdates {
			sendSlackMessageBatch(batch, edge, token)
		}
		previousState = currentState
		log.Println("Atualizações processadas com sucesso.")
	} else {
		log.Println("Nenhuma atualização detectada.")
	}
}

func detectUpdates(current, previous map[string]Row) []Row {
	var updates []Row
	for id, currentRow := range current {
		if prevRow, exists := previous[id]; !exists || !rowsEqual(currentRow, prevRow) {
			updates = append(updates, currentRow)
		}
	}
	return updates
}

func rowsEqual(a, b Row) bool {
	return a.DisplayName == b.DisplayName && a.ClientCode == b.ClientCode
}

// Agrupa as linhas em lotes de tamanho definido
func batchRows(rows []Row, batchSize int) [][]Row {
	var batches [][]Row
	for batchSize < len(rows) {
		rows, batches = rows[batchSize:], append(batches, rows[0:batchSize:batchSize])
	}
	batches = append(batches, rows)
	return batches
}

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

// Envia uma mensagem com múltiplas linhas para o Slack
func sendSlackMessageBatch(batch []Row, edge, token string) {
	edge = strings.ToUpper(edge)

	mainText := fmt.Sprintf("*%s - Virtual Servers Órfãos Detectados*", edge)
	var attachmentText strings.Builder

	for _, row := range batch {
		attachmentText.WriteString(
			fmt.Sprintf(
				"Floating: %s - CCODE: %s\n ID: %s\n",
				row.DisplayName, row.ClientCode, row.ID,
			),
		)
	}

	payload := SlackPayload{
		Channel: SlackChannel,
		Text:    mainText,
		Attachments: []SlackAttachment{
			{
				Text:           fmt.Sprintf("```%s```", attachmentText.String()),
				Color:          "#D00000",
				AttachmentType: "default",
			},
		},
	}

	// Converte para JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao serializar o payload para o Slack: %v\n", err)
		return
	}

	// Envia a requisição ao Slack
	req, err := http.NewRequest("POST", SlackWebhook, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Erro ao criar requisição Slack: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar mensagem para o Slack: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Erro ao fechar a resposta do Slack: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("Resposta inesperada do Slack: %v\n", resp.Status)
	} else {
		log.Printf("Mensagem enviada ao Slack: %s\n", mainText)
	}
}
