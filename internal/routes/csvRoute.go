package routes

import (
	"github.com/chmenegatti/nsxt-vs/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupCSVRoutes(r *gin.Engine, h *handlers.CSVHandler) {

	r.GET("/csv/:edge", h.GetCSVData)
	r.DELETE("/csv/:edge/:id", h.DeleteCSVRecord)

	r.GET("/run/:edge", h.PopulateCSVData)
}
