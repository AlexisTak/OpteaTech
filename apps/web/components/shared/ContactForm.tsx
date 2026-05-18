'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

const contactSchema = z.object({
  name: z.string().min(1, 'Le nom est requis'),
  email: z.string().email('Adresse email invalide'),
  company: z.string().optional(),
  serviceInterest: z.string().min(1, 'Sélectionnez un service'),
  budgetRange: z.string().optional(),
  message: z.string().min(20, 'Le message doit faire au moins 20 caractères'),
  website: z.string().optional(),
});

type ContactValues = z.infer<typeof contactSchema>;

const defaultValues: ContactValues = {
  name: '',
  email: '',
  company: '',
  serviceInterest: '',
  budgetRange: '',
  message: '',
  website: '',
};

export function ContactForm() {
  const [feedback, setFeedback] = useState<{ type: 'success' | 'error'; message: string } | null>(null);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ContactValues>({
    resolver: zodResolver(contactSchema),
    defaultValues,
  });

  async function onSubmit(values: ContactValues) {
    setFeedback(null);
    const response = await fetch('/api/contact', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(values),
    });

    if (!response.ok) {
      setFeedback({ type: 'error', message: 'Impossible d envoyer votre message pour le moment.' });
      return;
    }

    reset(defaultValues);
    setFeedback({ type: 'success', message: 'Message envoye. Retour sous 48h.' });
  }

  return (
    <form className="card grid gap-4 p-6" onSubmit={handleSubmit(onSubmit)}>
      <div className="grid gap-4 sm:grid-cols-2">
        <label className="field">
          Nom
          <input {...register('name')} aria-invalid={Boolean(errors.name)} />
          {errors.name ? <span className="field-error">{errors.name.message}</span> : null}
        </label>

        <label className="field">
          Email
          <input type="email" {...register('email')} aria-invalid={Boolean(errors.email)} />
          {errors.email ? <span className="field-error">{errors.email.message}</span> : null}
        </label>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <label className="field">
          Société
          <input {...register('company')} />
        </label>

        <label className="field">
          Service
          <select {...register('serviceInterest')} aria-invalid={Boolean(errors.serviceInterest)}>
            <option value="">Choisir</option>
            <option value="Sites Web">Sites Web</option>
            <option value="Logiciel">Logiciel</option>
            <option value="IA">IA</option>
            <option value="Conseil">Conseil</option>
            <option value="Autre">Autre</option>
          </select>
          {errors.serviceInterest ? <span className="field-error">{errors.serviceInterest.message}</span> : null}
        </label>
      </div>

      <label className="field">
        Budget
        <select {...register('budgetRange')}>
          <option value="">Choisir</option>
          <option value="lt-2k">&lt; 2k EUR</option>
          <option value="2k-5k">2-5k EUR</option>
          <option value="5k-15k">5-15k EUR</option>
          <option value="15k-plus">15k+ EUR</option>
        </select>
      </label>

      <label className="field sr-only" aria-hidden="true">
        Website
        <input tabIndex={-1} autoComplete="off" {...register('website')} />
      </label>

      <label className="field">
        Message
        <textarea rows={6} {...register('message')} aria-invalid={Boolean(errors.message)} />
        {errors.message ? <span className="field-error">{errors.message.message}</span> : null}
      </label>

      <button type="submit" className="btn-primary w-full justify-center sm:w-fit" disabled={isSubmitting}>
        {isSubmitting ? 'Envoi en cours...' : 'Envoyer la demande'}
      </button>

      {feedback ? (
        <p className={feedback.type === 'success' ? 'text-[var(--success)]' : 'text-[var(--danger)]'}>{feedback.message}</p>
      ) : null}
    </form>
  );
}
