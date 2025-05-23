package pdfgenerator

import (
	"bytes"
	"facturapid-api/dto"
	"fmt"
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const (
	fontArial           = "Arial"
	fontCourier         = "Courier"
	styleBold           = "B"
	styleItalic         = "I"
	styleBoldItalic     = "BI"
	styleRegular        = ""
	defaultFontSize     = 10
	headerFontSize      = 16
	smallFontSize       = 8
	lineHeight          = 5.5 // mm
	cellGap             = 2   // mm, gap between cells
	defaultLeftMargin   = 10  // mm
	defaultTopMargin    = 10  // mm
	defaultRightMargin  = 10  // mm
	footerHeight        = 15  // mm
	tableHeaderColorR   = 240
	tableHeaderColorG   = 240
	tableHeaderColorB   = 240
)

// Helper to safely get string from pointer, returns "" if nil
func getString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Helper to safely get float64 from pointer, returns 0.0 if nil
func getFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0.0
}

// GenerateInvoicePDF creates a PDF document for the given invoice.
func GenerateInvoicePDF(invoice dto.FullInvoiceDTO) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "") // Portrait, mm, A4 size
	pdf.SetMargins(defaultLeftMargin, defaultTopMargin, defaultRightMargin)
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, defaultTopMargin+footerHeight) // Auto page break with margin for footer

	// Register basic fonts
	pdf.AddFont(fontArial, "", "arial.json") // Ensure arial.json (or other .json) and .z files are available
	pdf.AddFont(fontArial, styleBold, "arialbd.json")
	pdf.AddFont(fontCourier, "", "cour.json")
	pdf.AddFont(fontCourier, styleBold, "courbd.json")

	// --- Invoice Header ---
	pdf.SetFont(fontArial, styleBold, headerFontSize)
	pdf.Cell(0, 10, "FACTURA") // 0 width = full width
	pdf.Ln(12)

	// --- Restaurant/Company Info (Hardcoded) ---
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.Cell(100, lineHeight, "Restaurante Ejemplo, S.L.")
	pdf.SetX(pdf.GetPageWidth() - defaultRightMargin - 80) // Align right for invoice details
	pdf.SetFont(fontArial, styleBold, defaultFontSize)
	pdf.Cell(40, lineHeight, "Factura Nº:")
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.CellFormat(40, lineHeight, fmt.Sprintf("%d", invoice.Header.Codigo), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight)

	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.Cell(100, lineHeight, "Calle Falsa 123, Ciudad Ejemplo")
	pdf.SetX(pdf.GetPageWidth() - defaultRightMargin - 80)
	pdf.SetFont(fontArial, styleBold, defaultFontSize)
	pdf.Cell(40, lineHeight, "Fecha:")
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.CellFormat(40, lineHeight, getString(invoice.Header.Fecha), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight)

	pdf.Cell(100, lineHeight, "NIF: B12345678")
	pdf.SetX(pdf.GetPageWidth() - defaultRightMargin - 80)
	pdf.SetFont(fontArial, styleBold, defaultFontSize)
	pdf.Cell(40, lineHeight, "Hora:")
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.CellFormat(40, lineHeight, getString(invoice.Header.Hora), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight * 2) // Extra space

	// --- Customer Info ---
	pdf.SetFont(fontArial, styleBold, defaultFontSize)
	pdf.Cell(40, lineHeight, "Cliente:")
	pdf.Ln(lineHeight)
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.Cell(0, lineHeight, fmt.Sprintf("Nombre: %s", getString(invoice.Header.Cliente1))) // ClienteNombre
	pdf.Ln(lineHeight)
	pdf.Cell(0, lineHeight, fmt.Sprintf("Dirección: %s", getString(invoice.Header.Cliente2))) // ClienteDireccion
	pdf.Ln(lineHeight)
	pdf.Cell(0, lineHeight, fmt.Sprintf("NIF/CIF: %s", getString(invoice.Header.Cliente3))) // ClienteNIF
	pdf.Ln(lineHeight)
	if email := getString(invoice.Header.Cliente4); email != "" { // ClienteEmail
		pdf.Cell(0, lineHeight, fmt.Sprintf("Email: %s", email))
		pdf.Ln(lineHeight)
	}
	pdf.Ln(lineHeight) // Extra space

	// --- Line Items Table ---
	// Table Header
	pdf.SetFont(fontArial, styleBold, defaultFontSize)
	pdf.SetFillColor(tableHeaderColorR, tableHeaderColorG, tableHeaderColorB)
	colWidths := []float64{95, 20, 25, 20, 30} // Description, Qty, Unit Price, IVA%, Subtotal
	headerTitles := []string{"Descripción", "Cant.", "P. Unit.", "IVA%", "Subtotal"}
	for i, title := range headerTitles {
		pdf.CellFormat(colWidths[i], lineHeight*1.5, title, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(lineHeight * 1.5)

	// Table Rows
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.SetFillColor(255, 255, 255) // White background for rows
	for _, line := range invoice.Lines {
		// Handle potential multi-line descriptions
		descLines := pdf.SplitText(line.Producto, colWidths[0]-cellGap) // -cellGap for padding
		
		// Calculate max lines for this row to keep cells aligned
		numLinesInRow := len(descLines)
		// For simplicity, assume other cells are single line, or truncate / handle overflow
		
		currentY := pdf.GetY()
		currentX := pdf.GetX()

		// Description (can be multi-line)
		for i, descLine := range descLines {
			pdf.SetXY(currentX, currentY+float64(i)*lineHeight)
			pdf.CellFormat(colWidths[0], lineHeight, descLine, "LR", 0, "L", false, 0, "")
		}
		pdf.SetXY(currentX+colWidths[0], currentY) // Move to next cell position

		// Other cells (single line for this example)
		pdf.CellFormat(colWidths[1], lineHeight*float64(numLinesInRow), fmt.Sprintf("%.2f", line.Unidades), "LR", 0, "R", false, 0, "")
		pdf.CellFormat(colWidths[2], lineHeight*float64(numLinesInRow), fmt.Sprintf("%.2f", getUnitPrice(line)), "LR", 0, "R", false, 0, "") // Calculate unit price if not direct in DTO
		pdf.CellFormat(colWidths[3], lineHeight*float64(numLinesInRow), fmt.Sprintf("%.2f", getFloat64(line.IvaAplicado)), "LR", 0, "R", false, 0, "")
		pdf.CellFormat(colWidths[4], lineHeight*float64(numLinesInRow), fmt.Sprintf("%.2f", line.Subtotal), "LR", 0, "R", false, 0, "")
		pdf.Ln(lineHeight * float64(numLinesInRow))
	}
	// Draw bottom line for the table
	pdf.CellFormat(colWidths[0]+colWidths[1]+colWidths[2]+colWidths[3]+colWidths[4], 0, "", "T", 0, "", false, 0, "")
	pdf.Ln(lineHeight)


	// --- Totals Section ---
	// Calculate right alignment position
	totalsLabelWidth := 40.0
	totalsValueWidth := 30.0
	pageWidth, _ := pdf.GetPageSize()
	contentWidth := pageWidth - defaultLeftMargin - defaultRightMargin
	totalsXPos := defaultLeftMargin + contentWidth - totalsLabelWidth - totalsValueWidth

	// Base Imponible (Sum of Bases)
	// For simplicity, assume Base1 is the main base. A real system might sum multiple bases.
	pdf.SetX(totalsXPos)
	pdf.SetFont(fontArial, styleRegular, defaultFontSize)
	pdf.CellFormat(totalsLabelWidth, lineHeight, "Base Imponible:", "", 0, "R", false, 0, "")
	pdf.CellFormat(totalsValueWidth, lineHeight, fmt.Sprintf("%.2f", getFloat64(invoice.Header.Base1)), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight)

	// IVA (Sum of CuotasIVA)
	// For simplicity, assume CuotaIva1 is the main IVA amount.
	pdf.SetX(totalsXPos)
	pdf.CellFormat(totalsLabelWidth, lineHeight, "Total IVA:", "", 0, "R", false, 0, "")
	pdf.CellFormat(totalsValueWidth, lineHeight, fmt.Sprintf("%.2f", getFloat64(invoice.Header.CuotaIva1)), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight * 1.5) // Extra space before grand total

	// Grand Total
	pdf.SetX(totalsXPos)
	pdf.SetFont(fontArial, styleBold, defaultFontSize+2) // Slightly larger and bold
	pdf.CellFormat(totalsLabelWidth, lineHeight, "TOTAL:", "", 0, "R", false, 0, "")
	pdf.CellFormat(totalsValueWidth, lineHeight, fmt.Sprintf("%.2f EUR", invoice.Header.Total), "", 0, "R", false, 0, "")
	pdf.Ln(lineHeight)

	// --- Footer (Example) ---
	pdf.SetY(-(defaultTopMargin + footerHeight - 5)) // Position from bottom
	pdf.SetFont(fontArial, styleItalic, smallFontSize)
	pdf.SetTextColor(128, 128, 128) // Grey
	pdf.CellFormat(0, 10, "Gracias por su preferencia.", "T", 0, "C", false, 0, "")
	// Page number can be added using pdf.PageNo() in a custom Footer function via pdf.SetFooterFunc()

	// Output PDF to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		log.Printf("Error generating PDF output: %v", err)
		return nil, fmt.Errorf("error generating PDF output: %w", err)
	}

	if pdf.Error() != nil {
		log.Printf("Error in PDF generation: %v", pdf.Error())
        return nil, fmt.Errorf("error in PDF generation: %w", pdf.Error())
    }

	return buf.Bytes(), nil
}

