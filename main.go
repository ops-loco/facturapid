package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log" // Standard logger for initial setup
	"os"
	"path/filepath"
	"time"

	_ "github.com/alexbrainman/odbc/driver"
	"github.com/kardianos/service"
	qrcode "github.com/skip2/go-qrcode"
)

// Configuration Note:
// Ideally, settings like database connection strings, API endpoints, polling intervals,
// and QR code directories would be externalized into a configuration file (e.g., JSON, YAML, TOML)
// or environment variables, rather than being hardcoded.
// For this example, they remain hardcoded for simplicity.

const (
	serviceName        = "FacturapidSynchronizer"
	serviceDisplayName = "Facturapid Synchronizer Service"
	serviceDescription = "Monitors MS Access DB for new invoices, processes them, and generates QR codes."
	qrCodeDir          = "qrcodes" // Directory to store QR codes
	// Note: In a real service, qrCodeDir should be an absolute path or configurable.
	// For example, it could be relative to the executable or a system app data folder.
	// For simplicity here, it's relative to the working dir of the service.
)

// --- Structs (FacturaData, FacturaLinData, FullInvoice) remain the same ---
type FacturaData struct {
	Codigo      int     `json:"codigo"`
	Fecha       string  `json:"fecha"`
	Hora        string  `json:"hora"`
	Total       float64 `json:"total"`
	Cliente1    string  `json:"cliente1"`
	Impresa     string  `json:"impresa"`
	Base1       float64 `json:"base1"`
	Iva1        float64 `json:"iva1"`
	CuotaIva1   float64 `json:"cuotaIva1"`
	TotalConIva float64 `json:"totalConIva"`
	FormaPago   string  `json:"formaPago"`
}
type FacturaLinData struct {
	ID            int     `json:"id"`
	CodigoFactura int     `json:"codigoFactura"`
	Producto      string  `json:"producto"`
	Descripcion   string  `json:"descripcion"`
	Unidades      float64 `json:"unidades"`
	PrecioUnidad  float64 `json:"precioUnidad"`
	Subtotal      float64 `json:"subtotal"`
	IvaAplicado   float64 `json:"ivaAplicado"`
	TotalLinea    float64 `json:"totalLinea"`
}
type FullInvoice struct {
	Header FacturaData      `json:"header"`
	Lines  []FacturaLinData `json:"lines"`
}

// --- Simulated Database Data (dummyFacturaHeaders, dummyFacturaTable, dummyFacturasLinTable) remain the same ---
var dummyFacturaHeaders = []FacturaData{
	{Codigo: 1, Fecha: "2023-01-15", Hora: "10:00", Total: 121.00, Cliente1: "QR", Impresa: "S"},
	{Codigo: 4, Fecha: "2023-01-15", Hora: "10:15", Total: 242.00, Cliente1: "QR", Impresa: "S"},
	{Codigo: 5, Fecha: "2023-01-15", Hora: "10:20", Total: 60.50, Cliente1: "QR", Impresa: "S"},
}
var dummyFacturaTable = map[int]FacturaData{
	1: {Codigo: 1, Fecha: "2023-01-15", Hora: "10:00", Cliente1: "QR", Impresa: "S", Base1: 100.00, Iva1: 21.0, CuotaIva1: 21.00, Total: 121.00, FormaPago: "Efectivo"},
	4: {Codigo: 4, Fecha: "2023-01-15", Hora: "10:15", Cliente1: "QR", Impresa: "S", Base1: 200.00, Iva1: 21.0, CuotaIva1: 42.00, Total: 242.00, FormaPago: "Tarjeta"},
	5: {Codigo: 5, Fecha: "2023-01-15", Hora: "10:20", Cliente1: "QR", Impresa: "S", Base1: 50.00, Iva1: 21.0, CuotaIva1: 10.50, Total: 60.50, FormaPago: "Efectivo"},
}
var dummyFacturasLinTable = []FacturaLinData{
	{ID: 1, CodigoFactura: 1, Producto: "PROD001", Descripcion: "Product A", Unidades: 2, PrecioUnidad: 25.00, Subtotal: 50.00, IvaAplicado: 21.0, TotalLinea: 60.50},
	{ID: 2, CodigoFactura: 1, Producto: "PROD002", Descripcion: "Product B", Unidades: 1, PrecioUnidad: 50.00, Subtotal: 50.00, IvaAplicado: 21.0, TotalLinea: 60.50},
	{ID: 4, CodigoFactura: 4, Producto: "PROD004", Descripcion: "Product D", Unidades: 4, PrecioUnidad: 50.00, Subtotal: 200.00, IvaAplicado: 21.0, TotalLinea: 242.00},
	{ID: 5, CodigoFactura: 5, Producto: "PROD005", Descripcion: "Product E", Unidades: 1, PrecioUnidad: 50.00, Subtotal: 50.00, IvaAplicado: 21.0, TotalLinea: 60.50},
}
var globalSimCounter = 0 // Renamed to avoid conflict if program struct has its own simCounter

