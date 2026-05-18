import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export function CommandPalette() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const handler = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
        event.preventDefault();
        setOpen((v) => !v);
      }
    };

    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, []);

  if (!open) return null;

  const go = (path: string) => {
    navigate(path);
    setOpen(false);
  };

  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.2)', display: 'grid', placeItems: 'start center', paddingTop: 120 }}>
      <div className="card" style={{ width: 460 }}>
        <div className="stack">
          <button className="button secondary" onClick={() => go('/dashboard')}>Dashboard</button>
          <button className="button secondary" onClick={() => go('/requests')}>Demandes</button>
          <button className="button secondary" onClick={() => go('/members')}>Clients</button>
          <button className="button secondary" onClick={() => go('/messages')}>Messages</button>
          <button className="button secondary" onClick={() => go('/contacts')}>Contacts</button>
          <button className="button secondary" onClick={() => setOpen(false)}>Fermer</button>
        </div>
      </div>
    </div>
  );
}
