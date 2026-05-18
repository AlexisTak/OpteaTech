import { createBrowserRouter, Navigate } from 'react-router-dom';
import App from '@/App';
import { LoginPage } from '@/pages/auth/Login';
import { DashboardPage } from '@/pages/dashboard/Dashboard';
import { RequestsPage } from '@/pages/requests/RequestsList';
import { MembersPage } from '@/pages/members/MembersList';
import { MessagesPage } from '@/pages/messages/Messages';
import { ContactsPage } from '@/pages/contacts/ContactsList';
import { SettingsPage } from '@/pages/settings/Settings';

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'requests', element: <RequestsPage /> },
      { path: 'members', element: <MembersPage /> },
      { path: 'messages', element: <MessagesPage /> },
      { path: 'contacts', element: <ContactsPage /> },
      { path: 'settings', element: <SettingsPage /> },
    ],
  },
]);
