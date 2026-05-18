import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/lib/store/auth.store';

const GO_URL = import.meta.env.VITE_GO_API_URL ?? 'http://localhost:3001';

export const goClient = axios.create({
  baseURL: `${GO_URL}/api`,
  timeout: 12000,
  headers: {
    'Content-Type': 'application/json',
  },
});

goClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = useAuthStore.getState().accessToken;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let refreshing = false;
let queue: Array<{ resolve: (token: string) => void; reject: (error: unknown) => void }> = [];

goClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status !== 401 || original._retry) {
      return Promise.reject(error);
    }

    if (refreshing) {
      return new Promise((resolve, reject) => {
        queue.push({ resolve: (token) => {
          original.headers.Authorization = `Bearer ${token}`;
          resolve(goClient(original));
        }, reject });
      });
    }

    original._retry = true;
    refreshing = true;

    try {
      const { refreshToken, setTokens } = useAuthStore.getState();
      if (!refreshToken) throw new Error('Missing refresh token');

      const refreshResponse = await axios.post(`${GO_URL}/api/auth/refresh`, {
        refresh_token: refreshToken,
      });

      const access = refreshResponse.data?.access_token;
      const nextRefresh = refreshResponse.data?.refresh_token ?? refreshToken;
      if (!access) throw new Error('Refresh failed');

      await setTokens(access, nextRefresh);
      queue.forEach((job) => job.resolve(access));
      queue = [];

      original.headers.Authorization = `Bearer ${access}`;
      return goClient(original);
    } catch (refreshError) {
      queue.forEach((job) => job.reject(refreshError));
      queue = [];
      await useAuthStore.getState().logout();
      return Promise.reject(refreshError);
    } finally {
      refreshing = false;
    }
  },
);
