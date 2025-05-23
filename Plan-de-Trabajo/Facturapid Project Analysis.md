## Facturapid Project Analysis

This document summarizes the requirements and constraints for the Facturapid project based on the provided prompt.

### 1. Go Synchronizer (Backend Service)

*   **Objective:** Monitor an existing TPV database for specific invoices, extract data, send it to a central API, generate a QR code linking to a frontend, and print the QR code on a ticket printer.
*   **Environment:** Windows 7 (must run as a service with minimal OS dependencies).
*   **Database Interaction:**
    *   Connect to Microsoft Access database (`.mdb`) located at `c:\tpv\tpv.mdb`.
    *   Database password: `mcdqn`.
    *   Use an appropriate Go library for MS Access connection (e.g., via ODBC).
*   **Monitoring:**
    *   Target Table: `Factura`.
    *   Trigger Condition: New records where `Cliente1 = "QR"` AND `Impresa = "S"`.
    *   Mechanism: Implement polling or leverage system events if possible (polling is more likely feasible).
*   **Data Processing:**
    *   Extract all fields from the qualifying `Factura` record.
    *   Retrieve corresponding line items from the `FacturasLin` table using `CodigoFactura`.
*   **API Communication:**
    *   Send the extracted invoice data (header and lines) to the central REST API.
    *   Implement secure communication (authentication needed).
*   **QR Code Generation:**
    *   Generate a QR code containing a unique URL pointing to the frontend application for the specific invoice.
    *   Use a suitable Go library for QR code generation.
*   **Printing:**
    *   Print the generated QR code on an Epson TM-T88 ticket printer.
    *   Printer connection: LPT1 port.
    *   Use appropriate communication protocol (e.g., ESC/POS commands).
*   **Deployment:**
    *   Provide an installer or script to easily install the application as a Windows service.

### 2. REST API (Backend)

*   **Objective:** Receive invoice data from the synchronizer, store it, manage fiscal data updates from the frontend, and generate final PDF invoices.
*   **Technology:** Framework like Echo, Gin, or Fiber (Go suggested). Database: PostgreSQL or MySQL.
*   **Endpoints:**
    *   `POST /invoices`: Receive and store invoice data from the synchronizer.
    *   `GET /invoices/{id}`: Retrieve specific invoice data (likely for the frontend).
    *   `PUT /invoices/{id}`: Update an invoice with customer's fiscal data provided via the frontend.
    *   `GET /invoices/{id}/pdf`: Generate and return the final invoice PDF with all data.
*   **Data Storage:**
    *   Design an optimized database schema for `Factura` and `FacturasLin` data.
*   **Security:**
    *   Implement secure authentication between the synchronizer and the API (e.g., API keys, JWT).
    *   Implement HTTPS.
    *   Validate and sanitize all inputs.
    *   Protect against common web vulnerabilities (SQLi, XSS, CSRF).
*   **Other Considerations:** Robust error handling, data validation, PDF generation library.

### 3. React Frontend

*   **Objective:** Provide a web interface for customers to view their initial invoice, enter their fiscal details, and download the final PDF invoice.
*   **Hosting:** Plesk server.
*   **Access:** Via unique URLs embedded in QR codes.
*   **User Interface:**
    *   Display initial invoice details (retrieved from the API).
    *   Provide a form for customers to input fiscal data (Name, Address, NIF/CIF, etc.).
    *   Allow download of the final PDF invoice after submitting fiscal data.
*   **Technical Considerations:**
    *   Responsive design (mobile-first).
    *   Client-side form validation.
    *   Routing using React Router.
    *   State management (Context API, Redux, Zustand, etc.).
    *   Intuitive and user-friendly design.
    *   Secure communication with the API (HTTPS).

### 4. Overall System

*   **Workflow:** Follows the 8 steps outlined: Request -> Mark QR -> Generate/Print Invoice -> Sync -> Generate/Print QR -> Scan QR -> Enter Fiscal Data -> Download PDF.
*   **Security:** HTTPS everywhere, input validation/sanitization, vulnerability prevention (XSS, CSRF, SQLi), unique/secure invoice links.
*   **Non-Functional Requirements:**
    *   High degree of automation (minimal staff intervention).
    *   Scalability for increased invoice volume.
    *   Stability and robustness, especially on the Windows 7 synchronizer component.
    *   Consider potential connectivity/performance limitations.
*   **Deliverables:**
    *   Source Code: Go Synchronizer, REST API, React Frontend.
    *   Documentation: Technical details for each component.
    *   Manuals: Installation, configuration, and Plesk deployment instructions.
*   **Priorities:** The Go synchronizer is the critical first component due to its integration with the existing TPV system on Windows 7.

