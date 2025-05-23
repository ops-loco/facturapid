package database

import (
	"database/sql"
	"facturapid-api/dto" // Import DTO package
	"fmt"
	"log"
	"time" // For parsing string dates to time.Time if necessary

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Custom error for not found, to be specific from db layer
var ErrInvoiceNotFound = sql.ErrNoRows // Reuse sql.ErrNoRows for semantic clarity or define a new one

const (
	createInvoicesTableSQL = `
CREATE TABLE IF NOT EXISTS invoices (
    codigo INTEGER PRIMARY KEY,
    cuenta VARCHAR(20),
    fecha TIMESTAMP WITHOUT TIME ZONE,
    hora TIMESTAMP WITHOUT TIME ZONE,
    total NUMERIC(12,4), 
    tipo_cobro VARCHAR(50),
    vendedor VARCHAR(50),
    cuota_iva NUMERIC(12,4),
    abonado NUMERIC(12,4),
    terminal VARCHAR(15),
    traspasada VARCHAR(1),
    tarifa VARCHAR(10) NOT NULL,
    base1 NUMERIC(12,4),
    base2 NUMERIC(12,4),
    base3 NUMERIC(12,4),
    iva1 NUMERIC(5,2),    -- Assuming IVA percentages like 21.00
    iva2 NUMERIC(5,2),
    iva3 NUMERIC(5,2),
    cuota_iva1 NUMERIC(12,4),
    cuota_iva2 NUMERIC(12,4),
    cuota_iva3 NUMERIC(12,4),
    serie VARCHAR(1) NOT NULL,
    cliente1 VARCHAR(30), 
    cliente2 VARCHAR(30),
    cliente3 VARCHAR(30),
    cliente4 VARCHAR(30), -- This will be used for ClienteEmail
    revisable VARCHAR(1),
    impresa VARCHAR(1),   
    cobro_mixto NUMERIC(12,4),
    efectivo_mixto NUMERIC(12,4),
    tipo_cobro_mixto VARCHAR(50),
    tipo_cobro_mixto2 VARCHAR(50),
    comensales INTEGER,
    codigo_de_factura INTEGER,
    fecha_de_factura TIMESTAMP WITHOUT TIME ZONE,
    hora_de_factura TIMESTAMP WITHOUT TIME ZONE,
    cobro_mixto2 NUMERIC(12,4),
    base4 NUMERIC(12,4),
    base5 NUMERIC(12,4),
    base6 NUMERIC(12,4),
    iva4 NUMERIC(5,2),
    iva5 NUMERIC(5,2),
    iva6 NUMERIC(5,2),
    cuota_iva4 NUMERIC(12,4),
    cuota_iva5 NUMERIC(12,4),
    cuota_iva6 NUMERIC(12,4)
);`

	createInvoiceLinesTableSQL = `
CREATE TABLE IF NOT EXISTS invoice_lines (
    codigo_factura INTEGER NOT NULL,
    unidades_old SMALLINT,
    subtotal NUMERIC(12,4),
    codigo_producto VARCHAR(15),
    producto VARCHAR(200) NOT NULL,
    iva_aplicado NUMERIC(5,2), 
    linea INTEGER NOT NULL,
    unidades NUMERIC(10,4),
    combinado_con VARCHAR(15) NOT NULL, 
    liga_siguiente VARCHAR(1),
    serie VARCHAR(1),
    PRIMARY KEY (codigo_factura, producto, linea),
    FOREIGN KEY (codigo_factura) REFERENCES invoices(codigo) ON DELETE CASCADE
);`

	createInvoicesIndexesSQL = `
CREATE INDEX IF NOT EXISTS idx_invoices_cliente1 ON invoices (cliente1);
CREATE INDEX IF NOT EXISTS idx_invoices_cliente4 ON invoices (cliente4); -- Index on email
CREATE INDEX IF NOT EXISTS idx_invoices_fecha ON invoices (fecha);
CREATE INDEX IF NOT EXISTS idx_invoices_impresa ON invoices (impresa); 
`

	createInvoiceLinesIndexesSQL = `
CREATE INDEX IF NOT EXISTS idx_invoice_lines_codigo_factura ON invoice_lines (codigo_factura);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_codigo_producto ON invoice_lines (codigo_producto);
`
)

// InitDB initializes and returns a PostgreSQL database connection.
func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Println("Successfully connected to the PostgreSQL database.")
	return db, nil
}

