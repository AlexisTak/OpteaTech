import { ContactForm } from '@/components/shared/ContactForm';

export default function ContactPage() {
  return (
    <section className="section-shell pt-24">
      <div className="mb-8 max-w-3xl space-y-4">
        <p className="eyebrow">Contact</p>
        <h1 className="font-display text-5xl text-[var(--text-primary)]">Parlons de votre projet</h1>
        <p className="text-[var(--text-secondary)]">
          Décris ton besoin, ton contexte et tes objectifs. Retour clair sous 48h avec une proposition réaliste.
        </p>
      </div>
      <ContactForm />
    </section>
  );
}
