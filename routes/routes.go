package routes

import (
	"github.com/Franch62/urls-centralizer/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) *gin.Engine {

	router.Use(cors.Default())

	router.POST("/api/urls", controllers.CreateURL)
	router.GET("/api/urls", controllers.GetAllURLs)
	router.GET("/api/urls/:id/fetch", controllers.FetchYAMLFromURL)
	router.PUT("/api/urls/:id", controllers.UpdateURL)
	router.DELETE("/api/urls/:id", controllers.DeleteURL)
	router.GET("/api/urls/:id/endpoints", controllers.GetURLEndpoints)
	router.GET("/docs/:id", controllers.ServeSwaggerUI)

	router.StaticFile("/swagger.html", "./public/swagger.html")

	return router
}
