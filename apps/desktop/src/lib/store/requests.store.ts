import { create } from 'zustand';

type RequestItem = { id: string; title: string; status: string };

type RequestsState = {
  requests: RequestItem[];
  setRequests: (items: RequestItem[]) => void;
};

export const useRequestsStore = create<RequestsState>((set) => ({
  requests: [],
  setRequests: (items) => set({ requests: items }),
}));
