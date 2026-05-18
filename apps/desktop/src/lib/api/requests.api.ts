import { goClient } from '@/lib/api/client';

export const adminApi = {
  getDashboard: () => goClient.get('/admin/dashboard').then((r) => r.data),
  listRequests: (params?: Record<string, unknown>) => goClient.get('/admin/requests', { params }).then((r) => r.data),
  getRequest: (id: string) => goClient.get(`/admin/requests/${id}`).then((r) => r.data),
  updateStatus: (id: string, status: string) => goClient.put(`/admin/requests/${id}/status`, { status }).then((r) => r.data),
};
