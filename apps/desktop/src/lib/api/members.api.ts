import { goClient } from '@/lib/api/client';

export const membersApi = {
  list: (params?: Record<string, unknown>) => goClient.get('/admin/members', { params }).then((r) => r.data),
  get: (email: string) => goClient.get(`/admin/members/${encodeURIComponent(email)}`).then((r) => r.data),
};
