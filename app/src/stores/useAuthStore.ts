import { create } from "zustand";
import { getApiRoot } from "../lib/api";

export interface AuthUser {
  id: string;
  username: string;
}

export type RefreshResult = "ok" | "signed-out" | "offline";

interface AuthState {
  status: "checking" | "anonymous" | "authenticated";
  user: AuthUser | null;
  initialize: () => Promise<void>;
  /** "ok", "signed-out" (401) or "offline" (server unreachable). */
  refresh: () => Promise<RefreshResult>;
  signup: (username: string, email: string, password: string) => Promise<AuthUser>;
  login: (email: string, password: string) => Promise<AuthUser>;
  logout: () => Promise<void>;
}

const apiRoot = getApiRoot();
// Access cookies last an hour; refresh well before that while a tab is open.
const refreshEveryMs = 20 * 60 * 1000;
let refreshTimer: number | undefined;

function keepFresh(refresh: () => Promise<RefreshResult>) {
  window.clearInterval(refreshTimer);
  refreshTimer = window.setInterval(() => void refresh(), refreshEveryMs);
}

export const useAuthStore = create<AuthState>((set, get) => ({
  status: "checking",
  user: null,
  initialize: async () => {
    // /session answers 200 (signed in, access cookie rotated if needed) or
    // 204 (signed out), so a normal page load logs no failed requests.
    const result = await get().refresh();
    if (result !== "ok") set({ status: "anonymous", user: null });
  },
  refresh: async () => {
    try {
      const response = await authFetch("/v1/auth/session");
      if (response.status === 200) {
        set({ status: "authenticated", user: await readUser(response) });
        keepFresh(get().refresh);
        return "ok";
      }
      if (response.status === 204) {
        window.clearInterval(refreshTimer);
        set({ status: "anonymous", user: null });
        return "signed-out";
      }
      return "offline";
    } catch {
      return "offline";
    }
  },
  signup: async (username, email, password) => {
    const response = await authFetch("/v1/auth/signup", jsonPost({ username, email, password }));
    if (!response.ok) throw new Error(await readError(response));
    const user = await readUser(response);
    set({ status: "authenticated", user });
    keepFresh(get().refresh);
    return user;
  },
  login: async (email, password) => {
    const response = await authFetch("/v1/auth/login", jsonPost({ email, password }));
    if (!response.ok) throw new Error(await readError(response));
    const user = await readUser(response);
    set({ status: "authenticated", user });
    keepFresh(get().refresh);
    return user;
  },
  logout: async () => {
    window.clearInterval(refreshTimer);
    try {
      await authFetch("/v1/auth/logout", { method: "POST" });
    } finally {
      set({ status: "anonymous", user: null });
    }
  },
}));

function jsonPost(body: unknown): RequestInit {
  return { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) };
}

function authFetch(path: string, init?: RequestInit) {
  return fetch(`${apiRoot}${path}`, { ...init, credentials: "include" });
}

async function readUser(response: Response) {
  return (await response.json()) as AuthUser;
}

async function readError(response: Response) {
  if (response.status === 409) return "That email already has an account.";
  if (response.status === 401) return "Wrong email or password.";
  if (response.status === 429) return "Too many attempts. Wait a minute and try again.";
  if (response.status === 400) return "Use a 2–24 character name, a valid email and a 12–72 character password.";
  return "Something went wrong. Try again.";
}
