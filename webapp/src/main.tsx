import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

// Catch global errors for debugging
window.onerror = function(message, _source, _lineno, _colno, error) {
  const root = document.getElementById('root');
  if (root) {
    root.innerHTML = `
      <div style="background: #111; color: red; padding: 20px; font-family: monospace; height: 100vh;">
        <h1 style="font-size: 14px;">CRITICAL_ERROR</h1>
        <p style="color: #666; font-size: 12px;">${message}</p>
        <pre style="font-size: 10px; color: #444;">${error?.stack || ''}</pre>
      </div>
    `;
  }
  return false;
};

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
