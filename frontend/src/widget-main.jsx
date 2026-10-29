import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './widget.css'
import Widget from './Widget.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <Widget />
  </StrictMode>,
)
