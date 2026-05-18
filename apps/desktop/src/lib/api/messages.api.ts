import { goClient } from '@/lib/api/client';

export const messagesApi = {
  listByRequest: (requestId: string) => goClient.get(`/admin/requests/${requestId}/messages`).then((r) => r.data),
  send: (requestId: string, content: string, attachments?: string[]) =>
    goClient.post(`/admin/requests/${requestId}/messages`, { content, attachments }).then((r) => r.data),
};
