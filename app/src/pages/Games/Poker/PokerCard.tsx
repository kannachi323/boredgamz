import type { PokerCard as Card } from "../../../hooks/usePokerMatch";

const ranks: Record<number, string> = { 11: "J", 12: "Q", 13: "K", 14: "A" };
const suits = { s: "♠", h: "♥", d: "♦", c: "♣" };

export function PokerCard({
  card,
  hidden,
  selected,
  compact = false,
}: {
  card?: Card;
  hidden?: boolean;
  selected?: boolean;
  compact?: boolean;
}) {
  if (hidden || !card)
    return (
      <span
        className={`pkr-card is-back ${compact ? "is-compact" : ""}`}
        aria-label="Hidden card"
      >
        <i>BG</i>
        <b>♠</b>
      </span>
    );
  const red = card.suit === "h" || card.suit === "d";
  const label = ranks[card.rank] ?? String(card.rank);
  return (
    <span
      className={`pkr-card ${red ? "is-red" : ""} ${selected ? "is-selected" : ""} ${compact ? "is-compact" : ""}`}
      aria-label={`${label} of ${suits[card.suit]}`}
    >
      <span>
        <b>{label}</b>
        <i>{suits[card.suit]}</i>
      </span>
      <em>{suits[card.suit]}</em>
      <span>
        <b>{label}</b>
        <i>{suits[card.suit]}</i>
      </span>
    </span>
  );
}
