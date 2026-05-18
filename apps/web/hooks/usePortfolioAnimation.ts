'use client';

import { RefObject, useEffect } from 'react';
import { gsap } from '@/lib/gsap';

interface PortfolioAnimationRefs {
  sectionRef: RefObject<HTMLElement | null>;
}

export function usePortfolioAnimation({ sectionRef }: PortfolioAnimationRefs) {
  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const ctx = gsap.context(() => {
      gsap.utils.toArray<HTMLElement>('.project-img-wrap').forEach((wrap) => {
        const img = wrap.querySelector('img');
        if (!img) return;

        gsap.fromTo(
          img,
          { yPercent: -8 },
          {
            yPercent: 8,
            ease: 'none',
            scrollTrigger: {
              trigger: wrap,
              start: 'top bottom',
              end: 'bottom top',
              scrub: 1.2,
            },
          },
        );
      });

      gsap.fromTo(
        '.project-card',
        { opacity: 0, y: 60 },
        {
          opacity: 1,
          y: 0,
          duration: 0.9,
          ease: 'opteaOut',
          stagger: { each: 0.15, from: 'start' },
          scrollTrigger: {
            trigger: '.portfolio-grid',
            start: 'top 80%',
            once: true,
          },
        },
      );
    }, section);

    return () => ctx.revert();
  }, [sectionRef]);
}
