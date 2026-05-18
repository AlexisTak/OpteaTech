import { FinalCta } from '@/components/sections/FinalCta';
import { Hero } from '@/components/sections/Hero';
import { Portfolio } from '@/components/sections/Portfolio';
import { Process } from '@/components/sections/Process';
import { Services } from '@/components/sections/Services';
import { Testimonials } from '@/components/sections/Testimonials';
import { getProjects, getServices, getTestimonials } from '@/lib/api';

export default async function HomePage() {
  const [services, projects, testimonials] = await Promise.all([getServices(), getProjects(), getTestimonials()]);

  return (
    <>
      <Hero />
      <Services items={services} />
      <Portfolio items={projects} />
      <Process />
      <Testimonials items={testimonials} />
      <FinalCta />
    </>
  );
}
