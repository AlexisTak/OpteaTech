'use client';

import { RefObject, useEffect } from 'react';
import { gsap } from '@/lib/gsap';

interface HeroAnimationRefs {
  sectionRef: RefObject<HTMLElement | null>;
}

export function useHeroAnimation({ sectionRef }: HeroAnimationRefs) {
  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const ctx = gsap.context(() => {
      const tl = gsap.timeline({ delay: 0.15 });

      tl.fromTo('.hero-label', { opacity: 0, y: 12 }, { opacity: 1, y: 0, duration: 0.6, ease: 'opteaOut' })
        .fromTo(
          '.hero-word',
          { yPercent: 110, opacity: 0 },
          { yPercent: 0, opacity: 1, duration: 0.85, ease: 'opteaOut', stagger: 0.08 },
          '-=0.25',
        )
        .fromTo('.hero-divider', { scaleX: 0, transformOrigin: 'left' }, { scaleX: 1, duration: 0.7, ease: 'opteaInOut' }, '-=0.2')
        .fromTo(['.hero-desc', '.hero-cta'], { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.6, ease: 'opteaOut', stagger: 0.12 }, '-=0.3')
        .fromTo('.hero-stat', { opacity: 0, y: 16 }, { opacity: 1, y: 0, duration: 0.5, ease: 'opteaOut', stagger: 0.08 }, '-=0.2');

      gsap.fromTo('.scroll-indicator', { opacity: 0.3, y: 0 }, { opacity: 1, y: 8, duration: 1.2, ease: 'power1.inOut', yoyo: true, repeat: -1 });

      document.querySelectorAll<HTMLElement>('[data-count]').forEach((el, index) => {
        const target = Number(el.dataset.count || '0');
        const suffix = el.dataset.suffix || '';
        const value = { val: 0 };
        gsap.to(value, {
          val: target,
          duration: 1.8,
          delay: 0.5 + index * 0.12,
          ease: 'power2.out',
          onUpdate: () => {
            el.textContent = `${Math.ceil(value.val)}${suffix}`;
          },
        });
      });
    }, section);

    return () => ctx.revert();
  }, [sectionRef]);
}
