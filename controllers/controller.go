package controller

import (
	"net/http"

	"github.com/chmenegatti/nsxt-vs/database"
	"github.com/chmenegatti/nsxt-vs/nsxt"
	"github.com/chmenegatti/nsxt-vs/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CSVController struct {
	db         *gorm.DB
	logger     *zap.Logger
	nsxtClient *nsxt.Client
}

func NewCSVController(db *gorm.DB, logger *zap.Logger, nsxtClient *nsxt.Client) *CSVController {
	return &CSVController{
		db:         db,
		logger:     logger,
		nsxtClient: nsxtClient,
	}
}

func (c *CSVController) GenerateCSV(ctx echo.Context) error {
	// Busca dados do banco
	records, err := database.FetchLoadBalances(c.db, c.logger)
	if err != nil {
		c.logger.Error("Failed to fetch data", zap.Error(err))
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "failed to fetch data",
			},
		)
	}

	// Prepara dados para o CSV do VS
	data := [][]string{{"vip_port", "address_load_balance", "nsxt_virtual_server_id"}}
	for _, record := range records {
		data = append(
			data, []string{
				record.VIPPort,
				record.AddressLoadBalance,
				record.NSXTServerID,
			},
		)
	}

	// Gera o arquivo CSV do VS
	vs := "vs.csv"
	if err := utils.WriteToCSV(vs, data, c.logger); err != nil {
		c.logger.Error("Failed to write data to CSV", zap.Error(err))
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "failed to write data to CSV",
			},
		)
	}

	// Busca dados do NSXT
	nsxtData, err := c.nsxtClient.GetVirtualServers()
	if err != nil {
		c.logger.Error("Failed to fetch NSXT data", zap.Error(err))
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "failed to fetch NSXT data",
			},
		)
	}

	// Prepara dados para o CSV do NSXT
	data = [][]string{{"id", "display_name", "lb_service_path"}}
	for _, register := range nsxtData {
		data = append(
			data, []string{
				register.ID,
				register.DisplayName,
				register.LBServicePath,
			},
		)
	}

	// Gera o arquivo CSV do NSXT
	c.logger.Info("Data successfully saved to CSV file")
	nsxtfile := "nsxt.csv"
	if err := utils.WriteToCSV(nsxtfile, data, c.logger); err != nil {
		c.logger.Error("Failed to write CSV", zap.Error(err))
		return ctx.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "failed to write CSV",
			},
		)
	}

	return ctx.JSON(
		http.StatusOK, map[string]string{
			"message":      "CSV generated",
			"nemesis data": vs,
			"nsxtData":     nsxtfile,
		},
	)
}
