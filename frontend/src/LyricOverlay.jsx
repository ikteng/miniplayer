import { useEffect, useRef, useState } from 'react'
import { getDominantColor } from './coverColor.js'
import './lyric-overlay.css'

const EMPTY_TRACK = { title: '', author: '', duration_human: '0:00' }
const EMPTY_PLAYER = { has_song: false, is_paused: true, seekbar_current_position: 0 }

// Parses yrc_lyric lines of the form "[<ms>,0]text" into {ms, text} pairs.
function parseYrcLyric(yrc) {
  if (!yrc) return []
  const lineRe = /^\[(\d+),\d+\](.*)$/
  return yrc
    .split('\n')
    .map((raw) => {
      const m = lineRe.exec(raw)
      if (!m) return null
      return { ms: Number(m[1]), text: m[2] }
    })
    .filter(Boolean)
}

function parsePlainLyric(plain) {
  if (!plain) return []
  return plain
    .split('\n')
    .filter((text) => text.trim() !== '')
    .map((text) => ({ ms: null, text }))
}

// status: lets an embedding page (e.g. the Lyrics settings preview) show
// why nothing is rendering, instead of silently returning null like the
// bare OBS overlay does.
export default function LyricOverlay({ onStatusChange } = {}) {
  const [player, setPlayer] = useState(EMPTY_PLAYER)
  const [track, setTrack] = useState(EMPTY_TRACK)
  const [lyricStyle, setLyricStyle] = useState('classic')
  const [coverTint, setCoverTint] = useState(null)
  const [lines, setLines] = useState([])
  const lastKeyRef = useRef('')

  useEffect(() => {
    fetch('/widgets/amuse/settings')
      .then((res) => res.json())
      .then((settings) => {
        setLyricStyle(settings.lyric_style || 'classic')
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
        if (cancelled) return
        setPlayer(player)
        setTrack(track)

        const key = `${track.title}_${track.author}`
        if (track.has_song !== false && track.title && key !== lastKeyRef.current) {
          lastKeyRef.current = key
          try {
            const lyricRes = await fetch(
              `/api/lyric?title=${encodeURIComponent(track.title)}&author=${encodeURIComponent(track.author)}`,
            )
            const { data } = await lyricRes.json()
            if (cancelled) return
            if (data?.has_karaoke_lyric) {
              setLines(parseYrcLyric(data.yrc_lyric))
            } else if (data?.has_lyric) {
              setLines(parsePlainLyric(data.lyric))
            } else {
              setLines([])
            }
          } catch {
            setLines([])
          }
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

  useEffect(() => {
    if (!onStatusChange) return
    if (!player.has_song) onStatusChange('no-song')
    else if (lines.length === 0) onStatusChange('no-lyrics')
    else onStatusChange('ok')
  }, [player.has_song, lines.length, onStatusChange])

  if (!player.has_song) return null
  if (lines.length === 0) return null

  const positionMs = player.seekbar_current_position
  const timed = lines[0]?.ms !== null

  let activeIndex = 0
  if (timed) {
    for (let i = 0; i < lines.length; i++) {
      if (lines[i].ms <= positionMs) activeIndex = i
      else break
    }
  }

  const prev = lines[activeIndex - 1]
  const current = lines[activeIndex]
  const next = lines[activeIndex + 1]

  if (lyricStyle === 'bubbles') {
    return (
      <div className="lyric-overlay lyric-overlay--bubbles" style={{ '--tint': coverTint || '#40a4d8' }}>
        {prev && <div key={activeIndex - 1} className="lyric-bubble">{prev.text}</div>}
        {current && (
          <div key={activeIndex} className="lyric-bubble lyric-bubble--active">
            {current.text}
          </div>
        )}
      </div>
    )
  }

  return (
    <div className="lyric-overlay" style={{ '--tint': coverTint || '#40a4d8' }}>
      {prev && <div key={activeIndex - 1} className="lyric-line">{prev.text}</div>}
      {current && <div key={activeIndex} className="lyric-line active">{current.text}</div>}
      {next && <div key={activeIndex + 1} className="lyric-line">{next.text}</div>}
    </div>
  )
}
