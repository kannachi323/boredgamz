import { useEffect, useRef, useState } from "react";

/** Counts down to a server deadline (per-move clocks). */
export function DeadlineClock({ deadline, urgentAt = 5 }: { deadline?: string; urgentAt?: number }) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    if (!deadline) return;
    setNow(Date.now());
    const timer = window.setInterval(() => setNow(Date.now()), 250);
    return () => window.clearInterval(timer);
  }, [deadline]);
  if (!deadline) return null;
  const seconds = Math.max(0, Math.ceil((new Date(deadline).getTime() - now) / 1000));
  const label = seconds >= 60 ? `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, "0")}` : `${seconds}s`;
  return <time className={`clock num ${seconds <= urgentAt ? "is-low" : ""}`}>{label}</time>;
}

/** A chess-style bank clock that ticks locally between server updates. */
export function BankClock({ milliseconds, running }: { milliseconds?: number; running: boolean }) {
  const receivedAt = useRef(Date.now());
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    receivedAt.current = Date.now();
    setNow(receivedAt.current);
    if (!running) return;
    const timer = window.setInterval(() => setNow(Date.now()), 250);
    return () => window.clearInterval(timer);
  }, [milliseconds, running]);
  if (milliseconds === undefined) return null;
  const live = Math.max(0, milliseconds - (running ? now - receivedAt.current : 0));
  const minutes = Math.floor(live / 60000);
  const seconds = Math.floor((live % 60000) / 1000);
  return <time className={`clock num ${running ? "is-running" : ""} ${live < 20000 ? "is-low" : ""}`}>{minutes}:{String(seconds).padStart(2, "0")}</time>;
}
