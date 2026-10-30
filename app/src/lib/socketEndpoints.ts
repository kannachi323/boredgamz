export type SocketService = "arena" | "poker" | "mahjong";

interface SocketEndpointConfig {
  arenaRoot?: string;
  legacyApiRoot?: string;
  protocol: string;
  host: string;
}

export function resolveSocketRoot(
  service: SocketService,
  config: SocketEndpointConfig,
) {
  const arenaRoot = trimTrailingSlash(config.arenaRoot);
  if (arenaRoot) {
    if (service === "arena") return arenaRoot;
    return arenaRoot.replace(/\/arena$/, `/${service}`);
  }

  const legacyApiRoot = trimTrailingSlash(config.legacyApiRoot);
  if (legacyApiRoot) return `${legacyApiRoot}/v1/${service}`;

  const socketProtocol = config.protocol === "https:" ? "wss:" : "ws:";
  return `${socketProtocol}//${config.host}/api/v1/${service}`;
}

export function getSocketRoot(service: SocketService) {
  return alignLoopbackHost(resolveSocketRoot(service, {
    arenaRoot: import.meta.env.VITE_WS_ROOT as string | undefined,
    legacyApiRoot: import.meta.env.VITE_WEBSOCKET_ROOT as string | undefined,
    protocol: window.location.protocol,
    host: window.location.host,
  }), window.location.hostname);
}

function trimTrailingSlash(value: string | undefined) {
  return value?.trim().replace(/\/+$/, "") || undefined;
}
import { alignLoopbackHost } from "./loopback";
