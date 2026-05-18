'use client';

import { FormEvent, useState } from 'react';
import { submitPublicRequest, type PublicRequestPayload } from '@/lib/bff';

type FormState = PublicRequestPayload;

const initialState: FormState = {
  client_name: '',
  client_email: '',
  client_company: '',
  client_phone: '',
  service_type: 'site_web',
  title: '',
  description: '',
  budget_range: '2k_5k',
  deadline: '',
};

export function PublicRequestForm() {
  const [form, setForm] = useState<FormState>(initialState);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<{ ok: boolean; message: string } | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setFeedback(null);

    const payload: PublicRequestPayload = {
      ...form,
      client_company: form.client_company || undefined,
      client_phone: form.client_phone || undefined,
      budget_range: form.budget_range || undefined,
      deadline: form.deadline || undefined,
    };

    const result = await submitPublicRequest(payload);
    setFeedback(result);
    setIsSubmitting(false);

    if (result.ok) {
      setForm(initialState);
    }
  }

  return (
    <form onSubmit={onSubmit} className="grid gap-3 rounded-xl border border-slate-300 bg-white p-4 shadow-sm">
      <h2 className="text-xl font-semibold text-slate-900">Start a project</h2>
      <p className="text-sm text-slate-600">This form is routed through Next.js proxy at /api/go/requests to Go API.</p>

      <input
        required
        value={form.client_name}
        onChange={(event) => setForm((prev) => ({ ...prev, client_name: event.target.value }))}
        placeholder="Your name"
        className="rounded-md border border-slate-300 px-3 py-2"
      />

      <input
        required
        type="email"
        value={form.client_email}
        onChange={(event) => setForm((prev) => ({ ...prev, client_email: event.target.value }))}
        placeholder="Email"
        className="rounded-md border border-slate-300 px-3 py-2"
      />

      <input
        value={form.client_company}
        onChange={(event) => setForm((prev) => ({ ...prev, client_company: event.target.value }))}
        placeholder="Company (optional)"
        className="rounded-md border border-slate-300 px-3 py-2"
      />

      <select
        value={form.service_type}
        onChange={(event) => setForm((prev) => ({ ...prev, service_type: event.target.value as FormState['service_type'] }))}
        className="rounded-md border border-slate-300 px-3 py-2"
      >
        <option value="site_web">Website</option>
        <option value="logiciel">Software</option>
        <option value="ia">AI</option>
        <option value="conseil">Consulting</option>
        <option value="autre">Other</option>
      </select>

      <input
        required
        value={form.title}
        onChange={(event) => setForm((prev) => ({ ...prev, title: event.target.value }))}
        placeholder="Project title"
        className="rounded-md border border-slate-300 px-3 py-2"
      />

      <textarea
        required
        value={form.description}
        onChange={(event) => setForm((prev) => ({ ...prev, description: event.target.value }))}
        placeholder="Describe your project"
        className="min-h-28 rounded-md border border-slate-300 px-3 py-2"
      />

      <button
        type="submit"
        disabled={isSubmitting}
        className="rounded-md bg-slate-900 px-4 py-2 font-medium text-white disabled:opacity-60"
      >
        {isSubmitting ? 'Sending...' : 'Send request'}
      </button>

      {feedback ? (
        <p className={feedback.ok ? 'text-sm text-emerald-700' : 'text-sm text-red-700'}>{feedback.message}</p>
      ) : null}
    </form>
  );
}
