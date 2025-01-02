package cmd

import "github.com/chmenegatti/nsxt-vs/config"

func Run(configuration *config.Config, edge, token string) {
	PopulateCSV(configuration, edge)
	VerifyAndSendSlackMessage(edge, token)
}