type program struct {
	quit   chan struct{}
	db     *sql.DB
	logger service.Logger
}

func (p *program) Start(s service.Service) error {
	p.logger.Info("Starting service: ", serviceDisplayName)
	p.quit = make(chan struct{})

	connStr := "driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=c:\\tpv\\tpv.mdb;pwd=mcdqn"
	var err error
	p.db, err = sql.Open("odbc", connStr)
	if err != nil {
		p.logger.Errorf("Error opening database connection: %v", err)
		return fmt.Errorf("error opening database connection: %w", err)
	}

	err = p.db.Ping()
	if err != nil {
		p.logger.Warningf("Error pinging database: %v. Proceeding with simulation logic.", err)
		// Depending on requirements, might return error here to prevent service start
	} else {
		p.logger.Info("Successfully connected to the database!")
	}

	go p.runLoop() // Launch the main application logic in a goroutine
	return nil
}

func (p *program) runLoop() {
	p.logger.Info("Application logic loop started.")
	defer p.logger.Info("Application logic loop stopped.")

	lastProcessedCodigo := 0
	// maxPolls is removed for continuous running; service stop will terminate.
	// For demonstration, we might re-introduce a counter or time limit if run interactively.
	ticker := time.NewTicker(10 * time.Second) // Polling interval
	defer ticker.Stop()

	// Initial simCounter for this run
	// Note: If service restarts, simCounter resets. For persistent state, it would need to be stored.
	localSimCounter := 0 

	for {
		select {
		case <-p.quit:
			p.logger.Info("Quit signal received, stopping application logic loop.")
			return
		case <-ticker.C:
			p.logger.Infof("--- Polling iteration (SimCounter: %d) ---", localSimCounter)
			
			// Pass logger to functions that need it, or use p.logger directly
			newInvoiceHeaders, err := p.fetchNewInvoiceHeaders(lastProcessedCodigo, localSimCounter)
			if err != nil {
				p.logger.Errorf("Error fetching new invoice headers: %v", err)
				continue // Continue to next tick
			}

			if len(newInvoiceHeaders) > 0 {
				p.logger.Infof("Found %d new invoice header(s) for processing:", len(newInvoiceHeaders))
				for _, header := range newInvoiceHeaders {
					p.logger.Infof("Processing new invoice header - Codigo: %d", header.Codigo)

					facturaDetails, err := p.getFacturaDetails(header.Codigo)
					if err != nil {
						p.logger.Errorf("Error getting full details for invoice %d: %v. Skipping.", header.Codigo, err)
						continue
					}

					facturaLines, err := p.getFacturaLines(header.Codigo)
					if err != nil {
						p.logger.Warningf("Error getting lines for invoice %d: %v.", header.Codigo, err)
					}

					fullInvoice := FullInvoice{
						Header: facturaDetails,
						Lines:  facturaLines,
					}

					apiCallSuccessful := false
					if err := p.sendInvoiceToAPI(fullInvoice); err != nil {
						p.logger.Errorf("Error sending invoice %d to API: %v", fullInvoice.Header.Codigo, err)
					} else {
						p.logger.Infof("Successfully processed and sent invoice %d to API.", fullInvoice.Header.Codigo)
						apiCallSuccessful = true
					}

					if apiCallSuccessful {
						mockFrontendURL := fmt.Sprintf("https://facturapid.example.com/invoice/INV-%d", fullInvoice.Header.Codigo)
						// Ensure qrCodeDir is an absolute path or handled correctly by the service's working directory
						qrFilename := filepath.Join(qrCodeDir, fmt.Sprintf("invoice_%d.png", fullInvoice.Header.Codigo))
						
						p.logger.Infof("Attempting to generate QR code for URL: %s to file: %s", mockFrontendURL, qrFilename)
						if err := p.generateQRCode(mockFrontendURL, qrFilename); err != nil {
							p.logger.Errorf("Error generating QR code for invoice %d: %v", fullInvoice.Header.Codigo, err)
						} else {
							p.logger.Infof("Successfully generated QR code for invoice %d to %s", fullInvoice.Header.Codigo, qrFilename)
						}
					}

					if header.Codigo > lastProcessedCodigo {
						lastProcessedCodigo = header.Codigo
					}
				}
				p.logger.Infof("Last processed Codigo updated to: %d", lastProcessedCodigo)
			} else {
				p.logger.Info("No new invoice headers found.")
			}
			localSimCounter++ // Increment local simulation counter
		}
	}
}

func (p *program) Stop(s service.Service) error {
	p.logger.Info("Stopping service: ", serviceDisplayName)
	close(p.quit) // Signal the runLoop to exit

	if p.db != nil {
		p.logger.Info("Closing database connection.")
		if err := p.db.Close(); err != nil {
			p.logger.Errorf("Error closing database: %v", err)
			// Potentially return this error, but usually service stop should try to succeed
		}
	}
	return nil
}

// --- Helper methods for program struct (wrapping existing logic) ---
// These methods now use p.logger and p.db

