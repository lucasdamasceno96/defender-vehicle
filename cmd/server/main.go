// cmd/server/main.go
package main

import (
	"net/http"
	"time"

	"github.com/lucasdamasceno96/defender-vehicle/internal/handlers"
	"github.com/lucasdamasceno96/defender-vehicle/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Create service instance
	telemetryService := services.NewTelemetryService()

	// 2. Create handler instance (with dependency injection)
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService)

	// 3. Set up Gin router
	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLFiles("static/index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "UP", "time": time.Now()})
		})

		// 4. Register routes with handlers
		api.GET("/telemetry", telemetryHandler.GetTelemetry)
	}

	router.Run(":8080")
}
