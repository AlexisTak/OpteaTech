import { ReactNode } from 'react';
import { Outlet } from 'react-router-dom';
import { Sidebar } from '@/components/layout/Sidebar';
import { Titlebar } from '@/components/layout/Titlebar';
import { CommandPalette } from '@/components/layout/CommandPalette';
import { useKeyboardShortcuts } from '@/lib/hooks/useKeyboardShortcuts';

export function AppShell({ children }: { children?: ReactNode }) {
  useKeyboardShortcuts();

  return (
    <div className="desktop-shell">
      <Titlebar />
      <div className="main-layout">
        <Sidebar />
        <main className="content">{children ?? <Outlet />}</main>
      </div>
      <CommandPalette />
    </div>
  );
}