// CreateSchema creates the necessary database tables if they don't already exist.
func CreateSchema(db *sql.DB) error {
	log.Println("Attempting to create database schema...")
	if _, err := db.Exec(createInvoicesTableSQL); err != nil {
		return fmt.Errorf("error creating invoices table: %w", err)
	}
	log.Println("Table 'invoices' checked/created successfully.")
	if _, err := db.Exec(createInvoiceLinesTableSQL); err != nil {
		return fmt.Errorf("error creating invoice_lines table: %w", err)
	}
	log.Println("Table 'invoice_lines' checked/created successfully.")
	if _, err := db.Exec(createInvoicesIndexesSQL); err != nil {
		return fmt.Errorf("error creating indexes for invoices table: %w", err)
	}
	log.Println("Indexes for 'invoices' table checked/created successfully.")
	if _, err := db.Exec(createInvoiceLinesIndexesSQL); err != nil {
		return fmt.Errorf("error creating indexes for invoice_lines table: %w", err)
	}
	log.Println("Indexes for 'invoice_lines' table checked/created successfully.")
	log.Println("Database schema creation process completed.")
	return nil
}

// parseDateTime combines date and time strings into a time.Time object.
func parseDateTime(dateStrP *string, timeStrP *string) (sql.NullTime, error) {
	if dateStrP == nil || *dateStrP == "" {
		return sql.NullTime{}, nil
	}
	dateStr := *dateStrP
	timeStr := "00:00:00" 
	if timeStrP != nil && *timeStrP != "" {
		timeStr = *timeStrP
	}
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z07:00", "2006-01-02 15:04"}
	fullDateTimeStr := dateStr + " " + timeStr
	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, fullDateTimeStr)
		if err == nil {
			break
		}
	}
    if err != nil && timeStr == "00:00:00" {
        layoutsDateOnly := []string{"2006-01-02"}
        for _, layout := range layoutsDateOnly {
            parsedTime, err = time.Parse(layout, dateStr)
            if err == nil {
                break
            }
        }
    }
	if err != nil {
		log.Printf("Warning: could not parse date-time string '%s %s': %v. Field will be null.", dateStr, timeStr, err)
		return sql.NullTime{}, nil
	}
	return sql.NullTime{Time: parsedTime, Valid: true}, nil
}

func insertInvoiceHeader(tx *sql.Tx, header dto.InvoiceHeaderDTO) error {
	fecha, _ := parseDateTime(header.Fecha, header.Hora)
	fechaDeFactura, _ := parseDateTime(header.FechaDeFactura, header.HoraDeFactura)
	stmt := `
INSERT INTO invoices (
    codigo, cuenta, fecha, hora, total, tipo_cobro, vendedor, cuota_iva, abonado, terminal,
    traspasada, tarifa, base1, base2, base3, iva1, iva2, iva3, cuota_iva1, cuota_iva2, cuota_iva3,
    serie, cliente1, cliente2, cliente3, cliente4, revisable, impresa, cobro_mixto, efectivo_mixto,
    tipo_cobro_mixto, tipo_cobro_mixto2, comensales, codigo_de_factura, fecha_de_factura,
    hora_de_factura, cobro_mixto2, base4, base5, base6, iva4, iva5, iva6,
    cuota_iva4, cuota_iva5, cuota_iva6
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
    $41, $42, $43, $44, $45, $46
) ON CONFLICT (codigo) DO NOTHING;`
	_, err := tx.Exec(stmt,
		header.Codigo, header.Cuenta, fecha, fecha, header.Total, header.TipoCobro, header.Vendedor, header.CuotaIVA, header.Abonado, header.Terminal,
		header.Traspasada, header.Tarifa, header.Base1, header.Base2, header.Base3, header.Iva1, header.Iva2, header.Iva3, header.CuotaIva1, header.CuotaIva2, header.CuotaIva3,
		header.Serie, header.Cliente1, header.Cliente2, header.Cliente3, header.Cliente4, header.Revisable, header.Impresa, header.CobroMixto, header.EfectivoMixto,
		header.TipoCobroMixto, header.TipoCobroMixto2, header.Comensales, header.CodigoDeFactura, fechaDeFactura,
		fechaDeFactura, header.CobroMixto2, header.Base4, header.Base5, header.Base6, header.Iva4, header.Iva5, header.Iva6,
		header.CuotaIva4, header.CuotaIva5, header.CuotaIva6,
	)
	if err != nil {
		return fmt.Errorf("error inserting invoice header (codigo %d): %w", header.Codigo, err)
	}
	return nil
}

