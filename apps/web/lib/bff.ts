export type PublicRequestPayload = {
  client_name: string;
  client_email: string;
  client_company?: string;
  client_phone?: string;
  service_type: 'site_web' | 'logiciel' | 'ia' | 'conseil' | 'autre';
  title: string;
  description: string;
  budget_range?: 'moins_2k' | '2k_5k' | '5k_15k' | '15k_plus';
  deadline?: string;
};

const GO_PROXY_BASE = '/api/go';
const GO_HEALTH_URL = process.env.GO_API_INTERNAL_URL ?? process.env.API_INTERNAL_URL ?? 'http://127.0.0.1:3001';

export function getBffBaseUrl(): string {
  return GO_PROXY_BASE;
}

export async function getBffStatus(): Promise<'online' | 'degraded' | 'offline'> {
  try {
    const response = await fetch(`${GO_HEALTH_URL.replace(/\/$/, '')}/health`, {
      cache: 'no-store',
    });

    if (!response.ok) {
      return 'degraded';
    }

    return 'online';
  } catch {
    return 'offline';
  }
}

export async function submitPublicRequest(payload: PublicRequestPayload): Promise<{ ok: boolean; message: string }> {
  try {
    const response = await fetch(`${getBffBaseUrl()}/requests`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const body = (await response.json().catch(() => ({}))) as { error?: string };
      return {
        ok: false,
        message: body.error ?? 'Request failed. Please try again.',
      };
    }

    return {
      ok: true,
      message: 'Request sent successfully. Check your email for next steps.',
    };
  } catch {
    return {
      ok: false,
      message: 'Go API is unreachable. Verify that backend is running on port 3001.',
    };
  }
}
