"use client";

import { useRef } from 'react';
import { useProcessAnimation } from '@/hooks/useProcessAnimation';

const steps = [
  {
    title: 'Echange initial',
    text: 'On discute du projet, des besoins et des contraintes. Gratuit et sans engagement.',
  },
  {
    title: 'Proposition',
    text: 'Devis detaille, planning et stack proposes sous 48h.',
  },
  {
    title: 'Developpement',
    text: 'Iterations regulieres, demos et feedback continu.',
  },
  {
    title: 'Livraison',
    text: 'Deploiement, formation et suivi post-livraison.',
  },
];

export function Process() {
  const sectionRef = useRef<HTMLElement>(null);
  useProcessAnimation({ sectionRef });

  return (
    <section ref={sectionRef} className="process-section section-shell section-alt" data-reveal>
      <p className="section-label">
        <span>03 / Processus</span>
      </p>
      <h2 className="display-md">Simple. Transparent. Efficace.</h2>

      <div className="mt-12">
          <div className="hidden sm:block">
            <svg viewBox="0 0 800 2" fill="none" style={{ width: '100%', height: 2 }}>
              <path className="process-path" d="M 0 1 L 800 1" stroke="var(--accent)" strokeWidth="1" />
            </svg>
          </div>
          <div className="mt-8 grid grid-cols-2 gap-x-8 gap-y-10 sm:grid-cols-4">
            {steps.map((step, index) => (
              <article key={step.title} className={`process-step process-step-${index + 1}`}>
                <div className="mb-3 flex items-center gap-2">
                  <span className={`process-dot process-dot-${index + 1} inline-block h-2 w-2 rounded-full bg-[var(--accent)]`} />
                  <p className="label">{String(index + 1).padStart(2, '0')}</p>
                </div>
                <h3 className="mt-3 text-xl text-[var(--ink)]">{step.title}</h3>
                <p className="body-md mt-3">{step.text}</p>
              </article>
            ))}
          </div>
        </div>
    </section>
  );
}
