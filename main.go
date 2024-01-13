package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"jackpot-mab/experiments-params/controller"
	"jackpot-mab/experiments-params/db"
	"jackpot-mab/experiments-params/docs"
	"log"
	"net/http"
	"os"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "jackpot-mab:experiments-params")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	dbConnection := db.ConnectionParams{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   os.Getenv("DB_NAME"),
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router := gin.Default()
	router.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}}))

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	experimentParamsController := controller.ExperimentParamsController{
		DAO: db.MakeExperimentsDAO(dbConnection),
	}

	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/experiment")
		{
			eg.GET("/:experiment_id", experimentParamsController.GetExperiment)
			eg.POST("/", experimentParamsController.AddExperiment)
			eg.PUT("/", experimentParamsController.UpdateExperiment)
			eg.POST("/parameter", experimentParamsController.AddOrUpdateParameter)
		}

		v1.GET("/experiments", experimentParamsController.GetExperiments)
	}

	router.GET("/", healthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("0.0.0.0:8091")
}
