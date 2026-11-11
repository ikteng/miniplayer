// Extracts an approximate dominant color from a cover image by downscaling it
// onto a tiny canvas and averaging pixels, then darkens/saturates it slightly
// so it stays readable as a UI tint rather than a raw average (which tends
// toward washed-out gray for busy album art).
const cache = new Map()

export function getDominantColor(coverUrl) {
  if (!coverUrl) return Promise.resolve(null)
  if (cache.has(coverUrl)) return cache.get(coverUrl)

  const promise = new Promise((resolve) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      try {
        const size = 32
        const canvas = document.createElement('canvas')
        canvas.width = size
        canvas.height = size
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, size, size)
        const { data } = ctx.getImageData(0, 0, size, size)

        let r = 0, g = 0, b = 0, count = 0
        for (let i = 0; i < data.length; i += 4) {
          const alpha = data[i + 3]
          if (alpha < 32) continue
          r += data[i]
          g += data[i + 1]
          b += data[i + 2]
          count++
        }
        if (count === 0) return resolve(null)
        r = Math.round(r / count)
        g = Math.round(g / count)
        b = Math.round(b / count)

        resolve(boostColor(r, g, b))
      } catch {
        // tainted canvas (CORS) or decode failure — no tint available
        resolve(null)
      }
    }
    img.onerror = () => resolve(null)
    img.src = coverUrl
  })

  cache.set(coverUrl, promise)
  return promise
}

function boostColor(r, g, b) {
  const [h, s, l] = rgbToHsl(r, g, b)
  const [r2, g2, b2] = hslToRgb(h, Math.min(1, s * 1.4 + 0.15), Math.min(0.62, Math.max(0.38, l)))
  return `#${toHex(r2)}${toHex(g2)}${toHex(b2)}`
}

// Derives a UI accent pair (base + lighter variant) from a cover image,
// matching the same hue/saturation the widget tint uses.
export function getAccentPair(coverUrl) {
  return getDominantColor(coverUrl).then((hex) => {
    if (!hex) return null
    const [r, g, b] = [hex.slice(1, 3), hex.slice(3, 5), hex.slice(5, 7)].map((h) => parseInt(h, 16))
    const [h, s] = rgbToHsl(r, g, b)
    const [lr, lg, lb] = hslToRgb(h, s, 0.72)
    return { accent: hex, accent2: `#${toHex(lr)}${toHex(lg)}${toHex(lb)}` }
  })
}

function toHex(v) {
  return Math.round(v).toString(16).padStart(2, '0')
}

function rgbToHsl(r, g, b) {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h = 0, s = 0
  const l = (max + min) / 2
  const d = max - min
  if (d !== 0) {
    s = d / (1 - Math.abs(2 * l - 1))
    switch (max) {
      case r: h = ((g - b) / d) % 6; break
      case g: h = (b - r) / d + 2; break
      default: h = (r - g) / d + 4
    }
    h *= 60
    if (h < 0) h += 360
  }
  return [h, s, l]
}

function hslToRgb(h, s, l) {
  const c = (1 - Math.abs(2 * l - 1)) * s
  const x = c * (1 - Math.abs(((h / 60) % 2) - 1))
  const m = l - c / 2
  let r = 0, g = 0, b = 0
  if (h < 60) [r, g, b] = [c, x, 0]
  else if (h < 120) [r, g, b] = [x, c, 0]
  else if (h < 180) [r, g, b] = [0, c, x]
  else if (h < 240) [r, g, b] = [0, x, c]
  else if (h < 300) [r, g, b] = [x, 0, c]
  else [r, g, b] = [c, 0, x]
  return [(r + m) * 255, (g + m) * 255, (b + m) * 255]
}
