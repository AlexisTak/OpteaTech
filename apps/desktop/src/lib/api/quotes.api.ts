import { goClient } from '@/lib/api/client';

export const quotesApi = {
  set: (requestId: string, data: unknown) => goClient.post(`/admin/requests/${requestId}/quote`, data).then((r) => r.data),
};
