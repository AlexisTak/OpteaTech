import { PublicRequestForm } from '../components/public-request-form';
import { getBffBaseUrl, getBffStatus } from '@/lib/bff';

export default async function DemandePage() {
  const status = await getBffStatus();
  const bffUrl = getBffBaseUrl();

  const statusClass =
    status === 'online'
      ? 'bg-emerald-50 text-emerald-700 border-emerald-200'
      : status === 'degraded'
        ? 'bg-amber-50 text-amber-700 border-amber-200'
        : 'bg-red-50 text-red-700 border-red-200';

  return (
    <div className="min-h-screen bg-slate-100 text-slate-900">
      <main className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-10 sm:px-6">
        <section className="rounded-xl border border-slate-300 bg-white p-6 shadow-sm">
          <h1 className="text-3xl font-bold tracking-tight">Nouvelle demande</h1>
          <p className="mt-2 text-slate-600">Formulaire connecté a Go API via le proxy Next.js.</p>
          <div className="mt-4 flex flex-wrap items-center gap-3 text-sm">
            <span className={`rounded-full border px-3 py-1 font-medium ${statusClass}`}>
              Go API status: {status}
            </span>
            <span className="text-slate-600">{bffUrl}</span>
          </div>
        </section>

        <PublicRequestForm />
      </main>
    </div>
  );
}