// getUnitPrice calculates unit price if not directly available or needs calculation.
// This is a placeholder; actual logic might depend on DTO structure.
func getUnitPrice(line dto.InvoiceLineDTO) float64 {
	if line.Unidades != 0 {
		// Assuming Subtotal is pre-IVA. If Subtotal includes IVA, this logic needs adjustment.
		// Or, if there's a direct PrecioUnidad field, use that.
		// For now, let's assume we need to derive it if a PrecioUnidad field isn't in InvoiceLineDTO.
		// If InvoiceLineDTO had PrecioUnidad: return getFloat64(line.PrecioUnidad)
		return line.Subtotal / line.Unidades 
	}
	return 0.0
}

// Helper to add a cell with proper UTF-8 encoding (if text might contain non-latin chars)
func cellStr(pdf *gofpdf.Fpdf, w, h float64, str, borderStr, alignStr string, fill bool) {
	 // gofpdf handles UTF-8 internally if fonts support it.
	 // For special characters, ensure .json font definition file includes them.
	 // The default core fonts (like Arial, Courier) have limited charsets.
	 // For full UTF-8, one would use pdf.AddUTF8Font and pdf.SetFont("fontname", "", size)
	pdf.CellFormat(w, h, str, borderStr, 0, alignStr, fill, 0, "")
}

