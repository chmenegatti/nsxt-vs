package cmd

import (
	"github.com/chmenegatti/nsxt-vs/api"
	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/database"
)

func SetupDatabase(config *config.Config, edge string) (*database.DatabaseManager, error) {
	dbConfig, err := config.GetDatabaseConfig(edge)
	if err != nil {
		return nil, err
	}
	return database.NewDatabaseManager(dbConfig)
}

func SetupNSXtClient(config *config.Config, edge string) (*api.NSXtAPIClient, error) {
	nsxtConfig, err := config.GetNSXtConfig(edge)
	if err != nil {
		return nil, err
	}
	return api.NewNSXtAPIClient(nsxtConfig), nil
}
