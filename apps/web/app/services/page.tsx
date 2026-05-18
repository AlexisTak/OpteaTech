import { Services } from '@/components/sections/Services';
import { getServices } from '@/lib/api';

export default async function ServicesPage() {
  const services = await getServices();

  return (
    <div className="pt-16">
      <Services items={services} />
    </div>
  );
}
