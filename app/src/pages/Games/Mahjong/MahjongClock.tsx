import { useEffect, useState } from "react";

export function MahjongClock({ deadline }: { deadline?: string }) {
  const [now, setNow] = useState(Date.now());
  useEffect(() => {
    if (!deadline) return;
    const timer = window.setInterval(() => setNow(Date.now()), 250);
    return () => clearInterval(timer);
  }, [deadline]);
  if (!deadline) return <span className="mjh-clock is-untimed">∞ <small>UNTIMED</small></span>;
  const remaining = Math.max(0, new Date(deadline).getTime() - now);
  const seconds = Math.ceil(remaining / 1000);
  return (
    <span className={`mjh-clock ${seconds <= 60 ? "is-low" : ""}`}>
      {Math.floor(seconds / 60)}:{String(seconds % 60).padStart(2, "0")}
      <small>REMAINING</small>
    </span>
  );
}
