import { create } from 'zustand';

type Notification = { id: string; message: string };

type NotificationsState = {
  items: Notification[];
  push: (message: string) => void;
  remove: (id: string) => void;
};

export const useNotificationsStore = create<NotificationsState>((set) => ({
  items: [],
  push: (message) =>
    set((state) => ({
      items: [...state.items, { id: crypto.randomUUID(), message }],
    })),
  remove: (id) => set((state) => ({ items: state.items.filter((i) => i.id !== id) })),
}));
