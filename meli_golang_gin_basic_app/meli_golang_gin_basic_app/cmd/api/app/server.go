package app

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Start() {

	//Create default gin router
	router := gin.Default()

	//Set endpoint routes
	mapRoutes(router)

	//verify if a port is specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		panic("could not start app")
	}

}
