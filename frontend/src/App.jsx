import { useState } from 'react'
import { Icon } from '@iconify/react'
import { Routes, Route, Navigate, Link, useLocation } from 'react-router-dom'
import Settings from './Settings'
import './App.css'

const NAV = [{ id: 'general', path: '/settings/general', label: 'General', icon: 'mdi:cog' }]

export default function App() {
  const [open, setOpen] = useState(true)
  const location = useLocation()

  return (
    <div className="shell">
      <aside className={`sidebar ${open ? '' : 'collapsed'}`}>
        <button className="sidebar-toggle" onClick={() => setOpen((v) => !v)}>
          <Icon icon={open ? 'mdi:chevron-left' : 'mdi:chevron-right'} width={16} />
        </button>
        <nav className="sidebar-nav">
          {NAV.map((item) => (
            <Link
              key={item.id}
              to={item.path}
              className={location.pathname === item.path ? 'active' : ''}
            >
              <Icon icon={item.icon} width={16} />
              {open && <span>{item.label}</span>}
            </Link>
          ))}
        </nav>
      </aside>

      <main className="content">
        <Routes>
          <Route path="/settings/general" element={<Settings />} />
          <Route path="*" element={<Navigate to="/settings/general" replace />} />
        </Routes>
      </main>
    </div>
  )
}
