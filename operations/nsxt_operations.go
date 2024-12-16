package operations

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/chmenegatti/nsxt-vs/api"
	csvapi "github.com/chmenegatti/nsxt-vs/csv"
	"github.com/chmenegatti/nsxt-vs/utils"
)

func FetchAndSaveNSXtData(client *api.NSXtAPIClient) error {
	servers, err := client.GetVirtualServers()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}

	sort.SliceStable(servers, func(i, j int) bool {
		return utils.CompareIPPort(servers[i], servers[j])
	})

	records := make([][3]string, len(servers))
	for i, server := range servers {
		records[i] = [3]string{server.ID, server.DisplayName, server.LbServicePath}
	}

	if err := csvapi.WriteCSV(records, "nsxt.csv"); err != nil {
		return fmt.Errorf("failed to write CSV: %v", err)
	}

	fmt.Println("Data sorted and saved to nsxt.csv")
	return nil
}

func EnrichDiffCSV(client *api.NSXtAPIClient) error {
	records, err := csvapi.ReadCSVFile("diff.csv")
	if err != nil {
		return fmt.Errorf("error reading diff.csv: %v", err)
	}

	outputFile, err := os.Create("diff_enriched.csv")
	if err != nil {
		return fmt.Errorf("error creating diff_enriched.csv: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	enrichedHeader := append(records[0], "client_code")
	if err := writer.Write(enrichedHeader); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	for i, record := range records[1:] {
		enrichedRecord := record
		if len(record) > 2 {
			servicePath := strings.TrimSpace(record[2])
			clientCode, err := fetchClientCode(client, servicePath)
			if err != nil {
				log.Printf("Failed to fetch client code for service %s: %v", servicePath, err)
				clientCode = ""
			}
			enrichedRecord = append(record, clientCode)
		}

		if err := writer.Write(enrichedRecord); err != nil {
			return fmt.Errorf("error writing record %d: %v", i+1, err)
		}
	}

	fmt.Println("Diff CSV enriched successfully")
	return nil
}

func fetchClientCode(client *api.NSXtAPIClient, servicePath string) (string, error) {
	rawData, err := client.FetchData("/policy/api/v1" + servicePath)
	if err != nil {
		return "", err
	}

	var serviceData map[string]interface{}
	if err := json.Unmarshal(rawData, &serviceData); err != nil {
		return "", err
	}

	clientCode, ok := serviceData["display_name"].(string)
	if !ok {
		return "", fmt.Errorf("display_name not found or not a string")
	}

	return clientCode, nil
}
