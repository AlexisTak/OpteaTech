export type ProjectCategory = 'web' | 'logiciel' | 'ia' | 'conseil';

export interface Project {
  title: string;
  slug: string;
  shortDescription: string;
  coverImageUrl: string;
  tags: string[];
  category: ProjectCategory;
  featured?: boolean;
  projectUrl?: string;
}

export interface Service {
  name: string;
  slug: string;
  description: string;
  icon: 'globe' | 'blocks' | 'bot' | 'cpu';
  color: string;
  features: string[];
  startingPrice?: number;
}

export interface Testimonial {
  clientName: string;
  clientRole: string;
  clientCompany: string;
  content: string;
  rating: number;
}

export const services: Service[] = [
  {
    name: 'Sites Web & Vitrines',
    slug: 'sites-web',
    description:
      'Création de sites modernes, rapides et SEO-optimisés. De la landing page au site e-commerce complet.',
    icon: 'globe',
    color: 'var(--accent)',
    features: ['Next.js / Astro / WordPress', 'Design responsive premium', 'Optimisation SEO & Core Web Vitals'],
    startingPrice: 1900,
  },
  {
    name: 'Logiciels sur mesure',
    slug: 'logiciels-sur-mesure',
    description: 'Développement d’applications web et back-office adaptés à vos processus métier.',
    icon: 'blocks',
    color: 'var(--accent-alt)',
    features: ['Architecture scalable', 'API REST / GraphQL', 'Interface admin intégrée'],
    startingPrice: 4200,
  },
  {
    name: 'Solutions IA',
    slug: 'solutions-ia',
    description:
      'Intégration de solutions d’IA pour automatiser et améliorer vos processus sans complexifier vos équipes.',
    icon: 'bot',
    color: 'var(--accent)',
    features: ['Chatbots & agents IA', 'Fine-tuning de modèles', 'Détection & analyse intelligente'],
    startingPrice: 3500,
  },
  {
    name: 'Conseil & Architecture',
    slug: 'conseil-architecture',
    description:
      'Audit, conseil technique et accompagnement pour structurer et faire évoluer vos projets.',
    icon: 'cpu',
    color: 'var(--accent-alt)',
    features: ['Audit de code & sécurité', 'Choix de stack & architecture', 'Accompagnement équipes'],
    startingPrice: 800,
  },
];

export const projects: Project[] = [
  {
    title: 'Portail citoyen municipal',
    slug: 'portail-citoyen-municipal',
    shortDescription: 'Un portail web pour centraliser les démarches administratives et fluidifier la relation usager.',
    coverImageUrl: 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?auto=format&fit=crop&w=1200&q=80',
    tags: ['Next.js', 'Go', 'PostgreSQL'],
    category: 'web',
    featured: true,
  },
  {
    title: 'Back-office logistique B2B',
    slug: 'backoffice-logistique-b2b',
    shortDescription: 'Outil métier de suivi des flux, planning et alerting en temps réel.',
    coverImageUrl: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?auto=format&fit=crop&w=1200&q=80',
    tags: ['React', 'Fiber', 'Redis'],
    category: 'logiciel',
  },
  {
    title: 'Assistant IA support client',
    slug: 'assistant-ia-support-client',
    shortDescription: 'Agent IA connecté à la base documentaire interne avec escalade humaine intégrée.',
    coverImageUrl: 'https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1200&q=80',
    tags: ['LLM', 'RAG', 'Supabase'],
    category: 'ia',
    featured: true,
  },
  {
    title: 'Audit architecture SaaS',
    slug: 'audit-architecture-saas',
    shortDescription: 'Refonte du design système et plan de scalabilité pour une startup en hypercroissance.',
    coverImageUrl: 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?auto=format&fit=crop&w=1200&q=80',
    tags: ['Architecture', 'Performance', 'Cloud'],
    category: 'conseil',
  },
];

export const testimonials: Testimonial[] = [
  {
    clientName: 'Camille R.',
    clientRole: 'Directrice',
    clientCompany: 'Association Horizon',
    content:
      'Une collaboration limpide et efficace. Notre nouveau site est rapide, clair et enfin simple à maintenir.',
    rating: 5,
  },
  {
    clientName: 'Thomas L.',
    clientRole: 'Co-fondateur',
    clientCompany: 'FlowRetail',
    content:
      'Très bon niveau technique, décisions pragmatiques, et des délais tenus. On se sent accompagnés sans blabla.',
    rating: 5,
  },
  {
    clientName: 'Sarah M.',
    clientRole: 'Cheffe de projet',
    clientCompany: 'Collectivité Nova',
    content: 'Le back-office sur mesure a réduit nos tâches manuelles de façon visible dès les premières semaines.',
    rating: 4,
  },
];
