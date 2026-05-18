import { MagneticButton } from '@/components/shared/MagneticButton';

const stats = [
  { value: '+40', label: 'projets livrés' },
  { value: '100%', label: 'satisfaction client' },
  { value: '48h', label: 'délai de réponse' },
  { value: '5 ans', label: "d'expérience" },
];

export function FinalCta() {
  return (
    <section className="section-shell pt-0" data-reveal>
      <div className="cta-panel">
        {/* Left */}
        <div className="cta-panel__left">
          <p className="label">PRÊT À COMMENCER ?</p>
          <h2 className="display-lg mt-5">
            Parlons de
            <br />
            votre projet.
          </h2>
          <div className="mt-10">
            <MagneticButton
              href="/contact"
              className="bg-white text-[var(--ink)] hover:bg-[var(--accent)] hover:text-white"
            >
              Prendre contact
            </MagneticButton>
          </div>
          <p className="label mt-8">Premier échange gratuit · Réponse sous 48h · Devis détaillé</p>
        </div>

        {/* Right — stats */}
        <div className="cta-panel__right">
          <dl className="cta-stats">
            {stats.map((s) => (
              <div key={s.label} className="cta-stat">
                <dt className="cta-stat__value">{s.value}</dt>
                <dd className="cta-stat__label">{s.label}</dd>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </section>
  );
}
