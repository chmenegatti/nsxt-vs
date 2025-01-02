package main

import (
	"log"

	"github.com/chmenegatti/nsxt-vs/internal/handlers"
	"github.com/chmenegatti/nsxt-vs/internal/repositories"
	"github.com/chmenegatti/nsxt-vs/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const EDGE = "tesp03"

func main() {

	setupGin("diff_enriched.csv")

}

func setupGin(filename string) {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1", "http://localhost", "http://10.100.21.11"}
	config.AllowMethods = []string{"GET", "DELETE"}
	r.Use(cors.New(config))

	csvRepo := repositories.NewCSVRepository(filename)
	csvHandler := handlers.NewCSVHandler(csvRepo)

	routes.SetupCSVRoutes(r, csvHandler)

	if err := r.Run(":4040"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
