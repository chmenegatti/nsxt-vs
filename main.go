package main

import (
	"fmt"
	"log"

	"github.com/chmenegatti/nsxt-vs/api"
	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/csv"
	"github.com/chmenegatti/nsxt-vs/database"
)

const EdgeServer = "tece01"

func main() {
	// Load configuration
	var yamlConfig config.YAMLConfig
	if err := yamlConfig.Load("config.yaml"); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Database Configuration
	dbConfig, err := yamlConfig.GetDatabaseConfig(EdgeServer)
	if err != nil {
		log.Fatalf("failed to get database configuration: %v", err)
	}

	dbManager, err := database.NewDatabaseManager(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize database manager: %v", err)
	}
	defer func() {
		if err := dbManager.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	// Query Load Balancers
	loadBalances, err := dbManager.QueryLoadBalances()
	if err != nil {
		log.Fatalf("failed to query load balancers: %v", err)
	}

	csvWriter := csv.NewCSVWriter()
	if err := csvWriter.WriteToFile(loadBalances, "nemesis.csv", []string{"id", "display_name", "service"}); err != nil {
		log.Fatalf("failed to write load balancers to CSV: %v", err)
	}
	fmt.Println("Load balancer data saved to nemesis.csv")

	// NSX-T Configuration
	nsxtConfig, err := yamlConfig.GetNSXtConfig(EdgeServer)
	if err != nil {
		log.Fatalf("failed to get NSX-T configuration: %v", err)
	}

	nsxtClient := api.NewNSXtClient(nsxtConfig)
	virtualServers, err := nsxtClient.GetVirtualServers()
	if err != nil {
		log.Fatalf("failed to fetch virtual servers: %v", err)
	}

	virtualServerRecords := convertVirtualServersToRecords(virtualServers)
	if err := csvWriter.WriteToFile(
		virtualServerRecords, "nsxt.csv", []string{"id", "display_name", "service"},
	); err != nil {
		log.Fatalf("failed to write virtual servers to CSV: %v", err)
	}
	fmt.Println("NSX-T virtual server data saved to nsxt.csv")

	// Compare CSV files
	if err := csv.CompareCSVFiles("nemesis.csv", "nsxt.csv", "diff.csv"); err != nil {
		log.Fatalf("failed to compare CSV files: %v", err)
	}
	fmt.Println("Difference between nemesis.csv and nsxt.csv saved to diff.csv")
}

func convertVirtualServersToRecords(servers []api.VirtualServer) [][3]string {
	records := make([][3]string, len(servers))
	for i, server := range servers {
		records[i] = [3]string{server.ID, server.DisplayName, server.Path}
	}
	return records
}
