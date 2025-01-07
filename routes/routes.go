package routes

import (
	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/controllers"
	"github.com/chmenegatti/nsxt-vs/nsxt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, logger *zap.Logger, edge string) {
	server := cfg.Server[edge]
	nsxtClient := nsxt.NewClient(server.URL, server.SessionId, server.Auth, logger)

	csvController := controllers.NewCSVController(db, logger, nsxtClient)

	//pelo echo.Context passar os valord de server.token e server.server

	e.GET("/csv/vs", csvController.GenerateCSV)
	e.GET("/csv/diff", csvController.GetDiff)
	e.DELETE("/csv/vs/:id", csvController.DeleteVs)
}