func (p *program) fetchNewInvoiceHeaders(lastCodigo int, currentSimCounter int) ([]FacturaData, error) {
	p.logger.Infof("Simulating fetch for invoice headers with Codigo > %d", lastCodigo)
	var newHeaders []FacturaData
	endIndex := 1 + currentSimCounter*2 
	if endIndex > len(dummyFacturaHeaders) {
		endIndex = len(dummyFacturaHeaders)
	}
	currentBatch := dummyFacturaHeaders
    // Ensure that we don't slice beyond the actual available dummy headers in early stages
	if currentSimCounter < (len(dummyFacturaHeaders)/2)+1 && endIndex <= len(dummyFacturaHeaders) {
		currentBatch = dummyFacturaHeaders[:endIndex]
	} else if endIndex > len(dummyFacturaHeaders) {
        currentBatch = dummyFacturaHeaders // Use all if endIndex goes beyond
    }


	for _, invHeader := range currentBatch {
		if invHeader.Cliente1 == "QR" && invHeader.Impresa == "S" && invHeader.Codigo > lastCodigo {
			newHeaders = append(newHeaders, invHeader)
		}
	}
	return newHeaders, nil
}

func (p *program) getFacturaDetails(invoiceID int) (FacturaData, error) {
	p.logger.Infof("  Simulating getFacturaDetails for Codigo: %d", invoiceID)
	if fd, ok := dummyFacturaTable[invoiceID]; ok {
		return fd, nil
	}
	return FacturaData{}, fmt.Errorf("simulated error: no Factura found with Codigo %d", invoiceID)
}

func (p *program) getFacturaLines(invoiceID int) ([]FacturaLinData, error) {
	p.logger.Infof("  Simulating getFacturaLines for CodigoFactura: %d", invoiceID)
	var lines []FacturaLinData
	for _, line := range dummyFacturasLinTable {
		if line.CodigoFactura == invoiceID {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func (p *program) sendInvoiceToAPI(invoice FullInvoice) error {
	apiEndpoint := "http://localhost:8080/api/invoices" // Mock API endpoint (ideally from config)
	jsonData, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling invoice to JSON: %w", err)
	}
	p.logger.Infof("\n--- Attempting to send Invoice %d to API ---\nEndpoint: %s\nJSON Payload:\n%s",
		invoice.Header.Codigo, apiEndpoint, string(jsonData))
	p.logger.Infof("Simulated: Successfully sent JSON for invoice %d to %s", invoice.Header.Codigo, apiEndpoint)
	return nil
}

func (p *program) generateQRCode(url string, filename string) error {
    dir := filepath.Dir(filename)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        p.logger.Infof("Creating directory: %s", dir)
        err = os.MkdirAll(dir, 0755)
        if err != nil {
            return fmt.Errorf("error creating directory %s: %w", dir, err)
        }
    }
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filename)
	if err != nil {
		return fmt.Errorf("error generating QR code and saving to %s: %w", filename, err)
	}
	return nil
}


func main() {
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Failed to create new service: %v", err)
	}

	// Assign logger from service to program struct
	// This logger might not be available until s.Run() is called,
	// so initial logging in Start might use a temporary logger or standard log.
	// However, s.Logger() can be called after service.New.
	logger, err := s.Logger(nil)
	if err != nil {
		log.Printf("Failed to get service logger: %v", err)
		// Fallback or handle error, for now, prg.logger might be nil until Start
	}
	prg.logger = logger


	// Handle command-line arguments for service control (install, uninstall, etc.)
	// The service package handles 'start', 'stop', 'status' when run by SCM.
	// 'install' and 'uninstall' commands are typically handled explicitly.
	if len(os.Args) > 1 {
		serviceAction := os.Args[1]
		switch serviceAction {
		case "install":
			err := s.Install()
			if err != nil {
				log.Fatalf("Failed to install service: %v", err)
			}
			log.Printf("Service '%s' installed successfully.", serviceDisplayName)
			return
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				log.Fatalf("Failed to uninstall service: %v", err)
			}
			log.Printf("Service '%s' uninstalled successfully.", serviceDisplayName)
			return
		case "start": // For command line, though SCM usually handles this
			err := s.Start()
			if err != nil {
				log.Fatalf("Failed to start service: %v", err)
			}
			log.Printf("Service '%s' started successfully via command line.", serviceDisplayName)
			return
		case "stop": // For command line
			err := s.Stop()
			if err != nil {
				log.Fatalf("Failed to stop service: %v", err)
			}
			log.Printf("Service '%s' stopped successfully via command line.", serviceDisplayName)
			return
		}
	}

	// Run the service. This will block until the service is stopped.
	// Inside s.Run(), it will call prg.Start(), then prg.Run() (if defined, or manage itself), and prg.Stop() on exit.
	// The kardianos/service package's Run() method handles the actual service event loop.
	// Our program's main work loop is initiated in prg.Start() via a goroutine.
	err = s.Run()
	if err != nil {
		if prg.logger != nil {
			prg.logger.Errorf("Service run failed: %v", err)
		} else {
			log.Fatalf("Service run failed: %v", err)
		}
	}
}
