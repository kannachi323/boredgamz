import { useEffect, useState } from 'react'

export function TurnClock({ deadline, control }: { deadline?: string; control: string }) {
  const [remaining, setRemaining] = useState(() => secondsLeft(deadline))
  useEffect(() => {
    setRemaining(secondsLeft(deadline))
    if (!deadline) return
    const timer = window.setInterval(() => setRemaining(secondsLeft(deadline)), 200)
    return () => window.clearInterval(timer)
  }, [deadline])
  if (!deadline) return <span>{control === 'untimed' ? 'Untimed' : `${control} / move`}</span>
  return <span className={remaining <= 5 ? 'is-urgent' : ''}>{remaining.toFixed(1)}s</span>
}

function secondsLeft(deadline?: string) {
  if (!deadline) return 0
  return Math.max(0, (new Date(deadline).getTime() - Date.now()) / 1000)
}
