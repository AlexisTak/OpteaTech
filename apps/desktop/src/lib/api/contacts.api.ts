import { goClient } from '@/lib/api/client';

export const contactsApi = {
  list: (params?: Record<string, unknown>) => goClient.get('/admin/contacts', { params }).then((r) => r.data),
  get: (id: string) => goClient.get(`/admin/contacts/${id}`).then((r) => r.data),
  create: (data: unknown) => goClient.post('/admin/contacts', data).then((r) => r.data),
  update: (id: string, data: unknown) => goClient.put(`/admin/contacts/${id}`, data).then((r) => r.data),
  delete: (id: string) => goClient.delete(`/admin/contacts/${id}`).then((r) => r.data),
  addNote: (id: string, note: string) => goClient.post(`/admin/contacts/${id}/notes`, { note }).then((r) => r.data),
};
