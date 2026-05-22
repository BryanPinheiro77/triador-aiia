package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/BryanPinheiro77/triador-aiia/internal/database"

	"github.com/BryanPinheiro77/triador-aiia/internal/handler"
	"github.com/BryanPinheiro77/triador-aiia/internal/repository"
	"github.com/BryanPinheiro77/triador-aiia/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	router := gin.Default()

	analysisRepository := repository.NewAnalysisRepository()

		analysisService := service.NewAnalysisService(
		analysisRepository,
	)

	analysisHandler := handler.NewAnalysisHandler(
		analysisService,
	)

	router.POST("/analyses", analysisHandler.Create)
	router.GET("/analyses", analysisHandler.FindAll)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API running",
		})
	})

	log.Println("Server running on port 8080")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}