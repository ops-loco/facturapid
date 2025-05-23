# Facturapid Technical Specifications

This document provides detailed technical specifications for the components of the Facturapid system.

## 1. Database Schema (PostgreSQL/MySQL)

Two main tables will be used to store the invoice data received from the synchronizer and updated by the frontend.

**Table: `invoices`**

*   Stores the main invoice header information and customer fiscal data.

| Column Name         | Data Type                     | Constraints              | Description                                                                 |
| :------------------ | :---------------------------- | :----------------------- | :-------------------------------------------------------------------------- |
| `id`                | SERIAL / BIGINT AUTO_INCREMENT | PRIMARY KEY              | Unique internal identifier for the invoice in this system.                  |
| `tpv_codigo`        | VARCHAR(50)                   | NOT NULL, INDEX          | Original `Codigo` from the TPV `Factura` table (e.g., "40296").             |
| `tpv_fecha`         | DATE                          | NOT NULL                 | Original `Fecha` from TPV.                                                  |
| `tpv_hora`          | TIME                          | NULLABLE                 | Original `Hora` from TPV.                                                   |
| `total`             | DECIMAL(10, 2)                | NOT NULL                 | Total amount from TPV `Factura`.                                            |
| `base1`             | DECIMAL(10, 2)                | NULLABLE                 | Taxable base 1 from TPV.                                                    |
| `iva1_percent`      | DECIMAL(5, 2)                 | NULLABLE                 | VAT percentage 1 from TPV (`Iva1`).                                         |
| `cuota_iva1`        | DECIMAL(10, 2)                | NULLABLE                 | VAT amount 1 from TPV (`CuotaIva1`).                                        |
| `base2`             | DECIMAL(10, 2)                | NULLABLE                 | Taxable base 2 from TPV.                                                    |
| `iva2_percent`      | DECIMAL(5, 2)                 | NULLABLE                 | VAT percentage 2 from TPV (`Iva2`).                                         |
| `cuota_iva2`        | DECIMAL(10, 2)                | NULLABLE                 | VAT amount 2 from TPV (`CuotaIva2`).                                        |
| `base3`             | DECIMAL(10, 2)                | NULLABLE                 | Taxable base 3 from TPV.                                                    |
| `iva3_percent`      | DECIMAL(5, 2)                 | NULLABLE                 | VAT percentage 3 from TPV (`Iva3`).                                         |
| `cuota_iva3`        | DECIMAL(10, 2)                | NULLABLE                 | VAT amount 3 from TPV (`CuotaIva3`).                                        |
| `fiscal_name`       | VARCHAR(255)                  | NULLABLE                 | Customer's fiscal name (provided via frontend).                             |
| `fiscal_address`    | TEXT                          | NULLABLE                 | Customer's fiscal address (provided via frontend).                          |
| `fiscal_nif_cif`    | VARCHAR(20)                   | NULLABLE                 | Customer's fiscal ID (NIF/CIF) (provided via frontend).                     |
| `fiscal_data_added` | BOOLEAN                       | NOT NULL, DEFAULT FALSE  | Flag indicating if fiscal data has been added by the customer.              |
| `created_at`        | TIMESTAMP WITH TIME ZONE      | NOT NULL, DEFAULT NOW()  | Timestamp when the record was created in this system.                       |
| `updated_at`        | TIMESTAMP WITH TIME ZONE      | NOT NULL, DEFAULT NOW()  | Timestamp when the record was last updated (especially fiscal data).        |
| `raw_tpv_data`      | JSONB / JSON                  | NULLABLE                 | Optional: Store the full raw JSON payload received from the synchronizer. |

**Table: `invoice_lines`**

*   Stores the individual line items for each invoice.

