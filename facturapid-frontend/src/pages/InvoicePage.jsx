import React from 'react';
import { useParams } from 'react-router-dom';

function InvoicePage() {
  const { invoiceId } = useParams();

  return (
    <div>
      <h1>Invoice Details</h1>
      <p>Displaying invoice details for ID: <strong>{invoiceId}</strong></p>
      {/* Placeholder for actual invoice data fetching and display */}
    </div>
  );
}

export default InvoicePage;
