import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export function useKeyboardShortcuts() {
  const navigate = useNavigate();

  useEffect(() => {
    const handler = (event: KeyboardEvent) => {
      if (!event.metaKey && !event.ctrlKey) return;

      switch (event.key) {
        case '1':
          event.preventDefault();
          navigate('/dashboard');
          break;
        case '2':
          event.preventDefault();
          navigate('/requests');
          break;
        case '3':
          event.preventDefault();
          navigate('/members');
          break;
        case '4':
          event.preventDefault();
          navigate('/messages');
          break;
        case '8':
          event.preventDefault();
          navigate('/contacts');
          break;
        default:
          break;
      }
    };

    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [navigate]);
}
