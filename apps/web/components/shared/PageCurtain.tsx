'use client';

import { useEffect, useRef } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { gsap } from '@/lib/gsap';

export function PageCurtain() {
  const curtainRef = useRef<HTMLDivElement>(null);
  const isFirstRender = useRef(true);
  const isTransitioning = useRef(false);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const curtain = curtainRef.current;
    if (!curtain) return;

    if (isFirstRender.current) {
      isFirstRender.current = false;
      gsap.set(curtain, { scaleY: 0, transformOrigin: 'bottom', pointerEvents: 'none' });
      return;
    }

    gsap.killTweensOf(curtain);
    gsap.set(curtain, { transformOrigin: 'top', pointerEvents: 'none' });
    gsap.to(curtain, {
      scaleY: 0,
      duration: 0.8,
      ease: 'opteaInOut',
      onComplete: () => {
        isTransitioning.current = false;
      },
    });
  }, [pathname]);

  useEffect(() => {
    const onClick = (event: MouseEvent) => {
      if (event.defaultPrevented || event.button !== 0) return;
      if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;

      const target = event.target;
      if (!(target instanceof Element)) return;

      const anchor = target.closest('a[href]');
      const curtain = curtainRef.current;
      if (!(anchor instanceof HTMLAnchorElement) || !curtain) return;
      if (anchor.target && anchor.target !== '_self') return;
      if (anchor.hasAttribute('download')) return;
      if (anchor.getAttribute('rel')?.includes('external')) return;

      const url = new URL(anchor.href, window.location.href);
      if (url.origin !== window.location.origin) return;

      const nextHref = `${url.pathname}${url.search}${url.hash}`;
      const currentHref = `${window.location.pathname}${window.location.search}${window.location.hash}`;
      if (nextHref === currentHref) return;
      if (url.pathname === window.location.pathname && url.search === window.location.search && url.hash) return;
      if (isTransitioning.current) {
        event.preventDefault();
        return;
      }

      isTransitioning.current = true;
      event.preventDefault();

      gsap.killTweensOf(curtain);
      gsap.set(curtain, { pointerEvents: 'auto', transformOrigin: 'bottom' });
      gsap.to(curtain, {
        scaleY: 1,
        duration: 0.65,
        ease: 'opteaInOut',
        onComplete: () => {
          router.push(nextHref);
        },
      });
    };

    document.addEventListener('click', onClick, true);
    return () => document.removeEventListener('click', onClick, true);
  }, [router]);

  return (
    <div
      ref={curtainRef}
      id="page-curtain"
      style={{
        position: 'fixed',
        inset: 0,
        background: 'var(--ink)',
        zIndex: 10000,
        transform: 'scaleY(0)',
        transformOrigin: 'bottom',
        pointerEvents: 'none',
      }}
      aria-hidden="true"
    />
  );
}
