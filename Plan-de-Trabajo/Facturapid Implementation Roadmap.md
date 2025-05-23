# Facturapid Implementation Roadmap

This document outlines a suggested implementation roadmap with estimated timelines and key milestones for the Facturapid project. Timelines are estimates and may need adjustment based on development progress and unforeseen challenges.

**Overall Priority:** Phase 1 (Go Synchronizer) is the highest priority due to its critical integration role and dependency on the Windows 7 environment.

## Phase 1: Go Synchronizer Development (Estimated Duration: 3-5 Weeks)

*   **Week 1: Setup & Core Connectivity**
    *   Milestone 1.1: Development and Windows 7 test environments configured.
    *   Milestone 1.2: Successful connection to the MS Access database (`.mdb`) via ODBC established from Go.
    *   Milestone 1.3: Basic polling mechanism implemented to query the `Factura` table.
*   **Week 2: Data Processing & API Stub Integration**
    *   Milestone 2.1: Logic for extracting `Factura` and `FacturasLin` data completed.
    *   Milestone 2.2: Data structures (structs) for API communication defined.
    *   Milestone 2.3: Basic HTTP client implemented to send data to a *mock* API endpoint (for testing synchronizer logic independently).
*   **Week 3: QR Generation & Printing**
    *   Milestone 3.1: QR code generation library integrated; URL construction logic complete.
    *   Milestone 3.2: Initial implementation of LPT1 printing logic (ESC/POS commands for QR code).
    *   Milestone 3.3: Testing QR code generation and printing on the target Epson TM-T88 printer.
*   **Week 4: Service Implementation & Refinement**
    *   Milestone 4.1: Integration with `kardianos/service` for Windows service functionality.
    *   Milestone 4.2: Installation/uninstallation scripts created.
    *   Milestone 4.3: Configuration file implementation and logging refined.
*   **Week 5: Integration Testing (Synchronizer -> Real API) & Buffer**
    *   Milestone 5.1: Replace mock API calls with calls to the actual (partially developed) REST API `/invoices` endpoint.
    *   Milestone 5.2: Implement API Key authentication.
    *   Milestone 5.3: Conduct initial integration tests in the Windows 7 environment.
    *   Buffer time for addressing issues found during testing, especially related to ODBC or LPT1 interaction.

## Phase 2: REST API Development (Estimated Duration: 2-3 Weeks - Can partially overlap with Phase 1)

*   **Week 3 (Overlap): Core Endpoints & Database**
    *   Milestone 6.1: Technology stack (Go framework, DB) chosen and set up.
    *   Milestone 6.2: Database schema designed and implemented (tables `invoices`, `invoice_lines`).
    *   Milestone 6.3: `POST /invoices` endpoint implemented and tested (receives data, stores in DB, returns ID). API Key authentication implemented.
*   **Week 4 (Overlap): Remaining Endpoints**
    *   Milestone 7.1: `GET /invoices/{id}` endpoint implemented and tested.
    *   Milestone 7.2: `PUT /invoices/{id}` endpoint implemented and tested (updates fiscal data).
*   **Week 5: PDF Generation & Security**
    *   Milestone 8.1: PDF generation library integrated.
    *   Milestone 8.2: `GET /invoices/{id}/pdf` endpoint implemented and tested.
    *   Milestone 8.3: Implement HTTPS, review input validation and security measures.

## Phase 3: React Frontend Development (Estimated Duration: 2-3 Weeks - Can partially overlap with Phase 2)

*   **Week 5 (Overlap): Setup & Basic View**
    *   Milestone 9.1: React project setup (Vite/CRA), routing configured.
    *   Milestone 9.2: API client functions created.
    *   Milestone 9.3: `/invoice/:id` page created; fetches and displays initial invoice data from `GET /invoices/{id}`.
*   **Week 6: Fiscal Data Form & Submission**
    *   Milestone 10.1: Fiscal data input form component built.
    *   Milestone 10.2: Client-side validation implemented.
    *   Milestone 10.3: Form submission logic implemented (calls `PUT /invoices/{id}`).
*   **Week 7: PDF Download & Styling/Responsiveness**
    *   Milestone 11.1: PDF download functionality implemented (links to `GET /invoices/{id}/pdf`).
    *   Milestone 11.2: UI styling and responsiveness finalized.
    *   Milestone 11.3: Build process configured for deployment.

## Phase 4: Integration, Testing, Deployment & Documentation (Estimated Duration: 1-2 Weeks)

*   **Week 8: End-to-End Testing & Bug Fixing**
    *   Milestone 12.1: Full system integration in a staging environment (Win 7 Sync, API Server, DB, Plesk Frontend).
    *   Milestone 12.2: Execute end-to-end test cases covering the entire workflow.
    *   Milestone 12.3: Identify and fix bugs found during integration testing.
*   **Week 9: Deployment & Documentation**
    *   Milestone 13.1: Final security review.
    *   Milestone 13.2: Deploy API & DB to production server.
    *   Milestone 13.3: Deploy React frontend build to Plesk.
    *   Milestone 13.4: Install and configure Go Synchronizer service on production Windows 7 machine(s).
    *   Milestone 13.5: Finalize all documentation (technical docs, manuals).
    *   Milestone 13.6: User Acceptance Testing (UAT) and sign-off.
    *   Milestone 13.7: Project handover.

**Total Estimated Duration:** 9 Weeks (assuming some parallel work)

