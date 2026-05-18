"use client";

import type { Testimonial } from '@/lib/site-data';

interface TestimonialsProps {
  items: Testimonial[];
}

export function Testimonials({ items }: TestimonialsProps) {
  return (
    <section className="section-shell" data-reveal>
      <div className="mb-10">
        <p className="section-label mb-8">
          <span>04 / Témoignages</span>
        </p>
        <h2 className="display-md">Ce que disent</h2>
        <h2 className="display-md -mt-3">nos clients.</h2>
      </div>

      <div className="testi-scroll-wrap">
        <div className="testi-track">
          {items.map((item, index) => (
            <article
              key={`${item.clientName}-${item.clientCompany}`}
              className="testimonial-slide card flex-none p-8"
              style={{ ['--i' as string]: index }}
            >
              <p className="font-display text-7xl leading-none text-[color:var(--accent-muted)]">❝</p>
              <p className="mt-2 font-display text-3xl italic leading-snug text-[var(--ink)]">{item.content}</p>
              <div className="divider mt-8" />
              <div className="mt-5 flex items-center justify-between gap-4">
                <p className="body-md">
                  {item.clientName} · {item.clientRole} · {item.clientCompany}
                </p>
                <p className="label">{'★'.repeat(item.rating)}</p>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}


