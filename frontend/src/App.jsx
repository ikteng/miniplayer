import { useEffect, useState } from 'react'
import { Icon } from '@iconify/react'
import { Routes, Route, Navigate, Link, useLocation } from 'react-router-dom'
import WidgetSettings from './WidgetSettings'
import LyricsSettings from './LyricsSettings'
import { getAccentPair } from './coverColor.js'
import './App.css'

const NAV = [
  { id: 'widget', path: '/settings/widget', label: 'Widget', icon: 'mdi:picture-in-picture-bottom-right-outline' },
  { id: 'lyrics', path: '/settings/lyrics', label: 'Lyrics', icon: 'mdi:music-note-outline' },
]

export default function App() {
  const [open, setOpen] = useState(true)
  const location = useLocation()

  const current = NAV.find((item) => item.path === location.pathname)

  useEffect(() => {
    let cancelled = false
    let lastCover = null

    async function poll() {
      try {
        const res = await fetch('/api/query')
        const { track } = await res.json()
        const cover = track?.cover || null
        if (cover === lastCover) return
        lastCover = cover

        if (!cover) {
          document.documentElement.style.removeProperty('--accent')
          document.documentElement.style.removeProperty('--accent-2')
          return
        }
        const pair = await getAccentPair(cover)
        if (cancelled || !pair) return
        document.documentElement.style.setProperty('--accent', pair.accent)
        document.documentElement.style.setProperty('--accent-2', pair.accent2)
      } catch {
        // keep last known accent on a transient fetch error
      }
    }

    poll()
    const id = setInterval(poll, 2000)
    return () => {
      cancelled = true
      clearInterval(id)
    }
  }, [])

  return (
    <div className="shell">
      <aside className={`sidebar ${open ? '' : 'collapsed'}`}>
        <div className="brand">
          <div className="brand-mark">
            <Icon icon="mdi:play-circle" width={20} />
          </div>
          {open && <span className="brand-name">MiniPlayer</span>}
          <button className="sidebar-toggle" onClick={() => setOpen((v) => !v)}>
            <Icon icon={open ? 'mdi:chevron-left' : 'mdi:chevron-right'} width={16} />
          </button>
        </div>
        <nav className="sidebar-nav">
          {NAV.map((item) => (
            <Link
              key={item.id}
              to={item.path}
              className={location.pathname === item.path ? 'active' : ''}
              title={item.label}
            >
              <Icon icon={item.icon} width={17} />
              <span className={open ? '' : 'nav-label-collapsed'}>{item.label}</span>
            </Link>
          ))}
        </nav>
      </aside>

      <main className="content">
        <header className="topbar">
          <h1 className="topbar-title">{current?.label ?? ''}</h1>
        </header>
        <div className="content-body">
          <Routes location={location}>
            <Route path="/settings/widget" element={<WidgetSettings />} />
            <Route path="/settings/lyrics" element={<LyricsSettings />} />
            <Route path="*" element={<Navigate to="/settings/widget" replace />} />
          </Routes>
        </div>
      </main>
    </div>
  )
}
