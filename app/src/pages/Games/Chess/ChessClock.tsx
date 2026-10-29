import { useEffect, useRef, useState } from 'react'

export function ChessClock({ milliseconds = 0, active }: { milliseconds?: number; active: boolean }) {
  const receivedAt = useRef(Date.now())
  const [now, setNow] = useState(() => Date.now())
  useEffect(() => {
    receivedAt.current = Date.now()
    setNow(receivedAt.current)
    if (!active) return
    const timer = window.setInterval(() => setNow(Date.now()), 250)
    return () => window.clearInterval(timer)
  }, [milliseconds, active])
  const live = Math.max(0, milliseconds - (active ? now - receivedAt.current : 0))
  const minutes = Math.floor(live / 60000)
  const seconds = Math.floor((live % 60000) / 1000)
  return <time className={live < 30000 ? 'is-low' : ''}>{minutes}:{seconds.toString().padStart(2, '0')}</time>
}
