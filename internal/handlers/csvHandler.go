package handlers

import (
	"net/http"

	"github.com/chmenegatti/nsxt-vs/api"
	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/internal/models"
	"github.com/chmenegatti/nsxt-vs/internal/repositories"
	"github.com/gin-gonic/gin"
)

type CSVHandler struct {
	repo *repositories.CSVRepository
}

func NewCSVHandler(repo *repositories.CSVRepository) *CSVHandler {
	return &CSVHandler{repo: repo}
}

func (h *CSVHandler) GetCSVData(c *gin.Context) {
	edge := c.Param("edge")
	data, err := h.repo.LoadCSVData(edge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *CSVHandler) DeleteCSVRecord(c *gin.Context) {
	edge := c.Param("edge")
	id := c.Param("id")

	data, err := h.repo.LoadCSVData(edge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newRecords []models.RegisterCSV

	for _, record := range data {
		if record.ID != id {
			newRecords = append(newRecords, record)
		}
	}

	if err := h.repo.SaveCSVData(newRecords, edge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	configuration, err := config.LoadConfig("config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nsxtConfig, err := configuration.GetNSXtConfig(edge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nsxtClient := api.NewNSXtAPIClient(nsxtConfig)
	if err := nsxtClient.DeleteVirtualServer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}

func (h *CSVHandler) PopulateCSVData(c *gin.Context) {
	edge := c.Param("edge")

	h.repo.GetCSVData(edge)
	c.JSON(http.StatusOK, gin.H{"message": "CSV populated"})
}
