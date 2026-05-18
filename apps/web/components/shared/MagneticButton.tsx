'use client';

import Link from 'next/link';
import { ReactNode, useRef } from 'react';
import { useMagneticButton } from '@/hooks/useMagneticButton';

interface MagneticButtonProps {
  href: string;
  className?: string;
  children: ReactNode;
}

export function MagneticButton({ href, className = '', children }: MagneticButtonProps) {
  const ref = useRef<HTMLAnchorElement>(null);
  useMagneticButton({ elementRef: ref, strength: 0.35 });

  return (
    <Link
      ref={ref}
      href={href}
      className={`btn-primary btn-magnetic ${className}`.trim()}
    >
      <span className="btn-inner inline-flex items-center gap-2">
        {children}
        <span className="arrow">↗</span>
      </span>
    </Link>
  );
}
