import { Flag, RotateCcw } from "lucide-react";
import { useEffect, useState } from "react";

/** Two-step resign so a stray tap can't end a game. */
export function ResignButton({ onResign }: { onResign: () => void }) {
  const [armed, setArmed] = useState(false);
  useEffect(() => {
    if (!armed) return;
    const timer = window.setTimeout(() => setArmed(false), 3000);
    return () => window.clearTimeout(timer);
  }, [armed]);
  return (
    <button className={`btn ${armed ? "btn--ink" : ""}`} type="button" onClick={() => (armed ? onResign() : setArmed(true))}>
      <Flag size={15} /> {armed ? "Confirm resign" : "Resign"}
    </button>
  );
}

export function RematchButton({ sent, label, onClick }: { sent: boolean; label: string; onClick: () => void }) {
  return (
    <button className="btn btn--primary" type="button" disabled={sent} onClick={onClick}>
      <RotateCcw size={15} /> {sent ? "Waiting for them…" : label}
    </button>
  );
}
