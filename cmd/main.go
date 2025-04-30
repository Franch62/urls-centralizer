package main

import (
	"urls-centralizer/config"
	"urls-centralizer/models"
	"urls-centralizer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.URL{})

	router := gin.Default()

	r := routes.SetupRoutes(router)
	r.Run(":8080")
}
