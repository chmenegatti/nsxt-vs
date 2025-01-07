package main

import (
	"os"

	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/database"
	"github.com/chmenegatti/nsxt-vs/logs"
	"github.com/chmenegatti/nsxt-vs/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	log := logs.InitLogger()
	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {

		}
	}(log)

	edge := os.Getenv("ENV")
	if edge == "" {
		log.Fatal("ENV is not set")
	}

	//Load Config.yaml

	configs, err := config.LoadConfig(edge, log)
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	db, err := database.GetDatabaseConnection(edge, *configs, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*", "http://localhost", "http://localhost:4040", "http://172.0.0.1:4040"},
				AllowMethods: []string{"GET", "DELETE"},
			},
		),
	)

	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("token", configs.Token)
				c.Set("server", configs.Server[edge].Server)
				c.Set("edge", edge)
				return next(c)
			}
		},
	)

	routes.RegisterRoutes(e, db, configs, log, edge)

	log.Info("Starting server on port 4040")
	if e.Start(":4040") != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
