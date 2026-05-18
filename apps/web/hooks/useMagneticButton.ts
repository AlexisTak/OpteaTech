'use client';

import { RefObject, useEffect } from 'react';
import { gsap } from '@/lib/gsap';

interface MagneticButtonOptions {
  elementRef: RefObject<HTMLElement | null>;
  strength?: number;
}

export function useMagneticButton({ elementRef, strength = 0.35 }: MagneticButtonOptions) {
  useEffect(() => {
    const element = elementRef.current;
    if (!element) return;

    const onMove = (event: MouseEvent) => {
      const rect = element.getBoundingClientRect();
      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;
      const dx = (event.clientX - centerX) * strength;
      const dy = (event.clientY - centerY) * strength;

      gsap.to(element, { x: dx, y: dy, duration: 0.5, ease: 'opteaOut' });
      const inner = element.querySelector<HTMLElement>('.btn-inner');
      if (inner) gsap.to(inner, { x: dx * 0.4, y: dy * 0.4, duration: 0.5, ease: 'opteaOut' });
    };

    const onLeave = () => {
      gsap.to(element, { x: 0, y: 0, duration: 0.7, ease: 'opteaSnap' });
      const inner = element.querySelector<HTMLElement>('.btn-inner');
      if (inner) gsap.to(inner, { x: 0, y: 0, duration: 0.6, ease: 'opteaSnap' });
    };

    element.addEventListener('mousemove', onMove);
    element.addEventListener('mouseleave', onLeave);

    return () => {
      element.removeEventListener('mousemove', onMove);
      element.removeEventListener('mouseleave', onLeave);
    };
  }, [elementRef, strength]);
}
