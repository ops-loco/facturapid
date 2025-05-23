package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Configuration Note: API Key should ideally come from a secure configuration source.
const ExpectedAPIKey = "supersecretapikey" // Hardcoded for this example

// APIKeyAuthMiddleware checks for a valid API key in the X-API-Key header.
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		if apiKey != ExpectedAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Next() // Proceed to the next handler
	}
}