| Column Name     | Data Type        | Constraints                                  | Description                                       |
| :-------------- | :--------------- | :------------------------------------------- | :------------------------------------------------ |
| `id`            | SERIAL / BIGINT AUTO_INCREMENT | PRIMARY KEY                                  | Unique internal identifier for the line item.     |
| `invoice_id`    | BIGINT           | NOT NULL, FOREIGN KEY (invoices.id), INDEX | Links to the corresponding invoice header.        |
| `producto`      | VARCHAR(255)     | NOT NULL                                     | Product description from TPV `FacturasLin`.       |
| `unidades`      | DECIMAL(10, 3)   | NOT NULL                                     | Quantity from TPV `FacturasLin`.                  |
| `subtotal`      | DECIMAL(10, 2)   | NOT NULL                                     | Line item subtotal from TPV `FacturasLin`.        |
| `iva_aplicado`  | DECIMAL(5, 2)    | NULLABLE                                     | VAT percentage applied from TPV `FacturasLin`.    |

## 2. REST API Endpoints (Go - Gin/Echo/Fiber)

Base URL: `/api/v1` (example)
Authentication: API Key required for `POST /invoices` endpoint, sent via `X-API-Key` header.

**1. Submit Invoice Data (from Synchronizer)**

*   **Route:** `POST /invoices`
*   **Auth:** Required (API Key)
*   **Request Body (JSON):**
    ```json
    {
      "tpv_codigo": "40296",
      "tpv_fecha": "2025-05-02",
      "tpv_hora": "18:00:00",
      "total": 121.00,
      "base1": 100.00,
      "iva1_percent": 21.00,
      "cuota_iva1": 21.00,
      // ... other base/iva fields if present ...
      "lines": [
        {
          "producto": "Item 1",
          "unidades": 2,
          "subtotal": 50.00,
          "iva_aplicado": 21.00
        },
        {
          "producto": "Item 2",
          "unidades": 1,
          "subtotal": 50.00,
          "iva_aplicado": 21.00
        }
      ],
      "raw_tpv_data": { ... } // Optional: Full TPV record data
    }
    ```
*   **Success Response (201 Created):**
    ```json
    {
      "id": 12345 // The unique internal ID generated by the API
    }
    ```
*   **Error Responses:**
    *   `400 Bad Request`: Invalid input data.
    *   `401 Unauthorized`: Missing or invalid API Key.
    *   `500 Internal Server Error`: Database error or other server issue.

**2. Get Invoice Details (for Frontend)**

*   **Route:** `GET /invoices/{id}`
*   **Auth:** None required (URL is considered unique/unguessable)
*   **Success Response (200 OK):**
    ```json
    {
      "id": 12345,
      "tpv_codigo": "40296",
      "tpv_fecha": "2025-05-02",
      "tpv_hora": "18:00:00",
      "total": 121.00,
      "base1": 100.00,
      "iva1_percent": 21.00,
      "cuota_iva1": 21.00,
      // ... other base/iva fields ...
      "fiscal_data_added": false,
      "lines": [
        {
          "producto": "Item 1",
          "unidades": 2,
          "subtotal": 50.00,
          "iva_aplicado": 21.00
        },
        // ... other lines ...
      ]
      // Note: Do NOT return sensitive fiscal data here if already added
      // Only return data needed for initial display before form submission
    }
    ```
*   **Error Responses:**
    *   `404 Not Found`: Invoice with the given ID does not exist.
    *   `500 Internal Server Error`.

**3. Update Invoice with Fiscal Data (from Frontend)**

*   **Route:** `PUT /invoices/{id}`
*   **Auth:** None required (URL is considered unique/unguessable)
*   **Request Body (JSON):**
    ```json
    {
      "fiscal_name": "Cliente Ejemplo S.L.",
      "fiscal_address": "Calle Falsa 123, 28080 Madrid",
      "fiscal_nif_cif": "B12345678"
    }
    ```
*   **Success Response (200 OK):**
    ```json
    {
      "message": "Invoice updated successfully"
    }
    ```
*   **Error Responses:**
    *   `400 Bad Request`: Invalid input data.
    *   `404 Not Found`: Invoice with the given ID does not exist.
    *   `409 Conflict`: Fiscal data already added (optional, prevent overwrites).
    *   `500 Internal Server Error`.

**4. Get Invoice PDF (for Frontend)**

*   **Route:** `GET /invoices/{id}/pdf`
*   **Auth:** None required (URL is considered unique/unguessable)
*   **Success Response (200 OK):**
    *   **Headers:**
        *   `Content-Type: application/pdf`
        *   `Content-Disposition: attachment; filename="factura_{tpv_codigo}.pdf"`
    *   **Body:** Binary PDF data.
