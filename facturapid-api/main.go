package main

import (
	"log"
	"net/http"
	"os" // For checking file existence
	"time" 

	"facturapid-api/database" 
	"facturapid-api/handlers" 
	"facturapid-api/middleware" 

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" 
)

// Helper function to generate self-signed certificates (for local testing only)
// In a real environment, this would not be part of the application code.
// It's provided here as a comment for developer convenience.
/*
func generateSelfSignedCerts() {
	// Command to generate self-signed certificates using openssl:
	// openssl genpkey -algorithm RSA -out server.key -aes256 # With passphrase
	// openssl rsa -in server.key -out server.key # Remove passphrase
	// openssl req -new -key server.key -out server.csr
	// openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
	//
	// Or a simpler command for a basic self-signed cert without a CSR:
	// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"
	//
	// After generating server.crt and server.key, place them in the same directory as the executable.
	log.Println("---")
	log.Println("Development Mode: To run with HTTPS, generate self-signed certificates:")
	log.Println("1. Install OpenSSL.")
	log.Println("2. Run the following command in your project root:")
	log.Println("   openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj \"/CN=localhost\"")
	log.Println("3. Ensure server.crt and server.key are in the same directory as the executable.")
	log.Println("---")
}
*/


func main() {
	// --- Database Setup ---
	connStr := "postgresql://myuser:mypassword@localhost:5432/facturapid_db?sslmode=disable"

	db, err := database.InitDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Warning: Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}()

	err = database.CreateSchema(db)
	if err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}
	log.Println("Database schema initialized successfully.")
	// --- End Database Setup ---

	// Create a default Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.SecurityHeadersMiddleware()) // Add security headers to all responses

	// --- API Routes ---
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			dbStatus := "UP"
			if err_ping := db.Ping(); err_ping != nil {
				dbStatus = "DOWN"
				log.Printf("Health check: Database ping failed: %v", err_ping)
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":          "DOWN",
					"message":         "Facturapid API is running, but database connection is problematic.",
					"database_status": dbStatus,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":          "UP",
				"message":         "Facturapid API is running smoothly!",
				"database_status": dbStatus,
			})
		})

		invoicesGroup := apiV1.Group("/invoices")
		invoicesGroup.Use(middleware.APIKeyAuthMiddleware()) 
		{
			invoicesGroup.POST("", handlers.CreateInvoiceHandler(db))
			invoicesGroup.GET("/:id", handlers.GetInvoiceHandler(db)) 
			invoicesGroup.PUT("/:id", handlers.UpdateInvoiceFiscalDataHandler(db))
			invoicesGroup.GET("/:id/pdf", handlers.GetInvoicePDFHandler(db))
		}
	}
	// --- End API Routes ---

	port := ":8080" // Default HTTP port
	httpsPort := ":8443" // Default HTTPS port

	// HTTPS Server Setup
	// For production, obtain valid certificates from a Certificate Authority (CA).
	// For local development, you can use self-signed certificates.
	// Ensure server.crt and server.key are present in the application's root directory.
	certFile := "server.crt"
	keyFile := "server.key"

	// Check if certificate files exist. If not, log instructions and fall back to HTTP for local dev.
	// In a production environment, you might want to make HTTPS mandatory and fail if certs are missing.
	useHTTPS := false
	if _, errCert := os.Stat(certFile); errCert == nil {
		if _, errKey := os.Stat(keyFile); errKey == nil {
			useHTTPS = true
		}
	}

	if useHTTPS {
		log.Printf("Starting Facturapid API server on HTTPS port %s", httpsPort)
		// router.RunTLS() will block until the server is shut down.
		if err := router.RunTLS(httpsPort, certFile, keyFile); err != nil {
			log.Fatalf("Failed to run HTTPS server: %v", err)
		}
	} else {
		log.Println("---")
		log.Println("WARNING: Certificate files (server.crt, server.key) not found.")
		log.Println("Falling back to HTTP. For HTTPS in local development, generate self-signed certificates.")
		log.Println("Example (using OpenSSL):")
		log.Println("  openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj \"/CN=localhost\"")
		log.Println("---")
		log.Printf("Starting Facturapid API server on HTTP port %s", port)
		if err := router.Run(port); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}
}
