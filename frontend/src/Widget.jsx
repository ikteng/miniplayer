import { useEffect, useState } from 'react'

const EMPTY_TRACK = { title: '', author: '', cover: '', duration_human: '0:00' }
const EMPTY_PLAYER = { has_song: false, is_paused: true, seekbar_current_position_human: '0:00', state_percent: 0 }

export default function Widget() {
  const [player, setPlayer] = useState(EMPTY_PLAYER)
  const [track, setTrack] = useState(EMPTY_TRACK)
  const [tintColor, setTintColor] = useState('#40a4d8')
  const [hideOnPause, setHideOnPause] = useState(false)

  useEffect(() => {
    fetch('/widgets/amuse/settings')
      .then((res) => res.json())
      .then(({ settings }) => {
        setTintColor(settings.tint_color)
        setHideOnPause(settings.hide_on_pause)
      })
      .catch(() => {})
  }, [])

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
  if (!visible) return null

  return (
    <div className="widget-card" style={{ '--tint': tintColor }}>
      <div className="widget-cover">
        {track.cover && <img src={track.cover} alt="" />}
      </div>
      <div className="widget-meta">
        <div className="widget-title">{track.title || 'Unknown title'}</div>
        <div className="widget-author">{track.author || 'Unknown artist'}</div>
        <div className="widget-bar">
          <div className="widget-bar-fill" style={{ width: `${Math.min(100, Math.max(0, player.state_percent * 100))}%` }} />
        </div>
      </div>
    </div>
  )
}
