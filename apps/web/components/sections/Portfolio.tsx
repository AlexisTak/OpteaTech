'use client';

import Link from 'next/link';
import { useMemo, useRef, useState } from 'react';
import Image from 'next/image';
import type { Project, ProjectCategory } from '@/lib/site-data';
import { usePortfolioAnimation } from '@/hooks/usePortfolioAnimation';

interface PortfolioProps {
  items: Project[];
}

const filters: Array<{ label: string; value: 'all' | ProjectCategory }> = [
  { label: 'Tous', value: 'all' },
  { label: 'Web', value: 'web' },
  { label: 'Logiciel', value: 'logiciel' },
  { label: 'IA', value: 'ia' },
  { label: 'Conseil', value: 'conseil' },
];

export function Portfolio({ items }: PortfolioProps) {
  const sectionRef = useRef<HTMLElement>(null);
  usePortfolioAnimation({ sectionRef });

  const [activeFilter, setActiveFilter] = useState<'all' | ProjectCategory>('all');

  const visibleItems = useMemo(() => {
    if (activeFilter === 'all') return items;
    return items.filter((item) => item.category === activeFilter);
  }, [activeFilter, items]);

  return (
    <section ref={sectionRef} id="projets" className="section-shell" data-reveal>
      <div className="mb-10 flex flex-col gap-8 md:flex-row md:items-end md:justify-between">
        <div>
          <p className="section-label mb-8">
            <span>02 / Projets</span>
          </p>
          <h2 className="display-md">Quelques</h2>
          <h2 className="display-md -mt-3">travaux.</h2>
        </div>
        <Link href="/projets" className="body-md">Voir tout →</Link>
      </div>

      <div className="mb-10 flex flex-wrap gap-2">
          {filters.map((filter) => (
            <button
              key={filter.value}
              type="button"
              onClick={() => setActiveFilter(filter.value)}
              data-filter={filter.value}
              className={`filter-btn tag ${
                activeFilter === filter.value
                  ? 'active'
                  : 'hover:border-[var(--border-strong)]'
              }`}
            >
              {filter.label}
            </button>
          ))}
      </div>

      <div className="portfolio-grid grid gap-6 lg:grid-cols-12">
        {visibleItems.map((project, index) => {
          const isWide = index % 2 === 0;
          return (
            <Link
              key={project.slug}
              href={`/projets#${project.slug}`}
              id={project.slug}
              data-category={project.category}
              className={`project-card ${isWide ? 'lg:col-span-7' : 'lg:col-span-5'} bg-[var(--surface)]`}
            >
              <div className="project-img-wrap overflow-hidden">
                <Image
                  src={project.coverImageUrl}
                  alt={project.title}
                  width={1200}
                  height={900}
                  className="w-full"
                  sizes="(max-width: 1024px) 100vw, 50vw"
                />
                <div className="project-card-overlay">
                  <span className="project-hover-label">Voir le projet →</span>
                </div>
              </div>
              <div className="project-card-body p-4">
                <h3 className="project-card-title text-3xl leading-tight text-[var(--ink)]">{project.title}</h3>
                <p className="project-card-copy mt-2 body-md">{project.shortDescription}</p>
                <p className="project-card-tags label mt-4">{project.tags.join(' · ')}</p>
              </div>
            </Link>
          );
        })}
      </div>
    </section>
  );
}
