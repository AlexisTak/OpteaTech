export const primaryNavLinks = [
  { href: '/services', label: 'Services' },
  { href: '/projets', label: 'Projets' },
  { href: '/contact', label: 'Contact' },
] as const;

export const footerNavLinks = [
  { href: '/', label: 'Accueil' },
  ...primaryNavLinks,
  { href: '/a-propos', label: 'A propos' },
] as const;

export const footerServiceLinks = [
  { href: '/services#sites-web', label: 'Sites web' },
  { href: '/services#logiciels-sur-mesure', label: 'Logiciels' },
  { href: '/services#solutions-ia', label: 'IA' },
  { href: '/services#conseil-architecture', label: 'Conseil' },
] as const;

export const socialLinks = [
  { href: 'https://www.linkedin.com', label: 'LinkedIn' },
  { href: 'https://github.com', label: 'GitHub' },
] as const;
