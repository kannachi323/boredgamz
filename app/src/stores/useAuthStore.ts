import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: string
  username: string
}

interface AuthState {
  isAuthenticated: boolean
  user: User | null

  setIsAuthenticated: (val: boolean) => void
  setUser: (user: User | null) => void
  authLoading: boolean

  setAuthLoading(val: boolean): void
  checkAuth: (callback: () => void | Promise<void>) => Promise<boolean>
  login: (email: string, password: string) => Promise<boolean>
  logout: (callback: () => void | Promise<void>) => Promise<boolean>
  signup: (email: string, password: string) => Promise<boolean>
}

export const useAuthStore = create(
  persist<AuthState>(
    (set, get) => ({
      isAuthenticated: false,
      user: null,
      authLoading: false,
      setAuthLoading: (val) => set({ authLoading: val }),
      setIsAuthenticated: (val) => set({ isAuthenticated: val }),
      setUser: (user) => set({ user }),
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
        set({ isAuthenticated: false, user: null });
        callback();
        return res.ok;
      },
      checkAuth: async (onFail: () => void | Promise<void>) => {
        const { setAuthLoading } = get();
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
            })
            return true
          } else {
            set({ isAuthenticated: false })
            onFail()
            return false
          }
        } catch (err) {
          console.warn('Auth check failed:', err)
          set({ isAuthenticated: false })
          onFail()
          return false
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
      }) as unknown as AuthState,
    }
  )
)