func insertInvoiceLine(tx *sql.Tx, line dto.InvoiceLineDTO, headerCodigo int) error {
	stmt := `
INSERT INTO invoice_lines (
    codigo_factura, unidades_old, subtotal, codigo_producto, producto, iva_aplicado,
    linea, unidades, combinado_con, liga_siguiente, serie
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) ON CONFLICT (codigo_factura, producto, linea) DO NOTHING;`
	_, err := tx.Exec(stmt,
		headerCodigo, line.UnidadesOld, line.Subtotal, line.CodigoProducto, line.Producto, line.IvaAplicado,
		line.Linea, line.Unidades, line.CombinadoCon, line.LigaSiguiente, line.Serie,
	)
	if err != nil {
		return fmt.Errorf("error inserting invoice line (header_codigo %d, producto %s, linea %d): %w", headerCodigo, line.Producto, line.Linea, err)
	}
	return nil
}

func CreateFullInvoice(db *sql.DB, fullInvoice dto.FullInvoiceDTO) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting database transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) 
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Error committing transaction: %v", err)
			} else {
				log.Printf("Transaction committed successfully for invoice %d.", fullInvoice.Header.Codigo)
			}
		}
	}()
	if err = insertInvoiceHeader(tx, fullInvoice.Header); err != nil {
		return err
	}
	for _, line := range fullInvoice.Lines {
		if err = insertInvoiceLine(tx, line, fullInvoice.Header.Codigo); err != nil {
			return err
		}
	}
	return nil
}

func NullableString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NullableInt(i *int) sql.NullInt64 {
    if i == nil {
        return sql.NullInt64{}
    }
    return sql.NullInt64{Int64: int64(*i), Valid: true}
}
func NullableInt16(i *int16) sql.NullInt16 {
    if i == nil {
        return sql.NullInt16{}
    }
    return sql.NullInt16{Int16: *i, Valid: true}
}

func NullableFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func timeToNullTime(t time.Time) sql.NullTime {
    if t.IsZero() {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: t, Valid: true}
}

