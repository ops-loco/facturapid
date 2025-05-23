package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a default Gin router
	// gin.Default() comes with Logger and Recovery middleware.
	router := gin.Default()

	// Define the /health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"message": "Facturapid API is running smoothly!",
		})
	})

	// Define the port the server will listen on
	port := ":8080" // Standard port, can be made configurable

	// Log server startup
	log.Printf("Starting Facturapid API server on port %s", port)

	// Start the HTTP server
	// router.Run() will block until the server is shut down (e.g., by CTRL+C or an error).
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
