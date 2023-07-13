package app

import (
	"meli_golang_gin_basic_app/cmd/api/controller"
	"meli_golang_gin_basic_app/cmd/api/service"

	"github.com/gin-gonic/gin"
)

func mapRoutes(router *gin.Engine) {
	healthChecker := controller.HealthChecker{}

	router.GET("/ping", healthChecker.PingHandler)

	databaseService := service.NewDatabase()
	databaseController := controller.NewDatabase(databaseService)

	v1 := router.Group("/api/v1")
	{
		//Persist database connections
		v1.POST("/database", databaseController.Persist)

		//Classify database
		v1.POST("/database/scan/:id", databaseController.Scan)

		//Get database classification
		v1.GET("/database/scan/:id", databaseController.GetClassification)

		//Get palabras clave a validar
		v1.GET("/privateData/wordlist", databaseController.GetWordList)

		//Add new word
		v1.POST("/privateData/wordlist/:word", databaseController.AddNewWord)
	}
}
