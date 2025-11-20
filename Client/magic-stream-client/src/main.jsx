import "bootstrap/dist/css/bootstrap.min.css";
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from "react-router-dom"; // ← ADD THIS
import App from './App.jsx';
import './index.css';

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>              {/* ← WRAP APP HERE */}
      <App />
    </BrowserRouter>
  </StrictMode>,
);

