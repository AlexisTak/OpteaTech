import { cache } from 'react';
import { projects, services, testimonials, type Project, type ProjectCategory, type Service, type Testimonial } from '@/lib/site-data';

const API_BASE = process.env.GO_API_INTERNAL_URL ?? process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

const getProjectFallback = (category?: ProjectCategory) =>
  category ? projects.filter((project) => project.category === category) : projects;

async function safeJsonFetch<T>(url: string): Promise<T | null> {
  try {
    const response = await fetch(url, { next: { revalidate: 120 } });
    if (!response.ok) return null;
    return (await response.json()) as T;
  } catch {
    return null;
  }
}

const getProjectsCached = cache(async (category?: ProjectCategory): Promise<Project[]> => {
  if (!API_BASE) return getProjectFallback(category);
  const query = category ? `?category=${category}` : '';
  const data = await safeJsonFetch<Project[]>(`${API_BASE}/api/projects${query}`);
  return data ?? getProjectFallback(category);
});

const getServicesCached = cache(async (): Promise<Service[]> => {
  if (!API_BASE) return services;
  const data = await safeJsonFetch<Service[]>(`${API_BASE}/api/services`);
  return data ?? services;
});

const getTestimonialsCached = cache(async (): Promise<Testimonial[]> => {
  if (!API_BASE) return testimonials;
  const data = await safeJsonFetch<Testimonial[]>(`${API_BASE}/api/testimonials`);
  return data ?? testimonials;
});

export async function getProjects(category?: ProjectCategory): Promise<Project[]> {
  return getProjectsCached(category);
}

export async function getServices(): Promise<Service[]> {
  return getServicesCached();
}

export async function getTestimonials(): Promise<Testimonial[]> {
  return getTestimonialsCached();
}
