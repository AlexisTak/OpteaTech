import Link from 'next/link';
import { footerNavLinks, footerServiceLinks, socialLinks } from '@/lib/navigation';

export function Footer() {
  return (
    <footer className="mx-auto w-full max-w-[1440px] px-6 pb-12 pt-20 md:px-12">
      <div className="grid gap-10 md:grid-cols-12">
        <div className="md:col-span-4">
          <p className="text-[20px] tracking-tight text-[var(--ink)]">
            optea<span className="text-[var(--accent)]">·</span>tech
          </p>
          <p className="font-display mt-4 max-w-sm text-3xl italic text-[var(--ink)]">
            Code artisanal,
            <br />
            resultats durables.
          </p>
        </div>

        <div className="md:col-span-2">
          <p className="label">Navigation</p>
          <ul className="mt-4 space-y-2 body-md">
            {footerNavLinks.map((link) => (
              <li key={link.href}>
                <Link href={link.href}>{link.label}</Link>
              </li>
            ))}
          </ul>
        </div>

        <div className="md:col-span-2">
          <p className="label">Services</p>
          <ul className="mt-4 space-y-2 body-md">
            {footerServiceLinks.map((link) => (
              <li key={link.href}>
                <Link href={link.href}>{link.label}</Link>
              </li>
            ))}
          </ul>
        </div>

        <div className="md:col-span-3 md:col-start-10">
          <p className="label">Contact</p>
          <ul className="mt-4 space-y-2 body-md">
            <li>
              <a href="mailto:hello@opteaetch.fr">hello@opteaetch.fr</a>
            </li>
            {socialLinks.map((link) => (
              <li key={link.href}>
                <a href={link.href} target="_blank" rel="noopener noreferrer">
                  {link.label} ↗
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="divider mt-14" />
      <div className="mt-5 flex flex-wrap items-center justify-between gap-4">
        <p className="label">© 2026 optea·tech · Fait avec soin en France</p>
        <p className="label">
          <Link href="/mentions-legales">Mentions legales</Link> · <Link href="/rgpd">RGPD</Link>
        </p>
      </div>
    </footer>
  );
}
