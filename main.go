package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jackpot-mab/experiments-params/controller"
	"jackpot-mab/experiments-params/docs"
	"net/http"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "jackpot-mab:experiments-params")
}

func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/experiment")
		{
			eg.GET("/", controller.ExperimentsParamsController)
			eg.POST("/", controller.AddExperimentParametersController)
			eg.PUT("/", controller.UpdateExperimentParametersController)
		}
	}

	router.GET("/", healthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("localhost:8091")
}
