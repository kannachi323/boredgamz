import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  username: string
}

interface AuthState {
  isAuthenticated: boolean
  user: User | null
  isGuest: boolean

  setIsAuthenticated: (val: boolean) => void
  setUser: (user: User | null) => void
  setIsGuest: (val: boolean) => void
  authLoading: boolean

  setAuthLoading(val: boolean): void
  checkAuth: (callback: () => void | Promise<void>) => Promise<boolean>
  login: (email: string, password: string) => Promise<boolean>
  logout: (callback: () => void | Promise<void>) => Promise<boolean>
  signup: (email: string, password: string) => Promise<boolean>
  getOrCreateGuestSession: () => Promise<boolean>
}

export const useAuthStore = create(
  persist<AuthState>(
    (set, get) => ({
      isAuthenticated: false,
      user: null,
      isGuest: false,
      authLoading: false,
      setAuthLoading: (val) => set({ authLoading: val }),
      setIsAuthenticated: (val) => set({ isAuthenticated: val }),
      setUser: (user) => set({ user }),
      setIsGuest: (val) => set({ isGuest: val }),
      signup: async (email, password) => {
        const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/signup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email: email, password: password }),
        });
        return res.ok
      },

      login: async (email, password) => {
        const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/login`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ email: email, password: password }),
        });

        return res.ok;
      },
      logout: async (callback: () => void | Promise<void>) => {
        const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/logout`, {
          method: "GET",
          credentials: "include",
        });
        set({ isAuthenticated: false, user: null, isGuest: false });
        callback();
        return res.ok;
      },
      getOrCreateGuestSession: async () => {
        try {
          const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/guest-session`, {
            method: 'GET',
            credentials: 'include',
          })
          if (res.ok) {
            const data = await res.json()
            set({
              user: { id: data.id, username: data.username },
              isAuthenticated: true,
              isGuest: true,
            })
            return true
          }
          return false
        } catch (err) {
          console.warn('Failed to create guest session:', err)
          return false
        }
      },
      checkAuth: async (onFail: () => void | Promise<void>) => {
        const { setAuthLoading, getOrCreateGuestSession } = get();
        try {
          setAuthLoading(true);
          const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/check-auth`, {
            method: 'GET',
            credentials: 'include',
          })
          if (res.ok) {
            const data = await res.json()
            set({
              user: { id: data.id, username: data.username },
              isAuthenticated: true,
              isGuest: false,
            })
            return true
          } else {
            // Not authenticated, try to create guest session
            set({ isAuthenticated: false })
            const guestCreated = await getOrCreateGuestSession()
            if (!guestCreated) {
              onFail()
            }
            return guestCreated
          }
        } catch (err) {
          console.warn('Auth check failed:', err)
          set({ isAuthenticated: false })
          // Try guest session as fallback
          const guestCreated = await getOrCreateGuestSession()
          if (!guestCreated) {
            onFail()
          }
          return guestCreated
        } finally {
          setAuthLoading(false)
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user,
        isGuest: state.isGuest,
      }) as unknown as AuthState,
    }
  )
)
