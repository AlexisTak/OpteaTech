'use client';

import { RefObject, useEffect } from 'react';
import { gsap } from '@/lib/gsap';

interface ProcessAnimationRefs {
  sectionRef: RefObject<HTMLElement | null>;
}

export function useProcessAnimation({ sectionRef }: ProcessAnimationRefs) {
  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const path = section.querySelector<SVGPathElement>('.process-path');
    if (!path) return;

    const ctx = gsap.context(() => {
      const totalLength = path.getTotalLength();
      gsap.set(path, { strokeDasharray: totalLength, strokeDashoffset: totalLength });
      gsap.set('.process-dot', { scale: 0, opacity: 0 });
      gsap.set('.process-step', { opacity: 0, y: 20 });

      const tl = gsap.timeline({
        scrollTrigger: {
          trigger: '.process-section',
          start: 'top 70%',
          end: 'bottom 60%',
          scrub: 1.5,
        },
      });

      tl.to(path, { strokeDashoffset: 0, ease: 'none' });

      [1, 2, 3, 4].forEach((step, index) => {
        const position = [0.2, 0.45, 0.68, 0.88][index];
        tl.fromTo(`.process-dot-${step}`, { scale: 0, opacity: 0 }, { scale: 1, opacity: 1, duration: 0.15, ease: 'opteaSnap' }, position)
          .fromTo(`.process-step-${step}`, { opacity: 0, y: 16 }, { opacity: 1, y: 0, duration: 0.2, ease: 'opteaOut' }, position + 0.02);
      });
    }, section);

    return () => ctx.revert();
  }, [sectionRef]);
}
