package main

import (
	"log"

	"github.com/chmenegatti/nsxt-vs/cmd"
	"github.com/chmenegatti/nsxt-vs/config"
)

const EDGE = "tesp03"

func main() {
	configuration, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	//token := configuration.GetToken()

	cmd.PopulateCSV(configuration, EDGE)

	//cmd.VerifyAndSendSlackMessage(EDGE, token)

}
