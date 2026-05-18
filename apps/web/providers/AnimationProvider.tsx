'use client';

import { ReactNode, useEffect } from 'react';
import { gsap, ScrollTrigger } from '@/lib/gsap';

interface AnimationProviderProps {
  children: ReactNode;
}

export function AnimationProvider({ children }: AnimationProviderProps) {
  useEffect(() => {
    ScrollTrigger.refresh();

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (reduced) {
      gsap.globalTimeline.timeScale(100);
      ScrollTrigger.getAll().forEach((trigger) => trigger.kill());
      return;
    }

    return () => {
      ScrollTrigger.getAll().forEach((trigger) => trigger.kill());
    };
  }, []);

  return <>{children}</>;
}
