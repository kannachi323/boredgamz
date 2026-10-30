import { create } from "zustand";

/** Live tables take the whole viewport; the site chrome steps aside. */
export const useShell = create<{ immersive: boolean; setImmersive: (value: boolean) => void }>((set) => ({
  immersive: false,
  setImmersive: (immersive) => set({ immersive }),
}));
