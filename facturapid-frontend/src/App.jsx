import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import InvoicePage from './pages/InvoicePage'; // Assuming InvoicePage.jsx is in src/pages

function App() {
  return (
    <>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/invoice/123">Sample Invoice 123</Link> {/* Example link */}
          </li>
          <li>
            <Link to="/invoice/456">Sample Invoice 456</Link> {/* Example link */}
          </li>
        </ul>
      </nav>

      <hr />

      {/* Route definitions will be in main.jsx as per typical Vite setup, 
          but if we want to nest routes within App layout, they could be here.
          For this task, main.jsx will handle the top-level Routes.
          This App component will just define the layout and common navigation.
          The <Outlet /> component would be used here if this was a layout route.
          However, the requirement is to set up routes in main.jsx or App.jsx.
          Let's assume main.jsx handles the BrowserRouter and Routes, and App.jsx
          is just a component that might be rendered by one of those routes, or provides layout.

          Revisiting the requirement: "Modify App.jsx (or main.jsx if using Vite's default structure more directly for routing)"
          Vite's default `main.jsx` renders `<App />`. It's common to put `BrowserRouter` in `main.jsx`
          and then the primary `Routes` in `App.jsx`. I'll follow this common pattern.
      */}

      <Routes>
        <Route 
          path="/" 
          element={
            <div>
              <h2>Welcome to Facturapid Frontend!</h2>
              <p>This is the home page.</p>
            </div>
          } 
        />
        <Route path="/invoice/:invoiceId" element={<InvoicePage />} />
        <Route path="*" element={<div><h2>Page Not Found</h2><p>Sorry, the page you are looking for does not exist.</p></div>} />
      </Routes>
    </>
  );
}

export default App;
