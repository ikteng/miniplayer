import { useEffect, useState } from 'react'
import { Icon } from '@iconify/react'

const PLATFORMS = {
  netease: 'NetEase Cloud Music',
  qq: 'QQ Music',
  kugou: 'Kugou Music',
  kuwo: 'Kuwo Music',
  spotify: 'Spotify',
  ayna: 'Ayna Record Player',
  apple: 'Apple Music',
  potplayer: 'PotPlayer',
  foobar: 'Foobar2000',
  lx: 'LX Music',
  soda: 'Soda Music',
  huahua: 'Huahua Live Assistant',
  musicfree: 'MusicFree',
  bq: 'BQ Jukebox',
  aimp: 'AIMP',
  youtube: 'YouTube Music',
  miebo: 'Miebo',
  yesplay: 'YesPlayMusic',
  cider: 'Cider',
  wesing: 'WeSing',
  browser: 'Browser',
  salt: 'Salt Player',
  'splayer-next': 'SPlayer Next',
}

const PLATFORM_GROUPS = {
  chinese: { label: 'Chinese', platforms: ['netease', 'qq', 'kugou', 'kuwo', 'soda', 'wesing'] },
  international: { label: 'International', platforms: ['spotify', 'apple', 'youtube', 'browser'] },
}

const MORE_PLATFORM_GROUPS = {
  jukebox: { label: 'Jukebox', platforms: ['miebo', 'ayna', 'huahua', 'bq'] },
  local: { label: 'Local player', platforms: ['potplayer', 'foobar', 'aimp', 'salt'] },
  other: { label: 'Other', platforms: ['lx', 'musicfree', 'cider', 'yesplay', 'splayer-next'] },
}

const POLL_INTERVALS = [100, 200, 300, 500, 1000]

const DEFAULT_DEVICE = { id: 'default', name: 'Default Audio Driver' }

