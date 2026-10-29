import { alignLoopbackHost } from "./loopback";

export function getApiRoot() {
  const configured = (import.meta.env.VITE_API_ROOT ?? import.meta.env.VITE_SERVER_ROOT) as string | undefined;
  const root = configured?.trim().replace(/\/+$/, "") || "/api";
  return alignLoopbackHost(root, window.location.hostname);
}
