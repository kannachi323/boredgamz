const browserIdentityKey = "boredgamz-player-id";
const tabIdentityKey = "boredgamz-player-tab-id";

export function getPlayerConnectionID() {
  const browserID = getOrCreateID(window.localStorage, browserIdentityKey, false);
  const tabID = getOrCreateID(window.sessionStorage, tabIdentityKey, true);
  return combinePlayerIdentity(browserID, tabID);
}

export function combinePlayerIdentity(browserID: string, tabID: string) {
  return `${browserID.slice(0, 48)}:${tabID.slice(0, 12)}`;
}

function getOrCreateID(storage: Storage, key: string, compact: boolean) {
  const stored = storage.getItem(key)?.trim();
  if (stored) return stored;
  const generated = createUUID();
  const id = compact ? generated.replace(/-/g, "").slice(0, 12) : generated;
  storage.setItem(key, id);
  return id;
}

function createUUID() {
  return typeof window.crypto.randomUUID === "function"
    ? window.crypto.randomUUID()
    : fallbackUUID();
}

function fallbackUUID() {
  const bytes = new Uint8Array(16);
  window.crypto.getRandomValues(bytes);
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  bytes[8] = (bytes[8] & 0x3f) | 0x80;
  const hex = Array.from(bytes, (byte) =>
    byte.toString(16).padStart(2, "0"),
  ).join("");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}
