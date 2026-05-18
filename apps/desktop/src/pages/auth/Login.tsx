import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authApi } from '@/lib/api/auth.api';
import { useAuthStore } from '@/lib/store/auth.store';

export function LoginPage() {
  const navigate = useNavigate();
  const setTokens = useAuthStore((s) => s.setTokens);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);

    try {
      const payload = await authApi.login(email, password);
      const access = payload?.access_token ?? payload?.data?.access_token;
      const refresh = payload?.refresh_token ?? payload?.data?.refresh_token;

      if (!access || !refresh) {
        setError('Reponse de login invalide');
        return;
      }

      await setTokens(access, refresh);
      navigate('/dashboard');
    } catch {
      setError('Identifiants invalides');
    }
  };

  return (
    <div className="desktop-shell" style={{ justifyContent: 'center', alignItems: 'center' }}>
      <form onSubmit={onSubmit} className="card stack" style={{ width: 420 }}>
        <h1 style={{ margin: 0 }}>Connexion admin</h1>
        <input className="input" placeholder="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input className="input" type="password" placeholder="mot de passe" value={password} onChange={(e) => setPassword(e.target.value)} />
        {error ? <p style={{ color: '#b91c1c', margin: 0 }}>{error}</p> : null}
        <button className="button" type="submit">Se connecter</button>
      </form>
    </div>
  );
}