func GetFullInvoiceByID(db *sql.DB, invoiceID int) (dto.FullInvoiceDTO, error) {
	var fullInvoice dto.FullInvoiceDTO
	var header dto.InvoiceHeaderDTO
	var cuenta, tipoCobro, vendedor, terminal, traspasada, cliente1, cliente2, cliente3, cliente4, revisable, impresa, tipoCobroMixto, tipoCobroMixto2 sql.NullString
	var fecha, hora, fechaDeFactura, horaDeFactura sql.NullTime
	var cuotaIVA, abonado, base1, base2, base3, iva1, iva2, iva3, cuotaIva1, cuotaIva2, cuotaIva3, cobroMixto, efectivoMixto, cobroMixto2, base4, base5, base6, iva4, iva5, iva6, cuotaIva4, cuotaIva5, cuotaIva6 sql.NullFloat64
	var comensales, codigoDeFactura sql.NullInt64
	queryHeader := `
SELECT 
    codigo, cuenta, fecha, hora, total, tipo_cobro, vendedor, cuota_iva, abonado, terminal,
    traspasada, tarifa, base1, base2, base3, iva1, iva2, iva3, cuota_iva1, cuota_iva2, cuota_iva3,
    serie, cliente1, cliente2, cliente3, cliente4, revisable, impresa, cobro_mixto, efectivo_mixto,
    tipo_cobro_mixto, tipo_cobro_mixto2, comensales, codigo_de_factura, fecha_de_factura,
    hora_de_factura, cobro_mixto2, base4, base5, base6, iva4, iva5, iva6,
    cuota_iva4, cuota_iva5, cuota_iva6
FROM invoices WHERE codigo = $1;`
	row := db.QueryRow(queryHeader, invoiceID)
	err := row.Scan(
		&header.Codigo, &cuenta, &fecha, &hora, &header.Total, &tipoCobro, &vendedor, &cuotaIVA, &abonado, &terminal,
		&traspasada, &header.Tarifa, &base1, &base2, &base3, &iva1, &iva2, &iva3, &cuotaIva1, &cuotaIva2, &cuotaIva3,
		&header.Serie, &cliente1, &cliente2, &cliente3, &cliente4, &revisable, &impresa, &cobroMixto, &efectivoMixto,
		&tipoCobroMixto, &tipoCobroMixto2, &comensales, &codigoDeFactura, &fechaDeFactura,
		&horaDeFactura, &cobroMixto2, &base4, &base5, &base6, &iva4, &iva5, &iva6,
		&cuotaIva4, &cuotaIva5, &cuotaIva6,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.FullInvoiceDTO{}, ErrInvoiceNotFound
		}
		return dto.FullInvoiceDTO{}, fmt.Errorf("error querying invoice header (id %d): %w", invoiceID, err)
	}
	if cuenta.Valid { header.Cuenta = &cuenta.String }
	if fecha.Valid { t := fecha.Time.Format("2006-01-02"); header.Fecha = &t }
	if hora.Valid { t := hora.Time.Format("15:04:05"); header.Hora = &t } 
	if tipoCobro.Valid { header.TipoCobro = &tipoCobro.String }
	if vendedor.Valid { header.Vendedor = &vendedor.String }
	if cuotaIVA.Valid { header.CuotaIVA = &cuotaIVA.Float64 }
	if abonado.Valid { header.Abonado = &abonado.Float64 }
	if terminal.Valid { header.Terminal = &terminal.String }
	if traspasada.Valid { header.Traspasada = &traspasada.String }
	if base1.Valid { header.Base1 = &base1.Float64 }
	if base2.Valid { header.Base2 = &base2.Float64 }
	if base3.Valid { header.Base3 = &base3.Float64 }
	if iva1.Valid { header.Iva1 = &iva1.Float64 }
	if iva2.Valid { header.Iva2 = &iva2.Float64 }
	if iva3.Valid { header.Iva3 = &iva3.Float64 }
	if cuotaIva1.Valid { header.CuotaIva1 = &cuotaIva1.Float64 }
	if cuotaIva2.Valid { header.CuotaIva2 = &cuotaIva2.Float64 }
	if cuotaIva3.Valid { header.CuotaIva3 = &cuotaIva3.Float64 }
	if cliente1.Valid { header.Cliente1 = &cliente1.String }
	if cliente2.Valid { header.Cliente2 = &cliente2.String }
	if cliente3.Valid { header.Cliente3 = &cliente3.String }
	if cliente4.Valid { header.Cliente4 = &cliente4.String }
	if revisable.Valid { header.Revisable = &revisable.String }
	if impresa.Valid { header.Impresa = &impresa.String }
	if cobroMixto.Valid { header.CobroMixto = &cobroMixto.Float64 }
	if efectivoMixto.Valid { header.EfectivoMixto = &efectivoMixto.Float64 }
	if tipoCobroMixto.Valid { header.TipoCobroMixto = &tipoCobroMixto.String }
    if tipoCobroMixto2.Valid { header.TipoCobroMixto2 = &tipoCobroMixto2.String }
	if comensales.Valid { v := int(comensales.Int64); header.Comensales = &v }
	if codigoDeFactura.Valid { v := int(codigoDeFactura.Int64); header.CodigoDeFactura = &v }
	if fechaDeFactura.Valid { t := fechaDeFactura.Time.Format("2006-01-02"); header.FechaDeFactura = &t }
	if horaDeFactura.Valid { t := horaDeFactura.Time.Format("15:04:05"); header.HoraDeFactura = &t }
	if cobroMixto2.Valid { header.CobroMixto2 = &cobroMixto2.Float64 }
	if base4.Valid { header.Base4 = &base4.Float64 }
	if base5.Valid { header.Base5 = &base5.Float64 }
	if base6.Valid { header.Base6 = &base6.Float64 }
	if iva4.Valid { header.Iva4 = &iva4.Float64 }
	if iva5.Valid { header.Iva5 = &iva5.Float64 }
	if iva6.Valid { header.Iva6 = &iva6.Float64 }
	if cuotaIva4.Valid { header.CuotaIva4 = &cuotaIva4.Float64 }
	if cuotaIva5.Valid { header.CuotaIva5 = &cuotaIva5.Float64 }
	if cuotaIva6.Valid { header.CuotaIva6 = &cuotaIva6.Float64 }
	fullInvoice.Header = header
	queryLines := `
SELECT 
    codigo_factura, unidades_old, subtotal, codigo_producto, producto, iva_aplicado,
    linea, unidades, combinado_con, liga_siguiente, serie
FROM invoice_lines WHERE codigo_factura = $1 ORDER BY linea ASC;`
	rowsLines, err := db.Query(queryLines, invoiceID)
	if err != nil {
		return dto.FullInvoiceDTO{}, fmt.Errorf("error querying invoice lines for id %d: %w", invoiceID, err)
	}
	defer rowsLines.Close()
	var lines []dto.InvoiceLineDTO
	for rowsLines.Next() {
		var line dto.InvoiceLineDTO
		var unidadesOld sql.NullInt16
		var subtotal, ivaAplicado, unidades sql.NullFloat64
		var codigoProducto, combinadoCon, ligaSiguiente, serie sql.NullString
		err := rowsLines.Scan(
			&line.CodigoFactura, &unidadesOld, &subtotal, &codigoProducto, &line.Producto, &ivaAplicado,
			&line.Linea, &unidades, &combinadoCon, &ligaSiguiente, &serie,
		)
		if err != nil {
			return dto.FullInvoiceDTO{}, fmt.Errorf("error scanning invoice line for id %d: %w", invoiceID, err)
		}
		if unidadesOld.Valid { v:= int16(unidadesOld.Int16); line.UnidadesOld = &v}
		if subtotal.Valid { line.Subtotal = subtotal.Float64 } 
		if codigoProducto.Valid { line.CodigoProducto = &codigoProducto.String }
		if ivaAplicado.Valid { line.IvaAplicado = &ivaAplicado.Float64 }
		if unidades.Valid { line.Unidades = unidades.Float64 } 
		if combinadoCon.Valid { line.CombinadoCon = combinadoCon.String } 
		if ligaSiguiente.Valid { line.LigaSiguiente = &ligaSiguiente.String }
		if serie.Valid { line.Serie = &serie.String }
		lines = append(lines, line)
	}
	if err = rowsLines.Err(); err != nil {
		return dto.FullInvoiceDTO{}, fmt.Errorf("error iterating invoice lines for id %d: %w", invoiceID, err)
	}
	fullInvoice.Lines = lines
	return fullInvoice, nil
}

