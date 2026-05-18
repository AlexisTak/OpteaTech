"use client";

import { FormEvent, useMemo, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { getBffBaseUrl } from '@/lib/bff';

type Envelope<T> = {
  data?: T;
  error?: string;
};

type DashboardData = Record<string, unknown>;

type DashboardRequest = {
  id?: string;
  client_name?: string;
  clientName?: string;
};

export default function LoginPage() {
  const searchParams = useSearchParams();
  const [tokenInput, setTokenInput] = useState('');
  const [email, setEmail] = useState('');
  const [dashboard, setDashboard] = useState<DashboardData | null>(null);
  const [dashboardError, setDashboardError] = useState<string | null>(null);
  const [linkError, setLinkError] = useState<string | null>(null);
  const [linkMessage, setLinkMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [linkLoading, setLinkLoading] = useState(false);

  const resolvedToken = useMemo(
    () => (searchParams.get('token') ?? tokenInput).trim(),
    [searchParams, tokenInput],
  );

  const linkedRequest = useMemo(() => {
    if (!dashboard || typeof dashboard !== 'object') {
      return null;
    }

    const request = (dashboard as { request?: DashboardRequest }).request;
    if (!request || typeof request !== 'object') {
      return null;
    }

    return {
      requestId: typeof request.id === 'string' ? request.id : '',
      clientName:
        typeof request.client_name === 'string'
          ? request.client_name
          : typeof request.clientName === 'string'
            ? request.clientName
            : '',
    };
  }, [dashboard]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);
    setDashboardError(null);
    setDashboard(null);

    try {
      const response = await fetch(`${getBffBaseUrl()}/client/dashboard?token=${encodeURIComponent(resolvedToken)}`, {
        cache: 'no-store',
      });

      const body = (await response.json().catch(() => ({}))) as Envelope<DashboardData> | DashboardData;
      const payload = ('data' in body ? body.data : body) as DashboardData | undefined;

      if (!response.ok || !payload) {
        const message = 'error' in body && typeof body.error === 'string' ? body.error : undefined;
        setDashboardError(message ?? 'Impossible d acceder au portail membre avec ce token.');
        return;
      }

      setDashboard(payload);
    } catch {
      setDashboardError('Le service Go est injoignable. Verifie que l API tourne sur le port 3001.');
    } finally {
      setLoading(false);
    }
  };

  const onRequestNewLink = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLinkLoading(true);
    setLinkError(null);
    setLinkMessage(null);

    try {
      const response = await fetch(`${getBffBaseUrl()}/client/request-new-link`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          request_id: linkedRequest?.requestId,
        }),
      });

      const body = (await response.json().catch(() => ({}))) as Envelope<Record<string, unknown>>;

      if (!response.ok) {
        setLinkError(body.error ?? 'Impossible de demander un nouveau lien pour le moment.');
        return;
      }

      setLinkMessage('Nouveau lien demande. Verifiez votre boite mail.');
    } catch {
      setLinkError('Le service Go est injoignable.');
    } finally {
      setLinkLoading(false);
    }
  };

  return (
    <main className="section-shell pt-32">
      <div className="mx-auto max-w-2xl rounded-[24px] border border-[var(--border)] bg-[var(--surface)] p-8 md:p-10">
        <p className="section-label mb-6">
          <span>Membres / Portail client</span>
        </p>
        <h1 className="display-md">Connexion membre</h1>
        <p className="body-md mt-4 text-[var(--ink-60)]">
          Entrez votre token pour acceder a votre espace client.
        </p>

        <form className="mt-8 grid gap-4" onSubmit={onSubmit}>
          <label className="grid gap-2">
            <span className="label">Token client</span>
            <input
              required
              type="text"
              value={tokenInput}
              onChange={(event) => setTokenInput(event.target.value)}
              className="h-12 rounded-xl border border-[var(--border)] bg-white px-4 text-[var(--ink)] outline-none transition focus:border-[var(--accent)]"
              placeholder="Votre token client (64 caractères)"
            />
          </label>

          {dashboardError ? <p className="mt-1 text-sm text-red-600">{dashboardError}</p> : null}

          <button
            type="submit"
            disabled={loading || !resolvedToken}
            className="mt-2 inline-flex h-12 items-center justify-center rounded-xl border border-[var(--accent)] bg-[var(--accent)] px-5 text-sm font-medium text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? 'Connexion en cours...' : 'Acceder a mon espace'}
          </button>
        </form>

        {dashboard ? (
          <div className="mt-8 space-y-4">
            <p className="body-md text-[var(--ink-60)]">
              Connecte en tant que{' '}
              <span className="font-medium text-[var(--ink)]">
                {linkedRequest?.clientName || 'membre'}
              </span>{' '}
              via le proxy Next.js /api/go/client/dashboard.
            </p>
            <pre className="overflow-auto rounded-xl border border-[var(--border)] bg-white p-4 text-xs leading-6 text-[var(--ink)]">
              {JSON.stringify(dashboard, null, 2)}
            </pre>
          </div>
        ) : null}

        <div className="mt-10 border-t border-[var(--border)] pt-8">
          <h2 className="text-xl font-medium text-[var(--ink)]">Token perdu ?</h2>
          <p className="body-md mt-2 text-[var(--ink-60)]">
            Demandez un nouveau token. Le nom et la demande sont recuperes automatiquement depuis votre token actif.
          </p>
          <form className="mt-5 grid gap-4" onSubmit={onRequestNewLink}>
            <label className="grid gap-2">
              <span className="label">Email</span>
              <input
                required
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                className="h-12 rounded-xl border border-[var(--border)] bg-white px-4 text-[var(--ink)] outline-none transition focus:border-[var(--accent)]"
                placeholder="votre@email.com"
              />
            </label>

            <p className="text-sm text-[var(--ink-60)]">
              Nom lie au token: <span className="font-medium text-[var(--ink)]">{linkedRequest?.clientName || 'Non detecte'}</span>
            </p>

            {linkError ? <p className="text-sm text-red-600">{linkError}</p> : null}
            {linkMessage ? <p className="text-sm text-green-700">{linkMessage}</p> : null}

            <button
              type="submit"
              disabled={linkLoading || !linkedRequest?.requestId || !email.trim()}
              className="inline-flex h-11 items-center justify-center rounded-xl border border-[var(--border)] px-4 text-sm text-[var(--ink)] transition hover:border-[var(--accent)] hover:text-[var(--accent)] disabled:cursor-not-allowed disabled:opacity-60"
            >
              {linkLoading ? 'Envoi en cours...' : 'Demander un nouveau lien'}
            </button>
          </form>
        </div>
      </div>
    </main>
  );
}
