import { useQuery } from '@tanstack/react-query';
import { adminApi } from '@/lib/api/requests.api';

export function RequestsPage() {
  const query = useQuery({
    queryKey: ['requests-list'],
    queryFn: () => adminApi.listRequests({ page: 1, limit: 30 }),
  });

  return (
    <div className="stack">
      <h1 style={{ margin: 0 }}>Demandes</h1>
      <div className="card">
        {query.isLoading ? 'Chargement...' : <pre>{JSON.stringify(query.data, null, 2)}</pre>}
      </div>
    </div>
  );
}
