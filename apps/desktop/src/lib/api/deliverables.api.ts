import { goClient } from '@/lib/api/client';

export const deliverablesApi = {
  add: (requestId: string, data: unknown) => goClient.post(`/admin/requests/${requestId}/deliverables`, data).then((r) => r.data),
};
