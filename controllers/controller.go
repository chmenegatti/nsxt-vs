package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chmenegatti/nsxt-vs/database"
	"github.com/chmenegatti/nsxt-vs/nsxt"
	"github.com/chmenegatti/nsxt-vs/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Response struct {
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
	NemesisData string `json:"nemesis_data,omitempty"`
	NSXTData    string `json:"nsxt_data,omitempty"`
	Differences string `json:"differences,omitempty"`
	Updated     string `json:"updated,omitempty"`
}

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
	vsFilename, err := c.generateLoadBalanceCSV()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
	}

	nsxtFilename, err := c.generateNSXTCSV()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
	}

	diff, err := utils.CompareCSV(nsxtFilename, vsFilename, c.logger)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to compare CSV files"})
	}

	if err := c.addLBDisplayName(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to add LB display name"})
	}

	title := "*Virtual Servers Órfãos Detectados*"
	message := fmt.Sprintf(
		"Existe um total de %d virtual servers órfãos no NSXT\n Clique nesse endereço para verificar: %s", diff-1,
		ctx.Get("server").(string),
	)

	if err := utils.SendSlackMesage(ctx, title, message, c.logger); err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to send slack message"})
	}

	return ctx.JSON(
		http.StatusOK, Response{
			Message:     "CSV files generated successfully",
			NemesisData: vsFilename,
			NSXTData:    nsxtFilename,
			Differences: strconv.Itoa(diff),
			Updated:     "diff_updated.csv",
		},
	)
}

func (c *CSVController) GetDiff(ctx echo.Context) error {

	records, err := utils.ReadFromCSV("diff_updated.csv", c.logger)

	if err != nil {
		c.logger.Error("Failed to read diff CSV file", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to read diff CSV file"})
	}

	data := make([]map[string]string, 0, len(records)-1)

	headers := records[0]
	for _, record := range records[1:] {
		if len(record) != len(headers) {
			c.logger.Warn("Invalid record format", zap.Strings("record", record))
			continue
		}

		recordMap := make(map[string]string)
		for i, header := range headers {
			recordMap[header] = record[i]
		}

		data = append(data, recordMap)
	}

	return ctx.JSON(http.StatusOK, data)

}

func (c *CSVController) DeleteVs(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.nsxtClient.DeleteLbVs(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
	}

	records, err := utils.ReadFromCSV("diff_updated.csv", c.logger)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to read diff CSV file"})
	}

	updatedRecords := make([][]string, 0, len(records))

	headers := records[0]
	updatedRecords = append(updatedRecords, headers)

	data := make(map[string]string)

	for _, record := range records[1:] {
		if record[0] == id {
			for i, header := range headers {
				data[header] = record[i]
			}
		}
		if record[0] != id {
			updatedRecords = append(updatedRecords, record)
		}

		if len(data) == 0 {
			return ctx.JSON(http.StatusNotFound, Response{Error: "virtual server not found"})
		}
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		c.logger.Error("Failed to convert data to JSON", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to convert data to JSON"})
	}

	if err := utils.WriteToCSV("diff_updated.csv", updatedRecords, c.logger); err != nil {
		c.logger.Error("Failed to write updated records to CSV", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, Response{Error: "failed to write updated records to CSV"})
	}

	return ctx.JSON(
		http.StatusOK,
		Response{
			Message: fmt.Sprintf(
				"Virtual server with data: %s deleted successfully", strings.Replace(string(dataJSON), "\"", "", -1),
			),
		},
	)
}

func (c *CSVController) generateLoadBalanceCSV() (string, error) {
	records, err := database.FetchLoadBalances(c.db, c.logger)
	if err != nil {
		c.logger.Error("Failed to fetch load balance data", zap.Error(err))
		return "", fmt.Errorf("failed to fetch load balance data: %w", err)
	}

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

	filename := "vs.csv"
	if err := utils.WriteToCSV(filename, data, c.logger); err != nil {
		c.logger.Error("Failed to write load balance data to CSV", zap.Error(err))
		return "", fmt.Errorf("failed to write load balance data to CSV: %w", err)
	}

	c.logger.Info(
		"Load balance data successfully saved to CSV",
		zap.Int("records", len(data)),
		zap.String("filename", filename),
	)

	return filename, nil
}

func (c *CSVController) generateNSXTCSV() (string, error) {
	nsxtData, err := c.nsxtClient.GetVirtualServers()
	if err != nil {
		c.logger.Error("Failed to fetch NSXT data", zap.Error(err))
		return "", fmt.Errorf("failed to fetch NSXT data: %w", err)
	}

	data := [][]string{{"id", "display_name", "lb_service_path"}}
	for _, register := range nsxtData {
		data = append(
			data, []string{
				register.ID,
				register.DisplayName,
				register.LBServicePath,
			},
		)
	}

	filename := "nsxt.csv"
	if err := utils.WriteToCSV(filename, data, c.logger); err != nil {
		c.logger.Error("Failed to write NSXT data to CSV", zap.Error(err))
		return "", fmt.Errorf("failed to write NSXT data to CSV: %w", err)
	}

	c.logger.Info(
		"NSXT data successfully saved to CSV",
		zap.Int("records", len(nsxtData)),
		zap.String("filename", filename),
	)

	return filename, nil
}

func (c *CSVController) addLBDisplayName() error {

	records, err := utils.ReadFromCSV("diff.csv", c.logger)

	if err != nil {
		c.logger.Error("Failed to read diff CSV file", zap.Error(err))
		return fmt.Errorf("failed to read diff CSV file: %w", err)
	}

	updatedRecords := make([][]string, 0, len(records))

	headers := []string{"id", "display_name", "client_code"}
	updatedRecords = append(updatedRecords, headers)

	for _, record := range records[1:] {
		if len(record) < 3 {
			c.logger.Warn("Invalid record format", zap.Strings("record", record))
			continue
		}

		id := strings.Split(record[2], "/")
		if len(id) < 4 {
			c.logger.Warn("Invalid service path format", zap.String("path", record[2]))
			continue
		}

		serviceID := id[3]

		lbServiceName, err := c.nsxtClient.GetLbServices(serviceID)
		if err != nil {
			c.logger.Error(
				"Failed to fetch LB service name",
				zap.String("service_id", serviceID),
				zap.Error(err),
			)
			continue
		}

		newRecord := make([]string, len(record))
		copy(newRecord, record)
		newRecord[2] = lbServiceName

		updatedRecords = append(updatedRecords, newRecord)
	}

	err = utils.WriteToCSV("diff_updated.csv", updatedRecords, c.logger)

	if err != nil {
		c.logger.Error("Failed to write updated records to CSV", zap.Error(err))
		return fmt.Errorf("failed to write updated records to CSV: %w", err)
	}

	c.logger.Info(
		"Successfully updated LB service names",
		zap.Int("total_records", len(updatedRecords)-1),
		zap.String("output_file", "diff_updated.csv"),
	)

	return nil
}