// Ensure font files (e.g., arial.json, arialbd.json, cour.json, courbd.json and their .z companions)
// are available in the path gofpdf searches (usually where gofpdf package is, or specified by gofpdf.SetFontLocation()).
// For a production app, these might be embedded or placed in a known assets directory.
// The .json and .z files are typically generated from .ttf files using gofpdf's makefont utility.
// For this example, we assume they are discoverable by gofpdf. If not, AddFont will fail.
// If AddFont fails, it usually prints to stderr but doesn't return an error directly to the caller of AddFont.
// The error is then typically caught when pdf.Output() is called or pdf.Error() is checked.

// Example of how to add a custom font (if needed for better UTF-8 support):
// pdf.AddUTF8Font("DejaVuSans", "", "DejaVuSans.ttf") // Requires DejaVuSans.ttf
// pdf.SetFont("DejaVuSans", "", 10)
// But this requires the .ttf file. For standard fonts, the .json/.z method is used.
// For now, relying on core fonts and basic characters.
// If `arial.json` is not found, `AddFont` will cause `pdf.Error()` to be set.
// This error will be caught before returning from `GenerateInvoicePDF`.
// The `makefont` utility from `jung-kurt/gofpdf/makefont` directory is used to create these.
// Example: go run $GOPATH/pkg/mod/github.com/jung-kurt/gofpdf@v1.16.2/makefont/makefont.go /path/to/arial.ttf outputdir
// This would create arial.json and arial.z in outputdir.
// These files then need to be findable by gofpdf (e.g. in same dir as executable, or standard font path).
// For simplicity, this example assumes they are available. A robust setup would manage these font files.
