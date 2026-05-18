import { Portfolio } from '@/components/sections/Portfolio';
import { getProjects } from '@/lib/api';

export default async function ProjetsPage() {
  const projects = await getProjects();

  return (
    <div className="pt-16">
      <Portfolio items={projects} />
    </div>
  );
}
