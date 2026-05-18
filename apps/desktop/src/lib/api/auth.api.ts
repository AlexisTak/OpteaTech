import axios from 'axios';

const GO_URL = import.meta.env.VITE_GO_API_URL ?? 'http://localhost:3001';

export const authApi = {
  login: async (email: string, password: string) => {
    const { data } = await axios.post(`${GO_URL}/api/auth/login`, { email, password });
    return data;
  },
};
