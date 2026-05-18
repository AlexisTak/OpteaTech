'use client';

import { useRef } from 'react';
import { MagneticButton } from '@/components/shared/MagneticButton';
import { useHeroAnimation } from '@/hooks/useHeroAnimation';
import { useHeroParallax } from '@/hooks/useHeroParallax';

interface AnimatedCounterProps {
  count: number;
  suffix?: string;
  label: string;
}

function AnimatedCounter({ count, suffix = '', label }: AnimatedCounterProps) {
  return (
    <div className="hero-stat">
      <p className="font-display text-3xl text-[var(--text-primary)]">
        <span data-count={count} data-suffix={suffix}>
          {count}
          {suffix}
        </span>
      </p>
      <p className="label mt-2">{label}</p>
    </div>
  );
}

function WordReveal({ text, delayStep = 100 }: { text: string; delayStep?: number }) {
  const words = text.split(' ');
  return (
    <>
      {words.map((word, index) => (
        <span key={`${word}-${index}`} style={{ display: 'inline-block', overflow: 'hidden', verticalAlign: 'top' }}>
          <span className="hero-word" style={{ display: 'inline-block' }} data-delay={index * delayStep}>
            {word}
            &nbsp;
          </span>
        </span>
      ))}
    </>
  );
}

export function Hero() {
  const sectionRef = useRef<HTMLElement>(null);
  useHeroAnimation({ sectionRef });
  useHeroParallax({ sectionRef });

  return (
    <section ref={sectionRef} id="accueil" className="hero-section section-shell relative min-h-[100dvh] overflow-hidden pt-28" data-reveal>
      <div className="hero-bg-grid hero-dots absolute inset-0" aria-hidden="true" />
      <div className="relative flex min-h-[calc(100dvh-112px)] flex-col justify-start pt-10 md:pt-14 gap-10">
        <p className="hero-label label inline-flex items-center gap-2">
          <span aria-hidden="true">●</span>
          AGENCE TECH · FRANCE · 2025
        </p>

        <h1 className="hero-title hero-title-wrap display-xl max-w-[14ch] text-[var(--ink)]">
          <WordReveal text="Votre" />
          <span style={{ fontStyle: 'italic' }}>
            <WordReveal text="vision," delayStep={120} />
          </span>
          <br />
          <WordReveal text="notre travail." delayStep={120} />
        </h1>

        <div className="hero-divider divider" />

        <div className="grid gap-8 lg:grid-cols-[1.2fr_0.8fr]">
          <p className="hero-desc body-lg max-w-2xl">
            Sites web, logiciels et solutions IA sur mesure. Un interlocuteur direct, des resultats durables.
          </p>
          <MagneticButton href="#services" className="self-start">
            Découvrir nos services
          </MagneticButton>
        </div>

        <div className="divider" />

        <div className="grid grid-cols-3 gap-5 sm:max-w-xl">
          <AnimatedCounter count={12} suffix="+" label="Projets livres" />
          <AnimatedCounter count={100} suffix="%" label="Satisfaction" />
          <AnimatedCounter count={48} suffix="h" label="Delai reponse" />
        </div>
      </div>

      {/* Scroll indicator — bottom center */}
      <div className="scroll-indicator absolute bottom-8 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 text-[var(--ink-30)]" aria-hidden="true">
        <span className="label tracking-widest">DÉFILER</span>
        <span className="relative h-10 w-[1px] bg-[var(--ink-10)]">
          <span
            className="absolute left-0 top-0 h-3 w-[1px] bg-[var(--ink)]"
            style={{ animation: 'scroll-indicator 1.6s ease-in-out infinite' }}
          />
        </span>
      </div>
    </section>
  );
}
