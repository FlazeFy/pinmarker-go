package main

import (
	"log"
	"os"
	"pinmarker/configs"
	"pinmarker/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	// Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Pinmarker is running on port %s\n", port)
	router.Run(":" + port)
}
