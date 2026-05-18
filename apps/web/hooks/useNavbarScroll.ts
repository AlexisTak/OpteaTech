'use client';

import { useEffect } from 'react';
import { gsap, ScrollTrigger } from '@/lib/gsap';

export function useNavbarScroll() {
  useEffect(() => {
    const nav = document.querySelector<HTMLElement>('.navbar');
    if (!nav) return;

    const ctx = gsap.context(() => {
      const bgTrigger = ScrollTrigger.create({
        start: 'top -80px',
        onUpdate: (self) => {
          if (self.progress > 0) {
            gsap.to(nav, {
              backdropFilter: 'blur(20px)',
              borderBottomColor: 'rgba(14,14,12,0.1)',
              backgroundColor: 'rgba(250,250,249,0.92)',
              duration: 0.4,
              ease: 'opteaOut',
            });
          } else {
            gsap.to(nav, {
              backdropFilter: 'blur(0px)',
              borderBottomColor: 'transparent',
              backgroundColor: 'transparent',
              duration: 0.4,
              ease: 'opteaOut',
            });
          }
        },
      });

      let lastY = 0;
      const hideTrigger = ScrollTrigger.create({
        onUpdate: (self) => {
          const currentY = self.scroll();
          const isDown = currentY > lastY;

          if (currentY > 200) {
            gsap.to(nav, {
              y: isDown ? -80 : 0,
              duration: isDown ? 0.35 : 0.5,
              ease: isDown ? 'opteaIn' : 'opteaOut',
            });
          } else {
            gsap.to(nav, { y: 0, duration: 0.25, ease: 'opteaOut' });
          }
          lastY = currentY;
        },
      });

      return () => {
        bgTrigger.kill();
        hideTrigger.kill();
      };
    });

    return () => ctx.revert();
  }, []);
}
