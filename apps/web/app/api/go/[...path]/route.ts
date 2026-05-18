import { NextRequest, NextResponse } from 'next/server';

const GO_API_URL = process.env.GO_API_INTERNAL_URL ?? 'http://127.0.0.1:3001';

async function proxy(req: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params;
  const joinedPath = path.join('/');
  const requestUrl = new URL(req.url);
  const targetUrl = `${GO_API_URL}/api/public/${joinedPath}${requestUrl.search}`;

  const body = req.method === 'GET' || req.method === 'HEAD' ? undefined : await req.text();

  try {
    const upstream = await fetch(targetUrl, {
      method: req.method,
      headers: {
        'Content-Type': 'application/json',
        'X-Forwarded-For': req.headers.get('x-forwarded-for') ?? '',
        'X-Real-IP': req.headers.get('x-real-ip') ?? '',
      },
      body,
      redirect: 'manual',
      cache: 'no-store',
    });

    const payload = await upstream.text();
    return new NextResponse(payload, {
      status: upstream.status,
      headers: {
        'Content-Type': upstream.headers.get('content-type') ?? 'application/json',
      },
    });
  } catch {
    return NextResponse.json(
      {
        error: 'Service Go indisponible.',
      },
      {
        status: 503,
      },
    );
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
