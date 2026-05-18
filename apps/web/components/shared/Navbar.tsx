'use client';

import Link from 'next/link';
import { Menu, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useNavbarScroll } from '@/hooks/useNavbarScroll';
import { MagneticButton } from '@/components/shared/MagneticButton';
import { primaryNavLinks } from '@/lib/navigation';

export function Navbar() {
  useNavbarScroll();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : '';
    return () => {
      document.body.style.overflow = '';
    };
  }, [open]);

  return (
    <header className="navbar fixed top-0 z-50 w-full border-b border-transparent transition">
      <nav className="grid-shell flex h-16 items-center justify-between">
        <Link href="/" className="text-[15px] font-medium tracking-tight text-[var(--ink)]" aria-label="Accueil optea.tech">
          optea<span className="text-[var(--accent)]">·</span>tech
        </Link>

        <div className="hidden items-center gap-10 md:flex">
          {primaryNavLinks.map((link) => (
            <Link key={link.href} href={link.href} className="nav-link">
              {link.label}
            </Link>
          ))}
        </div>

        <div className="hidden items-center gap-4 md:flex">
          <MagneticButton href="/login" className="text-sm">
            Se connecter
          </MagneticButton>
        </div>

        <button
          className="inline-flex h-10 w-10 items-center justify-center rounded-lg border border-[var(--border)] text-[var(--ink)] md:hidden"
          type="button"
          onClick={() => setOpen((prev) => !prev)}
          aria-expanded={open}
          aria-controls="mobile-navigation"
          aria-label={open ? 'Fermer le menu' : 'Ouvrir le menu'}
        >
          {open ? <X size={18} /> : <Menu size={18} />}
        </button>
      </nav>

      {open ? (
        <div className="grid-shell mt-2 md:hidden">
          <div className="mobile-nav-backdrop rounded-[28px] border border-[var(--border)] bg-[color:rgba(255,255,255,0.92)] p-4 shadow-[0_20px_60px_rgba(14,14,12,0.08)] backdrop-blur-xl">
            <div id="mobile-navigation" className="flex flex-col gap-3">
              {primaryNavLinks.map((link) => (
                <Link key={link.href} href={link.href} className="rounded-2xl px-4 py-3 text-sm text-[var(--ink-60)] transition hover:bg-[var(--surface)] hover:text-[var(--ink)]" onClick={() => setOpen(false)}>
                  {link.label}
                </Link>
              ))}
              <Link href="/contact" className="btn-primary mt-2 w-full justify-center text-sm" onClick={() => setOpen(false)}>
                Démarrer un projet
              </Link>
            </div>
          </div>
        </div>
      ) : null}
    </header>
  );
}
