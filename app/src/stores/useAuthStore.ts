import { create } from "zustand";
import { getApiRoot } from "../lib/api";

export interface AuthUser {
  id: string;
  username: string;
}

interface AuthState {
  status: "checking" | "anonymous" | "authenticated";
  user: AuthUser | null;
  initialize: () => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<AuthUser>;
  login: (email: string, password: string) => Promise<AuthUser>;
  logout: () => Promise<void>;
}

const apiRoot = getApiRoot();

export const useAuthStore = create<AuthState>((set) => ({
  status: "checking",
  user: null,
  initialize: async () => {
    let response = await authFetch("/v1/auth/session");
    if (response.status === 401) {
      const refreshed = await authFetch("/v1/auth/refresh", { method: "POST" });
      if (refreshed.ok) response = await authFetch("/v1/auth/session");
    }
    if (!response.ok) {
      set({ status: "anonymous", user: null });
      return;
    }
    set({ status: "authenticated", user: await readUser(response) });
  },
  signup: async (username, email, password) => {
    const response = await authFetch("/v1/auth/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });
    if (!response.ok) throw new Error(await readError(response));
    const user = await readUser(response);
    set({ status: "authenticated", user });
    return user;
  },
  login: async (email, password) => {
    const response = await authFetch("/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!response.ok) throw new Error(await readError(response));
    const user = await readUser(response);
    set({ status: "authenticated", user });
    return user;
  },
  logout: async () => {
    await authFetch("/v1/auth/logout", { method: "POST" });
    set({ status: "anonymous", user: null });
  },
}));

function authFetch(path: string, init?: RequestInit) {
  return fetch(`${apiRoot}${path}`, { ...init, credentials: "include" });
}

async function readUser(response: Response) {
  return (await response.json()) as AuthUser;
}

async function readError(response: Response) {
  const message = (await response.text()).trim();
  if (response.status === 409) return "That email is already registered.";
  if (response.status === 401) return "The email or password is incorrect.";
  return message || "The account service could not complete that request.";
}
