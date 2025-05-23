// src/services/apiService.js

const API_BASE_URL = 'http://localhost:8080/api/v1'; // Ensure this matches your backend API
const API_KEY = 'supersecretapikey'; // Replace with your actual API key, consider using environment variables

/**
 * Helper function to handle common fetch responses and errors.
 * @param {Response} response - The response object from a fetch call.
 * @returns {Promise<any>} - A promise that resolves with the JSON data if response is ok.
 * @throws {Error} - Throws an error if the response is not ok.
 */
async function handleResponse(response) {
  if (!response.ok) {
    let errorMessage = `API request failed with status ${response.status}`;
    try {
      const errorData = await response.json();
      if (errorData && errorData.error) {
        errorMessage += `: ${errorData.error}`;
        if (errorData.details) {
          errorMessage += ` (${errorData.details})`;
        }
      } else if (response.statusText) {
        errorMessage += `: ${response.statusText}`;
      }
    } catch (e) {
      // Could not parse error JSON, use default status text if available
      if (response.statusText) {
        errorMessage += `: ${response.statusText}`;
      }
    }
    throw new Error(errorMessage);
  }
  // If response is ok, try to parse JSON.
  // Handle cases where response might be ok but have no content (e.g., 204 No Content)
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.indexOf("application/json") !== -1) {
    return response.json();
  }
  return response.text(); // Or handle as text, or return null/undefined for no content
}

/**
 * Fetches a single invoice by its ID.
 * @param {string|number} invoiceId - The ID of the invoice to fetch.
 * @returns {Promise<object>} - A promise that resolves to the invoice data.
 * @throws {Error} - Throws an error if the request fails.
 */
export async function fetchInvoice(invoiceId) {
  const url = `${API_BASE_URL}/invoices/${invoiceId}`;
  console.log(`Fetching invoice from: ${url}`); // For debugging

  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'X-API-Key': API_KEY,
      },
    });
    return handleResponse(response);
  } catch (error) {
    console.error('Error in fetchInvoice:', error);
    throw error; // Re-throw to be caught by the calling component
  }
}

/**
 * Submits updated fiscal data for an invoice.
 * @param {string|number} invoiceId - The ID of the invoice to update.
 * @param {object} fiscalData - The fiscal data object to submit.
 *   Example: { cliente_nombre: "New Name", cliente_nif: "New NIF", ... }
 * @returns {Promise<object>} - A promise that resolves to the response data from the server (if any).
 * @throws {Error} - Throws an error if the request fails.
 */
export async function submitFiscalData(invoiceId, fiscalData) {
  const url = `${API_BASE_URL}/invoices/${invoiceId}`;
  console.log(`Submitting fiscal data to: ${url}`, fiscalData); // For debugging

  try {
    const response = await fetch(url, {
      method: 'PUT',
      headers: {
        'X-API-Key': API_KEY,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(fiscalData),
    });
    return handleResponse(response);
  } catch (error) {
    console.error('Error in submitFiscalData:', error);
    throw error; // Re-throw
  }
}

/**
 * Constructs the URL to fetch the PDF for a given invoice ID.
 * This URL can be used with an <a> tag or window.open().
 * Note: This function itself doesn't make an API call.
 * For downloading via fetch and handling blobs, a different approach would be needed.
 * @param {string|number} invoiceId - The ID of the invoice.
 * @returns {string} - The URL to the invoice PDF.
 */
export function getInvoicePdfUrl(invoiceId) {
  const url = `${API_BASE_URL}/invoices/${invoiceId}/pdf`;
  console.log(`Generated PDF URL: ${url}`); // For debugging
  return url;
  // To make this URL work directly in <a> or window.open, the browser needs to handle the API key.
  // If the API key cannot be sent via query param (which is bad practice for keys),
  // then fetching the PDF as a blob and creating an object URL is the way to go for client-side download.
  // However, the current requirement is just to return the URL string.
  // If the PDF endpoint is protected by X-API-Key header, this URL won't work directly in an <a> tag
  // without further client-side logic to fetch the blob.
}


/**
 * Fetches the invoice PDF as a Blob.
 * This is an alternative for getInvoicePdfUrl if direct linking is not feasible due to auth headers.
 * @param {string|number} invoiceId - The ID of the invoice.
 * @returns {Promise<Blob>} - A promise that resolves to the PDF Blob.
 * @throws {Error} - Throws an error if the request fails.
 */
export async function fetchInvoicePdfBlob(invoiceId) {
    const url = `${API_BASE_URL}/invoices/${invoiceId}/pdf`;
    console.log(`Fetching PDF blob from: ${url}`);

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'X-API-Key': API_KEY,
            },
        });

        if (!response.ok) {
            // Try to parse error message if backend sends JSON error for PDF endpoint
            let errorMessage = `Failed to fetch PDF with status ${response.status}`;
            try {
                const errorData = await response.json(); // Assuming error might be JSON
                if (errorData && errorData.error) {
                    errorMessage += `: ${errorData.error}`;
                }
            } catch (e) {
                // Ignore if error response is not JSON
            }
            throw new Error(errorMessage);
        }
        return response.blob();
    } catch (error) {
        console.error('Error in fetchInvoicePdfBlob:', error);
        throw error;
    }
}

// Example of how to use fetchInvoicePdfBlob and trigger download:
/*
async function downloadInvoicePdf(invoiceId) {
  try {
    const blob = await fetchInvoicePdfBlob(invoiceId);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = `factura_${invoiceId}.pdf`; // Filename for download
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    a.remove();
  } catch (error) {
    console.error("Failed to download PDF:", error);
    // Display error to user
  }
}
*/