// UpdateInvoiceFiscalData updates the fiscal information for a specific invoice.
// It maps FiscalDataDTO fields to cliente1, cliente2, cliente3, and cliente4.
func UpdateInvoiceFiscalData(db *sql.DB, invoiceID int, fiscalData dto.FiscalDataDTO) error {
	stmt := `
UPDATE invoices 
SET 
    cliente1 = $1, 
    cliente2 = $2, 
    cliente3 = $3, 
    cliente4 = $4
WHERE codigo = $5;`

	result, err := db.Exec(stmt,
		NullableString(fiscalData.ClienteNombre),
		NullableString(fiscalData.ClienteDireccion),
		NullableString(fiscalData.ClienteNIF),
		NullableString(fiscalData.ClienteEmail),
		invoiceID,
	)
	if err != nil {
		return fmt.Errorf("error executing update for invoice id %d: %w", invoiceID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// This error is less common for UPDATE unless there's a driver issue or connection problem
		return fmt.Errorf("error fetching rows affected for invoice id %d: %w", invoiceID, err)
	}

	if rowsAffected == 0 {
		return ErrInvoiceNotFound // No rows updated means invoice with that ID was not found
	}

	log.Printf("Successfully updated fiscal data for invoice ID %d. Rows affected: %d", invoiceID, rowsAffected)
	return nil
}
