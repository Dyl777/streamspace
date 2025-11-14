import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface UserState {
  username: string | null;
  role: 'user' | 'admin' | null;
  setUser: (username: string, role: 'user' | 'admin') => void;
  clearUser: () => void;
  isAuthenticated: boolean;
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      username: null,
      role: null,
      isAuthenticated: false,
      setUser: (username, role) =>
        set({ username, role, isAuthenticated: true }),
      clearUser: () =>
        set({ username: null, role: null, isAuthenticated: false }),
    }),
    {
      name: 'streamspace-user',
    }
  )
);
