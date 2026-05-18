import { create } from 'zustand';
import { invoke } from '@tauri-apps/api/core';

type AuthState = {
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  setTokens: (access: string, refresh: string) => Promise<void>;
  initFromKeychain: () => Promise<void>;
  logout: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  setTokens: async (access, refresh) => {
    await invoke('store_tokens', {
      tokens: { access_token: access, refresh_token: refresh },
    });
    set({ accessToken: access, refreshToken: refresh, isAuthenticated: true });
  },
  initFromKeychain: async () => {
    const tokens = await invoke<{ access_token: string; refresh_token: string } | null>('get_tokens').catch(() => null);
    if (!tokens) return;
    set({ accessToken: tokens.access_token, refreshToken: tokens.refresh_token, isAuthenticated: true });
  },
  logout: async () => {
    await invoke('clear_tokens').catch(() => undefined);
    set({ accessToken: null, refreshToken: null, isAuthenticated: false });
  },
}));
