package cmd

import (
	"fmt"
	"log"

	"github.com/chmenegatti/nsxt-vs/config"
	csvapi "github.com/chmenegatti/nsxt-vs/csv"
	"github.com/chmenegatti/nsxt-vs/operations"
)

func PopulateCSV(config *config.Config, EDGE string) {
	dbManager, err := SetupDatabase(config, EDGE)
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer dbManager.Close()

	rows, err := dbManager.QueryLoadBalances()
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	if err := csvapi.WriteCSV(rows, "nemesis.csv"); err != nil {
		log.Fatalf("Could not write to CSV: %v", err)
	}
	fmt.Println("Data successfully saved to nemesis.csv")

	nsxtClient, err := SetupNSXtClient(config, EDGE)
	if err != nil {
		log.Fatalf("NSX-T client setup failed: %v", err)
	}

	if err := operations.FetchAndSaveNSXtData(nsxtClient); err != nil {
		log.Fatalf("Failed to fetch and save NSX-T data: %v", err)
	}

	if err := csvapi.CompareCSVFiles("nemesis.csv", "nsxt.csv", "diff.csv"); err != nil {
		log.Fatalf("Failed to generate diff CSV: %v", err)
	}
	fmt.Println("Diff CSV generated successfully")

	if err := operations.EnrichDiffCSV(nsxtClient); err != nil {
		log.Fatalf("Failed to enrich diff CSV: %v", err)
	}
}
