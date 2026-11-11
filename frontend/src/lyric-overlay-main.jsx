import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './lyric-overlay-page.css'
import LyricOverlay from './LyricOverlay.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <LyricOverlay />
  </StrictMode>,
)
