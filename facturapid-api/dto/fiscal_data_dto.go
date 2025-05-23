package dto

// FiscalDataDTO represents the updatable fiscal information for an invoice's customer.
type FiscalDataDTO struct {
	// ClienteNombre maps to invoices.cliente1
	ClienteNombre *string `json:"cliente_nombre" binding:"omitempty,max=30"` 
	// ClienteDireccion maps to invoices.cliente2
	ClienteDireccion *string `json:"cliente_direccion" binding:"omitempty,max=30"` 
	// ClienteNIF maps to invoices.cliente3
	ClienteNIF *string `json:"cliente_nif" binding:"omitempty,max=30"`       
	// ClienteEmail maps to invoices.cliente4 (assuming varchar(30) is sufficient for email)
	ClienteEmail *string `json:"cliente_email" binding:"omitempty,email,max=30"` 
}

// Note on field lengths: max=30 is based on the original DDL for Cliente1-4 VARCHAR(30).
// If ClienteEmail needs more space, the 'invoices.cliente4' column in PostgreSQL
// would need to be altered (e.g., to VARCHAR(255)).
// The `omitempty` tag means these fields are optional in the PUT request.
// If a field is not provided, it won't be updated (logic to be handled in UPDATE statement or by sending current values).
// For this implementation, if a field is nil in DTO, we might send NULL or skip updating it.
// The current plan is to update all specified fields in the UPDATE statement. If a DTO field is nil,
// the corresponding database field will be set to NULL by the database driver if the pointer is passed directly.
// This is generally acceptable for optional fields.
// If the requirement is to only update non-nil fields, the SQL query builder would be more complex.
// For now, we'll update all fields provided in the DTO, setting to NULL if the DTO field is nil.