*   **Error Responses:**
    *   `404 Not Found`: Invoice with the given ID does not exist.
    *   `409 Conflict`: Fiscal data has not been added yet (cannot generate full PDF).
    *   `500 Internal Server Error`: PDF generation error or database error.

## 3. Go Synchronizer Specifications

*   **Language:** Go (latest stable version).
*   **Database Connection:** `alexbrainman/odbc` library for connecting to MS Access `.mdb` via system ODBC driver on Windows 7. Requires appropriate 32/64-bit ODBC driver for Access installed on the Win 7 machine.
*   **Invoice Monitoring:** Timed polling loop (`time.Ticker`) querying the `Factura` table. State persistence (last processed ID or timestamp) stored locally (e.g., in a simple text file or SQLite) to avoid reprocessing.
*   **API Communication:** Standard Go `net/http` client. Use HTTPS. Implement API Key authentication via `X-API-Key` header.
*   **QR Code Generation:** `skip2/go-qrcode` library. Generate QR code image data (e.g., PNG).
*   **Printing (LPT1/ESC/POS):**
    *   Attempt direct file writing to `LPT1:` device if permissions allow.
    *   Alternatively, use `os/exec` to call Windows commands like `COPY /B <qr_image_file_or_commands> LPT1:`.
    *   Construct necessary ESC/POS commands for initializing the printer and printing the QR code image. Refer to Epson TM-T88 command reference.
*   **Windows Service:** `kardianos/service` library for service installation, start, stop, and management.
*   **Configuration:** Use a simple configuration file (e.g., `config.json` or `.env`) for database path, password, API endpoint, API key, polling interval, printer port.
*   **Logging:** Implement robust logging (e.g., using `log` package or a library like `logrus`) to a file for troubleshooting.

## 4. React Frontend Specifications

*   **Language:** JavaScript/TypeScript (TypeScript recommended).
*   **Framework:** React (latest stable version), potentially using Vite or Create React App.
*   **Routing:** `react-router-dom` for handling the `/invoice/:id` route.
*   **API Client:** `axios` or `fetch` API for making HTTPS requests to the backend REST API.
*   **State Management:** Context API with `useReducer` for simple cases, or Zustand/Redux for more complex state needs.
*   **UI Components:** A component library like Material UI, Ant Design, or build custom components styled with CSS Modules, Tailwind CSS, or Styled Components.
*   **Form Handling:** Use a library like `react-hook-form` for efficient form state management and validation.
*   **Validation:** Client-side validation for fiscal data fields (required, format checks for NIF/CIF).
*   **PDF Handling:** Trigger download by navigating to the `GET /invoices/{id}/pdf` endpoint, letting the browser handle the PDF download.
*   **Responsiveness:** Use CSS media queries or framework features to ensure usability on mobile devices.
*   **Deployment:** Build static assets (`npm run build` or `yarn build`) and deploy to Plesk web hosting.

## 5. Security Specifications

*   **HTTPS:** Mandatory for all communication between Frontend <-> API and Synchronizer <-> API. Use valid SSL/TLS certificates.
*   **API Authentication:** Secure API Key mechanism for the Synchronizer -> API connection. The key should be configurable and not hardcoded.
*   **Input Validation:** Rigorous validation on both the API (server-side) and Frontend (client-side) for all user inputs (fiscal data) and data received from the synchronizer.
*   **Input Sanitization:** Sanitize data before storing it in the database or rendering it in the frontend to prevent XSS attacks.
*   **SQL Injection Prevention:** Use parameterized queries or ORMs in the API backend to prevent SQL injection.
*   **Unique Invoice URLs:** Rely on the uniqueness and unguessability of the internal invoice ID in the frontend URL (`/invoice/{id}`) as a form of access control for individual invoices. Avoid sequential IDs if possible (use UUIDs or hashids).
*   **Error Handling:** Avoid exposing sensitive system details in error messages to the end-user.
*   **Dependency Management:** Keep all libraries and frameworks up-to-date to patch known vulnerabilities.

