package main

import (
	"api-server/db"
	"api-server/handlers"
	"api-server/middleware"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize database connections
	if err := db.Init(); err != nil {
		panic(err)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	r.POST("/jobs", middleware.AuthRequired(), handlers.CreateJob)
	r.GET("/jobs", handlers.GetJobStatus)
	r.POST("/clips", middleware.AuthRequired(), handlers.CreateClip)
	r.GET("/videos/:videoId/clips", handlers.GetClips)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}
