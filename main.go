package main

import (
	"log"
	"os"
	"pinmarker/configs"
	"pinmarker/routes"
	"time"

	_ "pinmarker/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       PinMarker API
// @version     1.0
// @description API for PinMarker Mobile, Web, and Telegram Bot tracking feature
// @host        localhost:9001
// @BasePath    /api/v1

func initLogging() {
	now := time.Now()
	logFileName := "logs/pinmarker-" + now.Format("January-2006") + ".log"

	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	initLogging()
	log.Println("Pinmarker API service is starting...")

	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("error loading ENV")
	}

	// Init Firebase
	configs.InitFirebaseApp()

	// Init Gin
	router := gin.Default()

	// Setup Dependencies
	routes.SetUpDependency(router)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Pinmarker is running on port %s\n", port)
	router.Run(":" + port)
}
