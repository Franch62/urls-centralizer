package main

import (
	"net/http"
	"os"

	"github.com/Franch62/urls-centralizer/config"
	"github.com/Franch62/urls-centralizer/models"
	"github.com/Franch62/urls-centralizer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.URL{})

	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback local
	}

	http.ListenAndServe(":"+port, router)
}
