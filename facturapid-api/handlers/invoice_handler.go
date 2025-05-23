package handlers

import (
	"database/sql"
	"errors" // For checking sql.ErrNoRows or custom db errors
	"facturapid-api/database"
	"facturapid-api/dto"
	"facturapid-api/pdfgenerator" // Import the pdfgenerator package
	"fmt"
	"log" // For logging errors
	"net/http"
	"strconv" // For parsing ID from string

	"github.com/gin-gonic/gin"
)

// CreateInvoiceHandler handles the creation of new invoices.
func CreateInvoiceHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fullInvoice dto.FullInvoiceDTO

		if err := c.ShouldBindJSON(&fullInvoice); err != nil {
			log.Printf("Error binding JSON for CreateInvoice: %v", err) 
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": err.Error(), 
			})
			return
		}

		if fullInvoice.Header.Codigo == 0 { 
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": "InvoiceHeader.Codigo is required and cannot be zero.",
			})
			return
		}

		for i := range fullInvoice.Lines {
			if fullInvoice.Lines[i].CodigoFactura != fullInvoice.Header.Codigo {
				errMsg := fmt.Sprintf("Mismatch in CodigoFactura for line item (Producto: %s, Linea: %d). Expected %d, got %d.",
					fullInvoice.Lines[i].Producto,
					fullInvoice.Lines[i].Linea,
					fullInvoice.Header.Codigo,
					fullInvoice.Lines[i].CodigoFactura,
				)
				log.Printf("Validation Error for CreateInvoice: %s", errMsg)
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request payload",
					"details": errMsg,
				})
				return
			}
		}

		err := database.CreateFullInvoice(db, fullInvoice)
		if err != nil {
			log.Printf("Error creating full invoice (Codigo: %d) in database: %v", fullInvoice.Header.Codigo, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process invoice",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"invoice_id": fullInvoice.Header.Codigo,
			"message":    "Invoice created successfully",
		})
	}
}

// GetInvoiceHandler handles retrieving a single invoice by its ID.
func GetInvoiceHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		invoiceID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error parsing invoice ID '%s': %v", idStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format"})
			return
		}

		if invoiceID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice ID must be a positive integer"})
			return
		}

		fullInvoice, err := database.GetFullInvoiceByID(db, invoiceID)
		if err != nil {
			if errors.Is(err, database.ErrInvoiceNotFound) { 
				log.Printf("Invoice not found for ID %d: %v", invoiceID, err)
				c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
				return
			}
			log.Printf("Error retrieving invoice (ID: %d) from database: %v", invoiceID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoice"})
			return
		}
		c.JSON(http.StatusOK, fullInvoice)
	}
}

// UpdateInvoiceFiscalDataHandler handles updating fiscal data for an existing invoice.
func UpdateInvoiceFiscalDataHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		invoiceID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error parsing invoice ID for update '%s': %v", idStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format"})
			return
		}
		if invoiceID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice ID must be a positive integer"})
			return
		}

		var fiscalData dto.FiscalDataDTO
		if err := c.ShouldBindJSON(&fiscalData); err != nil {
			log.Printf("Error binding JSON for UpdateInvoiceFiscalData (ID: %d): %v", invoiceID, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload for fiscal data",
				"details": err.Error(),
			})
			return
		}
		
		if fiscalData.ClienteNombre == nil &&
			fiscalData.ClienteDireccion == nil &&
			fiscalData.ClienteNIF == nil &&
			fiscalData.ClienteEmail == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": "At least one fiscal data field must be provided for update.",
			})
			return
		}

		err = database.UpdateInvoiceFiscalData(db, invoiceID, fiscalData)
		if err != nil {
			if errors.Is(err, database.ErrInvoiceNotFound) {
				log.Printf("Invoice not found for fiscal data update (ID: %d): %v", invoiceID, err)
				c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
				return
			}
			log.Printf("Error updating fiscal data for invoice (ID: %d) in database: %v", invoiceID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice fiscal data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"invoice_id": invoiceID,
			"message":    "Invoice fiscal data updated successfully",
		})
	}
}

// GetInvoicePDFHandler handles generating and returning a PDF for a single invoice.
func GetInvoicePDFHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		invoiceID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error parsing invoice ID for PDF '%s': %v", idStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format"})
			return
		}
		if invoiceID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice ID must be a positive integer"})
			return
		}

		// 1. Fetch invoice data
		fullInvoice, err := database.GetFullInvoiceByID(db, invoiceID)
		if err != nil {
			if errors.Is(err, database.ErrInvoiceNotFound) {
				log.Printf("Invoice not found for PDF generation (ID %d): %v", invoiceID, err)
				c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
				return
			}
			log.Printf("Error retrieving invoice for PDF (ID: %d) from database: %v", invoiceID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoice data"})
			return
		}

		// 2. Generate PDF
		pdfBytes, err := pdfgenerator.GenerateInvoicePDF(fullInvoice)
		if err != nil {
			log.Printf("Error generating PDF for invoice (ID: %d): %v", invoiceID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invoice PDF"})
			return
		}

		// 3. Set headers and return PDF
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"factura_%d.pdf\"", invoiceID))
		// Use "attachment" instead of "inline" to force download.
		
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}
