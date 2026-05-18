'use client';

import { RefObject, useEffect } from 'react';
import { gsap } from '@/lib/gsap';

interface HeroParallaxRefs {
  sectionRef: RefObject<HTMLElement | null>;
}

export function useHeroParallax({ sectionRef }: HeroParallaxRefs) {
  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const ctx = gsap.context(() => {
      gsap.to('.hero-bg-grid', {
        yPercent: -30,
        ease: 'none',
        scrollTrigger: {
          trigger: '.hero-section',
          start: 'top top',
          end: 'bottom top',
          scrub: true,
        },
      });

      gsap.to('.hero-title-wrap', {
        yPercent: 15,
        ease: 'none',
        scrollTrigger: {
          trigger: '.hero-section',
          start: 'top top',
          end: 'bottom top',
          scrub: 0.8,
        },
      });

      gsap.to('.scroll-indicator', {
        opacity: 0,
        ease: 'none',
        scrollTrigger: {
          trigger: '.hero-section',
          start: 'top top',
          end: '20% top',
          scrub: true,
        },
      });
    }, section);

    return () => ctx.revert();
  }, [sectionRef]);
}