export default function Settings() {
  const [form, setForm] = useState(null)
  const [status, setStatus] = useState('')
  const [devices, setDevices] = useState([DEFAULT_DEVICE])
  const [detecting, setDetecting] = useState(false)
  const [showMore, setShowMore] = useState(false)

  useEffect(() => {
    fetch('/api/settings/general')
      .then((res) => res.json())
      .then(setForm)
      .catch(() => setStatus('Failed to load settings'))

    fetch('/api/audio/devices')
      .then((res) => res.json())
      .then((json) => {
        const list = json.data ?? json
        if (Array.isArray(list) && list.length) setDevices(list)
      })
      .catch(() => {})
  }, [])

  function update(key, value) {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  async function save() {
    setStatus('Saving...')
    try {
      const res = await fetch('/api/settings/general', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      const json = await res.json()
      setStatus(json.code === 200 ? 'Saved' : `Error: ${json.msg}`)
    } catch {
      setStatus('Failed to save')
    }
    setTimeout(() => setStatus(''), 2000)
  }

  async function detectPlatform() {
    setDetecting(true)
    setStatus('Detecting...')
    try {
      const res = await fetch('/api/audio/detectPlatform')
      const json = await res.json()
      const detected = json.data?.platform
      if (detected) {
        update('platform', detected)
        setStatus(`Detected: ${PLATFORMS[detected] ?? detected}`)
      } else {
        setStatus('No active session detected')
      }
    } catch {
      setStatus('Detection failed')
    }
    setDetecting(false)
    setTimeout(() => setStatus(''), 2000)
  }

  if (!form) return <main className="settings-page">Loading...</main>

  const isMorePlatform = Object.values(MORE_PLATFORM_GROUPS).some((g) => g.platforms.includes(form.platform))
  const moreExpanded = showMore || isMorePlatform

  return (
    <main className="settings-page">
      <h1 className="page-title">General</h1>

      <section className="settings-section">
        <h2 className="section-title">Audio device</h2>
        <label className="field">
          <span>Device</span>
          <select
            value={form.device_id}
            onChange={(e) => {
              const device = devices.find((d) => d.id === e.target.value)
              setForm((prev) => ({ ...prev, device_id: e.target.value, device_name: device?.name ?? prev.device_name }))
            }}
          >
            {devices.map((d) => (
              <option key={d.id} value={d.id}>{d.name}</option>
            ))}
          </select>
        </label>
      </section>

      <section className="settings-section">
        <div className="section-title-row">
          <h2 className="section-title">Music service</h2>
          <button type="button" className="detect-btn" onClick={detectPlatform} disabled={detecting}>
            <Icon icon="mdi:magnify-scan" width={14} />
            {detecting ? 'Detecting...' : 'Detect automatically'}
          </button>
        </div>
        {Object.entries(PLATFORM_GROUPS).map(([groupId, group]) => (
          <div className="platform-group" key={groupId}>
            <span className="platform-group-label">{group.label}</span>
            <div className="platform-tabs">
              {group.platforms.map((id) => (
                <button
                  key={id}
                  type="button"
                  className={form.platform === id ? 'active' : ''}
                  onClick={() => update('platform', id)}
                >
                  {PLATFORMS[id]}
                </button>
              ))}
            </div>
          </div>
        ))}

        {moreExpanded && Object.entries(MORE_PLATFORM_GROUPS).map(([groupId, group]) => (
          <div className="platform-group" key={groupId}>
            <span className="platform-group-label">{group.label}</span>
            <div className="platform-tabs">
              {group.platforms.map((id) => (
                <button
                  key={id}
                  type="button"
                  className={form.platform === id ? 'active' : ''}
                  onClick={() => update('platform', id)}
                >
                  {PLATFORMS[id]}
                </button>
              ))}
            </div>
          </div>
        ))}

        <button type="button" className="show-more-btn" onClick={() => setShowMore((v) => !v)}>
          <Icon icon={moreExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'} width={14} />
          {moreExpanded ? 'Show less' : 'More'}
        </button>
      </section>

      <section className="settings-section">
        <h2 className="section-title">System settings</h2>

        <div className="toggle-row">
          <div className="toggle-text">
            <span>Run at startup</span>
            <span className="desc">Launch MiniPlayer automatically when the system starts</span>
          </div>
          <input
            type="checkbox"
            className="toggle"
            checked={form.run_at_startup}
            onChange={(e) => update('run_at_startup', e.target.checked)}
          />
        </div>

        <h3 className="subsection-title">Detection method</h3>

        <div className="toggle-row">
          <div className="toggle-text">
            <span>Prefer SMTC</span>
            <span className="desc">Keeps cover art identical to the one shown in the music app</span>
          </div>
          <input
            type="checkbox"
            className="toggle"
            checked={form.smtc}
            onChange={(e) => update('smtc', e.target.checked)}
          />
        </div>

        <div className="toggle-row">
          <div className="toggle-text">
            <span>Enable fallback platform</span>
            <span className="desc">Automatically detect a fallback platform when {PLATFORMS[form.platform] ?? 'the preferred platform'} isn't running</span>
          </div>
          <input
            type="checkbox"
            className="toggle"
            checked={form.fallback_platform_enabled}
            onChange={(e) => update('fallback_platform_enabled', e.target.checked)}
          />
        </div>

        {form.fallback_platform_enabled && (
          <label className="field">
            <span>Fallback platform</span>
            <select value={form.fallback_platform} onChange={(e) => update('fallback_platform', e.target.value)}>
              {Object.entries(PLATFORMS).map(([id, name]) => (
                <option key={id} value={id}>{name}</option>
              ))}
            </select>
          </label>
        )}

        <h3 className="subsection-title">Detection settings</h3>

        <label className="field">
          <span>Poll interval</span>
          <select value={form.poll_interval} onChange={(e) => update('poll_interval', Number(e.target.value))}>
            {POLL_INTERVALS.map((ms) => (
              <option key={ms} value={ms}>{ms} ms</option>
            ))}
          </select>
        </label>
      </section>

      <button className="save-btn" onClick={save}>
        <Icon icon="mdi:content-save" width={14} />
        Save
      </button>
      {status && <div className="status">{status}</div>}
    </main>
  )
}
