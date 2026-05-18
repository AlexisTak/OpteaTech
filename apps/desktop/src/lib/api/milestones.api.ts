import { goClient } from '@/lib/api/client';

export const milestonesApi = {
  create: (requestId: string, data: unknown) => goClient.post(`/admin/requests/${requestId}/milestones`, data).then((r) => r.data),
  update: (requestId: string, milestoneId: string, data: unknown) =>
    goClient.put(`/admin/requests/${requestId}/milestones/${milestoneId}`, data).then((r) => r.data),
};
