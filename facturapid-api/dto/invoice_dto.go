package dto

// InvoiceHeaderDTO corresponds to the data expected for an invoice header.
// It mirrors fields from the 'invoices' table and synchronizer's FacturaData.
type InvoiceHeaderDTO struct {
	Codigo         int     `json:"codigo" binding:"required"` // Primary key, essential
	Cuenta         *string `json:"cuenta"`                    // Use pointers for optional fields
	Fecha          *string `json:"fecha"`                     // Consider time.Time if specific format needed
	Hora           *string `json:"hora"`                      // Consider time.Time
	Total          float64 `json:"total" binding:"omitempty,gte=0"`
	TipoCobro      *string `json:"tipo_cobro"`
	Vendedor       *string `json:"vendedor"`
	CuotaIVA       *float64 `json:"cuota_iva" binding:"omitempty,gte=0"`
	Abonado        *float64 `json:"abonado" binding:"omitempty,gte=0"`
	Terminal       *string `json:"terminal"`
	Traspasada     *string `json:"traspasada" binding:"omitempty,len=1"`
	Tarifa         string  `json:"tarifa" binding:"required"`
	Base1          *float64 `json:"base1" binding:"omitempty,gte=0"`
	Base2          *float64 `json:"base2" binding:"omitempty,gte=0"`
	Base3          *float64 `json:"base3" binding:"omitempty,gte=0"`
	Iva1           *float64 `json:"iva1" binding:"omitempty,gte=0"`   // Percentage
	Iva2           *float64 `json:"iva2" binding:"omitempty,gte=0"`   // Percentage
	Iva3           *float64 `json:"iva3" binding:"omitempty,gte=0"`   // Percentage
	CuotaIva1      *float64 `json:"cuota_iva1" binding:"omitempty,gte=0"`
	CuotaIva2      *float64 `json:"cuota_iva2" binding:"omitempty,gte=0"`
	CuotaIva3      *float64 `json:"cuota_iva3" binding:"omitempty,gte=0"`
	Serie          string  `json:"serie" binding:"required,len=1"`
	Cliente1       *string `json:"cliente1"` // Corresponds to synchronizer's Cliente1
	Cliente2       *string `json:"cliente2"`
	Cliente3       *string `json:"cliente3"`
	Cliente4       *string `json:"cliente4"`
	Revisable      *string `json:"revisable" binding:"omitempty,len=1"`
	Impresa        *string `json:"impresa" binding:"omitempty,len=1"` // Corresponds to synchronizer's Impresa
	CobroMixto     *float64 `json:"cobro_mixto" binding:"omitempty,gte=0"`
	EfectivoMixto  *float64 `json:"efectivo_mixto" binding:"omitempty,gte=0"`
	TipoCobroMixto *string `json:"tipo_cobro_mixto"`
	TipoCobroMixto2 *string `json:"tipo_cobro_mixto2"`
	Comensales     *int    `json:"comensales" binding:"omitempty,gte=0"`
	CodigoDeFactura *int   `json:"codigo_de_factura"` // Assuming this is different from Codigo
	FechaDeFactura *string `json:"fecha_de_factura"`  // Consider time.Time
	HoraDeFactura  *string `json:"hora_de_factura"`   // Consider time.Time
	CobroMixto2    *float64 `json:"cobro_mixto2" binding:"omitempty,gte=0"`
	Base4          *float64 `json:"base4" binding:"omitempty,gte=0"`
	Base5          *float64 `json:"base5" binding:"omitempty,gte=0"`
	Base6          *float64 `json:"base6" binding:"omitempty,gte=0"`
	Iva4           *float64 `json:"iva4" binding:"omitempty,gte=0"` // Percentage
	Iva5           *float64 `json:"iva5" binding:"omitempty,gte=0"` // Percentage
	Iva6           *float64 `json:"iva6" binding:"omitempty,gte=0"` // Percentage
	CuotaIva4      *float64 `json:"cuota_iva4" binding:"omitempty,gte=0"`
	CuotaIva5      *float64 `json:"cuota_iva5" binding:"omitempty,gte=0"`
	CuotaIva6      *float64 `json:"cuota_iva6" binding:"omitempty,gte=0"`
}

// InvoiceLineDTO corresponds to the data expected for an invoice line item.
// It mirrors fields from the 'invoice_lines' table and synchronizer's FacturaLinData.
type InvoiceLineDTO struct {
	CodigoFactura  int     `json:"codigo_factura" binding:"required"` // Should match InvoiceHeaderDTO.Codigo
	UnidadesOld    *int16  `json:"unidades_old" binding:"omitempty,gte=0"`
	Subtotal       float64 `json:"subtotal" binding:"omitempty,gte=0"`
	CodigoProducto *string `json:"codigo_producto"`
	Producto       string  `json:"producto" binding:"required"`
	IvaAplicado    *float64 `json:"iva_aplicado" binding:"omitempty,gte=0"` // Percentage
	Linea          int     `json:"linea" binding:"required,gt=0"` // Line number, should be positive
	Unidades       float64 `json:"unidades" binding:"omitempty,gte=0"`
	// `combinado_con` from DDL seems to be NOT NULL, ensure it's handled.
	// If it's always system-generated or optional from client, adjust binding.
	CombinadoCon   string  `json:"combinado_con"` 
	LigaSiguiente  *string `json:"liga_siguiente" binding:"omitempty,len=1"`
	Serie          *string `json:"serie" binding:"omitempty,len=1"`
}

// FullInvoiceDTO is the top-level structure for the POST /invoices request payload.
type FullInvoiceDTO struct {
	Header InvoiceHeaderDTO   `json:"header" binding:"required"`
	Lines  []InvoiceLineDTO `json:"lines" binding:"omitempty,dive"` // dive validates each element in slice
}
