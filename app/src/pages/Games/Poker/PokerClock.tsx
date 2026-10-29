import { useEffect, useRef, useState } from "react";

export function PokerClock({
  deadline,
  active,
}: {
  deadline?: string;
  active: boolean;
}) {
  const [now, setNow] = useState(Date.now());
  const last = useRef(deadline);
  useEffect(() => {
    last.current = deadline;
    setNow(Date.now());
    if (!active || !deadline) return;
    const timer = window.setInterval(() => setNow(Date.now()), 250);
    return () => clearInterval(timer);
  }, [deadline, active]);
  if (!deadline) return <span>—</span>;
  const seconds = Math.max(
    0,
    Math.ceil((new Date(last.current ?? deadline).getTime() - now) / 1000),
  );
  return <time className={seconds <= 8 ? "is-low" : ""}>{seconds}s</time>;
}
