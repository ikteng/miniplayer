import { useEffect, useState } from 'react'
import { getDominantColor } from './coverColor.js'
import './widget.css'

const EMPTY_TRACK = { title: '', author: '', cover: '', duration_human: '0:00' }
const EMPTY_PLAYER = { has_song: false, is_paused: true, seekbar_current_position_human: '0:00', state_percent: 0 }

export default function Widget({ onStatusChange } = {}) {
  const [player, setPlayer] = useState(EMPTY_PLAYER)
  const [track, setTrack] = useState(EMPTY_TRACK)
  const [hideOnPause, setHideOnPause] = useState(false)
  const [widgetStyle, setWidgetStyle] = useState('compact')
  const [coverTint, setCoverTint] = useState(null)

  useEffect(() => {
    fetch('/widgets/amuse/settings')
      .then((res) => res.json())
      .then((settings) => {
        setHideOnPause(settings.hide_on_pause)
        setWidgetStyle(settings.widget_style || 'compact')
      })
      .catch(() => {})
  }, [])

  useEffect(() => {
    if (!track.cover) {
      setCoverTint(null)
      return
    }
    let cancelled = false
    getDominantColor(track.cover).then((color) => {
      if (!cancelled) setCoverTint(color)
    })
    return () => {
      cancelled = true
    }
  }, [track.cover])

  useEffect(() => {
    let cancelled = false
    async function poll() {
      try {
        const res = await fetch('/api/query')
        const { player, track } = await res.json()
        if (!cancelled) {
          setPlayer(player)
          setTrack(track)
        }
      } catch {
        // keep last known state on a transient fetch error
      }
    }
    poll()
    const id = setInterval(poll, 1000)
    return () => {
      cancelled = true
      clearInterval(id)
    }
  }, [])

  const visible = player.has_song && !(hideOnPause && player.is_paused)

  useEffect(() => {
    if (!onStatusChange) return
    if (!player.has_song) onStatusChange('no-song')
    else if (!visible) onStatusChange('hidden-on-pause')
    else onStatusChange('ok')
  }, [player.has_song, visible, onStatusChange])

  if (!visible) return null

  const progress = Math.min(100, Math.max(0, player.state_percent * 100))
  const tint = coverTint || '#40a4d8'

  if (widgetStyle === 'album') {
    return (
      <div className="widget-card widget-card--album" style={{ '--tint': tint }}>
        <div className="widget-cover widget-cover--album">
          {track.cover && <img src={track.cover} alt="" />}
          <div className="widget-album-shade" />
          <div className="widget-album-meta">
            <div className="widget-title">{track.title || 'Unknown title'}</div>
            {track.author && <div className="widget-author">{track.author}</div>}
            <div className="widget-bar">
              <div className="widget-bar-fill" style={{ width: `${progress}%` }} />
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (widgetStyle === 'minimal') {
    return (
      <div className="widget-card widget-card--minimal" style={{ '--tint': tint }}>
        <div className="widget-minimal-text">
          <span className="widget-title">{track.title || 'Unknown title'}</span>
          {track.author && <span className="widget-minimal-sep">·</span>}
          {track.author && <span className="widget-author">{track.author}</span>}
        </div>
        <div className="widget-bar widget-bar--minimal">
          <div className="widget-bar-fill" style={{ width: `${progress}%` }} />
        </div>
      </div>
    )
  }

  return (
    <div className="widget-card" style={{ '--tint': tint }}>
      <div className="widget-cover">
        {track.cover && <img src={track.cover} alt="" />}
      </div>
      <div className="widget-meta">
        <div className="widget-title">{track.title || 'Unknown title'}</div>
        {track.author && <div className="widget-author">{track.author}</div>}
        <div className="widget-bar">
          <div className="widget-bar-fill" style={{ width: `${progress}%` }} />
        </div>
      </div>
    </div>
  )
}
