import { useQuery } from '@tanstack/react-query';
import { adminApi } from '@/lib/api/requests.api';

export function DashboardPage() {
  const dashboard = useQuery({
    queryKey: ['desktop-dashboard'],
    queryFn: adminApi.getDashboard,
    refetchInterval: 30000,
  });

  return (
    <div className="stack">
      <h1 style={{ margin: 0 }}>Dashboard</h1>
      <div className="card">
        {dashboard.isLoading ? 'Chargement...' : <pre>{JSON.stringify(dashboard.data, null, 2)}</pre>}
      </div>
    </div>
  );
}
