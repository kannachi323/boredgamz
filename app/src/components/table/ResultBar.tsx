import type { ReactNode } from "react";

export function ResultBar({ title, detail, tone = "neutral", children }: { title: string; detail?: string; tone?: "win" | "loss" | "neutral"; children?: ReactNode }) {
  return (
    <div className={`result is-${tone}`} role="status">
      <div>
        <b>{title}</b>
        {detail && <span>{detail}</span>}
      </div>
      <div className="result__actions">{children}</div>
    </div>
  );
}
