package models

type RegisterCSV struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Service     string `json:"service"`
	ClientCode  string `json:"client_code"`
}
