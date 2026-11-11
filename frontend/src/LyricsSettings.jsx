import { useEffect, useRef, useState } from 'react'
import { Icon } from '@iconify/react'
import LyricOverlay from './LyricOverlay.jsx'

const STATUS_TEXT = {
  'no-song': 'No song is currently playing',
  'no-lyrics': 'No lyrics found for the current track',
  ok: null,
}

const STYLES = [
  { id: 'classic', label: 'Classic' },
  { id: 'bubbles', label: 'Bubbles' },
]

const DEFAULT_POS = { right: 24, bottom: 24 }

export default function LyricsSettings() {
  const [status, setStatus] = useState('no-song')
  const [lyricStyle, setLyricStyle] = useState(null)
  const [saveStatus, setSaveStatus] = useState('')
  const [previewKey, setPreviewKey] = useState(0)
  const [copied, setCopied] = useState(false)
  const [pos, setPos] = useState(DEFAULT_POS)
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
    const url = `${window.location.origin}/lyric-overlay`
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
      .then((settings) => setLyricStyle(settings.lyric_style || 'classic'))
      .catch(() => {})
  }, [])

  async function selectStyle(id) {
    setLyricStyle(id)
    setPreviewKey((k) => k + 1)
    setSaveStatus('Saving...')
    try {
      const res = await fetch('/widgets/amuse/settings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lyric_style: id }),
      })
      const json = await res.json()
      setSaveStatus(json.code === 200 ? 'Saved' : `Error: ${json.msg}`)
    } catch {
      setSaveStatus('Failed to save')
    }
    setTimeout(() => setSaveStatus(''), 2000)
  }

  return (
    <main className="settings-page">
      <div
        ref={dragRef}
        className="floating-preview"
        style={{ right: pos.right, bottom: pos.bottom }}
      >
        <div className="floating-preview-handle" onPointerDown={startDrag}>
          <Icon icon="mdi:drag" width={14} />
          <span>Preview</span>
        </div>
        <div className="lyric-preview-frame">
          <LyricOverlay key={previewKey} onStatusChange={setStatus} />
          {STATUS_TEXT[status] && <div className="lyric-preview-placeholder">{STATUS_TEXT[status]}</div>}
        </div>
      </div>

      <section className="settings-section">
        <h2 className="section-title">Preview</h2>
        <p className="desc">
          This mirrors the OBS overlay at <code>/lyric-overlay</code>, showing whatever is
          currently playing. Drag the floating preview in the corner to move it around.
        </p>

        <div className="link-row">
          <input
            type="text"
            className="link-input"
            value={`${window.location.origin}/lyric-overlay`}
            readOnly
            onFocus={(e) => e.target.select()}
          />
          <button className="detect-btn" onClick={copyOverlayLink}>
            <Icon icon={copied ? 'mdi:check' : 'mdi:content-copy'} width={14} />
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>
      </section>

      {lyricStyle && (
        <section className="settings-section">
          <h2 className="section-title">Style</h2>

          <div className="platform-tabs">
            {STYLES.map((s) => (
              <button
                key={s.id}
                className={lyricStyle === s.id ? 'active' : ''}
                onClick={() => selectStyle(s.id)}
              >
                {s.label}
              </button>
            ))}
          </div>

          {saveStatus && <div className="status">{saveStatus}</div>}
        </section>
      )}
    </main>
  )
}
