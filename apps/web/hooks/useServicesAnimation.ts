'use client';

import { RefObject, useEffect } from 'react';
import { gsap, ScrollTrigger } from '@/lib/gsap';

interface ServicesAnimationRefs {
  sectionRef: RefObject<HTMLElement | null>;
}

export function useServicesAnimation({ sectionRef }: ServicesAnimationRefs) {
  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const ctx = gsap.context(() => {
      ScrollTrigger.batch('.service-row', {
        onEnter: (elements) => {
          gsap.fromTo(
            elements,
            { opacity: 0, x: 40 },
            { opacity: 1, x: 0, duration: 0.6, ease: 'opteaOut', stagger: 0.08 },
          );
        },
        start: 'top 90%',
        once: true,
      });
    }, section);

    return () => ctx.revert();
  }, [sectionRef]);
}
