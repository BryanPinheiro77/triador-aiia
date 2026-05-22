package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/BryanPinheiro77/triador-aiia/internal/database"
	"github.com/BryanPinheiro77/triador-aiia/internal/handler"
	"github.com/BryanPinheiro77/triador-aiia/internal/llm"
	"github.com/BryanPinheiro77/triador-aiia/internal/repository"
	"github.com/BryanPinheiro77/triador-aiia/internal/service"
	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
		},
	}))

	llmClient := llm.NewOpenAIClient()

	analysisRepository := repository.NewAnalysisRepository()

	analysisService := service.NewAnalysisService(
		analysisRepository,
		llmClient,
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
