import { NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/lib/store/auth.store';

const links = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/requests', label: 'Demandes' },
  { to: '/members', label: 'Clients' },
  { to: '/messages', label: 'Messages' },
  { to: '/contacts', label: 'Contacts' },
  { to: '/settings', label: 'Parametres' },
];

export function Sidebar() {
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);

  return (
    <aside className="sidebar">
      <div className="stack">
        {links.map((item) => (
          <NavLink key={item.to} to={item.to} className={({ isActive }) => `sidebar-link${isActive ? ' active' : ''}`}>
            {item.label}
          </NavLink>
        ))}
      </div>
      <div style={{ marginTop: 16 }}>
        <button
          className="button secondary"
          onClick={async () => {
            await logout();
            navigate('/login');
          }}
        >
          Deconnexion
        </button>
      </div>
    </aside>
  );
}
