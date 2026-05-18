"use client";

import type { Service } from '@/lib/site-data';
import { useRef, useState } from 'react';
import { ArrowUpRight } from 'lucide-react';
import { useServicesAnimation } from '@/hooks/useServicesAnimation';

interface ServicesProps {
  items: Service[];
}

export function Services({ items }: ServicesProps) {
  const sectionRef = useRef<HTMLElement>(null);
  useServicesAnimation({ sectionRef });

  const [activeSlug, setActiveSlug] = useState(items[0]?.slug ?? '');
  const activeService = items.find((service) => service.slug === activeSlug) ?? items[0];

  const activateService = (slug: string) => {
    setActiveSlug((current) => (current === slug ? current : slug));
  };

  return (
    <section ref={sectionRef} id="services" className="services-section section-shell section-alt" data-reveal>
      <div className="grid gap-10 lg:grid-cols-12">
        <aside className="services-sticky-col lg:sticky lg:top-32 lg:col-span-4 lg:h-fit">
          <p className="section-label">
            <span>01 / Services</span>
          </p>
          <h2 className="display-md max-w-[10ch]">Tout ce qu&apos;il vous faut.</h2>
          <p className="body-md mt-6 max-w-md">
            De la maquette a la mise en production. Un seul interlocuteur, une execution nette et une vision long terme.
          </p>

          {activeService ? (
            <div className="service-preview mt-10">
              <div className="service-preview-content">
                <p className="label">Service actif</p>
                <h3 className="service-preview-title mt-3 text-2xl text-[var(--ink)]">{activeService.name}</h3>
                <p className="body-md mt-4 max-w-md">{activeService.description}</p>
                <ul className="service-preview-list mt-5 grid gap-2 max-w-md">
                  {activeService.features.slice(0, 2).map((feature) => (
                    <li key={feature}>
                      <span className="text-[var(--accent)]">- </span>
                      {feature}
                    </li>
                  ))}
                </ul>
                {activeService.startingPrice ? (
                  <p className="service-preview-price label mt-5">A partir de {activeService.startingPrice} EUR</p>
                ) : null}
              </div>
            </div>
          ) : null}
        </aside>

        <div className="lg:col-span-7 lg:col-start-6">
          {items.map((service, index) => {
            const isOpen = activeSlug === service.slug;
            const panelId = `service-panel-${service.slug}`;
            return (
              <article
                key={service.slug}
                className={`service-row ${isOpen ? 'is-open' : ''}`}
                data-active={isOpen ? 'true' : 'false'}
              >
                <button
                  type="button"
                  className="w-full text-left"
                  onPointerEnter={() => activateService(service.slug)}
                  onClick={() => activateService(service.slug)}
                  onFocus={() => activateService(service.slug)}
                  aria-expanded={isOpen}
                  aria-controls={panelId}
                >
                  <div className="flex items-center justify-between gap-4">
                    <div className="grid gap-3">
                      <div className="flex items-center gap-4">
                        <span className="service-index label">{String(index + 1).padStart(2, '0')}</span>
                        <span className="tag">{service.slug.split('-')[0]}</span>
                      </div>
                      <h3 className="service-title text-3xl leading-tight text-[var(--ink)]">{service.name}</h3>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className={`service-state label ${isOpen ? 'is-open' : ''}`}>
                        {isOpen ? 'Ouvert' : 'Apercu'}
                      </span>
                      <ArrowUpRight className="service-arrow" size={18} />
                    </div>
                  </div>
                </button>

                {!isOpen ? (
                  <p className="service-teaser body-md mt-3 max-w-2xl">{service.description}</p>
                ) : null}

                <div
                  id={panelId}
                  className="service-details"
                  aria-hidden={!isOpen}
                >
                  <div className="service-details-inner pt-4">
                    <p className="body-md mt-4 max-w-2xl">{service.description}</p>
                    <ul className="mt-4 grid gap-2 pb-2 text-[15px] text-[var(--ink-60)]">
                      {service.features.map((feature) => (
                        <li key={feature}>
                          <span className="text-[var(--accent)]">- </span>
                          {feature}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </article>
            );
          })}
        </div>
      </div>
    </section>
  );
}
