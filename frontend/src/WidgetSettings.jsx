import { useEffect, useRef, useState } from 'react'
import { Icon } from '@iconify/react'
import Widget from './Widget.jsx'

const STATUS_TEXT = {
  'no-song': 'No song is currently playing',
  'hidden-on-pause': 'Hidden — "Hide on pause" is on and playback is paused',
  ok: null,
}

const DEFAULT_POS = { right: 24, bottom: 24 }

const STYLES = [
  { id: 'compact', label: 'Compact' },
  { id: 'album', label: 'Album' },
  { id: 'minimal', label: 'Minimal' },
]

export default function WidgetSettings() {
  const [settings, setSettings] = useState(null)
  const [general, setGeneral] = useState(null)
  const [status, setStatus] = useState('no-song')
  const [saveStatus, setSaveStatus] = useState('')
  const [copied, setCopied] = useState(false)
  const [pos, setPos] = useState(DEFAULT_POS)
  const [previewKey, setPreviewKey] = useState(0)
  const dragRef = useRef(null)

  function startDrag(e) {
    const panel = dragRef.current
    if (!panel) return
    const rect = panel.getBoundingClientRect()
    const offsetX = e.clientX - rect.left
    const offsetY = e.clientY - rect.top

    function onMove(ev) {
      const right = window.innerWidth - ev.clientX - (rect.width - offsetX)
      const bottom = window.innerHeight - ev.clientY - (rect.height - offsetY)
      const maxRight = Math.max(8, window.innerWidth - rect.width - 8)
      const maxBottom = Math.max(8, window.innerHeight - rect.height - 8)
      setPos({
        right: Math.min(maxRight, Math.max(8, right)),
        bottom: Math.min(maxBottom, Math.max(8, bottom)),
      })
    }
    function onUp() {
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerup', onUp)
    }
    window.addEventListener('pointermove', onMove)
    window.addEventListener('pointerup', onUp)
  }

  useEffect(() => {
    function clamp() {
      const panel = dragRef.current
      if (!panel) return
      const rect = panel.getBoundingClientRect()
      const maxRight = Math.max(8, window.innerWidth - rect.width - 8)
      const maxBottom = Math.max(8, window.innerHeight - rect.height - 8)
      setPos((p) => ({
        right: Math.min(p.right, maxRight),
        bottom: Math.min(p.bottom, maxBottom),
      }))
    }
    window.addEventListener('resize', clamp)
    return () => window.removeEventListener('resize', clamp)
  }, [])

  async function copyOverlayLink() {
    const url = `${window.location.origin}/widget`
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // clipboard access denied or unavailable
    }
  }

  useEffect(() => {
    fetch('/widgets/amuse/settings')
      .then((res) => res.json())
      .then(setSettings)
      .catch(() => {})
    fetch('/api/settings/general')
      .then((res) => res.json())
      .then(setGeneral)
      .catch(() => {})
  }, [])

  async function selectStyle(id) {
    await save({ ...settings, widget_style: id })
    setPreviewKey((k) => k + 1)
  }

  async function save(next) {
    setSettings(next)
    setSaveStatus('Saving...')
    try {
      const res = await fetch('/widgets/amuse/settings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(next),
      })
      const json = await res.json()
      setSaveStatus(json.code === 200 ? 'Saved' : `Error: ${json.msg}`)
    } catch {
      setSaveStatus('Failed to save')
    }
    setTimeout(() => setSaveStatus(''), 2000)
  }

  async function saveGeneral(next) {
    setGeneral(next)
    try {
      await fetch('/api/settings/general', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(next),
      })
    } catch {
      // ignore transient save error
    }
  }

  if (!settings || !general) return <main className="settings-page">Loading...</main>

  return (
    <main className="settings-page">
      <div
        ref={dragRef}
        className="floating-preview"
        style={{ right: pos.right, bottom: pos.bottom }}
        onPointerDown={startDrag}
      >
        <Widget key={previewKey} onStatusChange={setStatus} />
        {STATUS_TEXT[status] && <div className="lyric-preview-placeholder">{STATUS_TEXT[status]}</div>}
      </div>

      <section className="settings-section">
        <h2 className="section-title">Preview</h2>
        <p className="desc">
          This mirrors the OBS overlay at <code>/widget</code>, showing whatever is currently
          playing. Drag the floating preview in the corner to move it around.
        </p>

        <div className="link-row">
          <input
            type="text"
            className="link-input"
            value={`${window.location.origin}/widget`}
            readOnly
            onFocus={(e) => e.target.select()}
          />
          <button className="detect-btn" onClick={copyOverlayLink}>
            <Icon icon={copied ? 'mdi:check' : 'mdi:content-copy'} width={14} />
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>
      </section>

      <section className="settings-section">
        <h2 className="section-title">Style</h2>

        <div className="style-swatches">
          {STYLES.map((s) => (
            <button
              key={s.id}
              className={`style-swatch style-swatch--${s.id} ${(settings.widget_style || 'compact') === s.id ? 'active' : ''}`}
              onClick={() => selectStyle(s.id)}
            >
              <span className="style-swatch-mock" aria-hidden="true">
                {s.id === 'compact' && (
                  <span className="mock-compact">
                    <span className="mock-cover" />
                    <span className="mock-lines">
                      <span className="mock-line mock-line--title" />
                      <span className="mock-line mock-line--author" />
                      <span className="mock-line mock-line--bar" />
                    </span>
                  </span>
                )}
                {s.id === 'album' && (
                  <span className="mock-album">
                    <span className="mock-album-shade" />
                    <span className="mock-lines">
                      <span className="mock-line mock-line--title" />
                      <span className="mock-line mock-line--bar" />
                    </span>
                  </span>
                )}
                {s.id === 'minimal' && (
                  <span className="mock-minimal">
                    <span className="mock-line mock-line--title" />
                    <span className="mock-line mock-line--bar" />
                  </span>
                )}
              </span>
              <span className="style-swatch-label">{s.label}</span>
            </button>
          ))}
        </div>
      </section>

      <section className="settings-section">
        <h2 className="section-title">Appearance</h2>

        <div className="toggle-row">
          <div className="toggle-text">
            <span>Hide on pause</span>
            <span className="desc">Hide the widget while playback is paused</span>
          </div>
          <input
            type="checkbox"
            className="toggle"
            checked={settings.hide_on_pause}
            onChange={(e) => save({ ...settings, hide_on_pause: e.target.checked })}
          />
        </div>

        {saveStatus && <div className="status">{saveStatus}</div>}
      </section>

      <section className="settings-section">
        <h2 className="section-title">System</h2>

        <div className="toggle-row">
          <div className="toggle-text">
            <span>Run at startup</span>
            <span className="desc">Launch MiniPlayer automatically when the system starts</span>
          </div>
          <input
            type="checkbox"
            className="toggle"
            checked={general.run_at_startup}
            onChange={(e) => saveGeneral({ ...general, run_at_startup: e.target.checked })}
          />
        </div>
      </section>
    </main>
  )
}
