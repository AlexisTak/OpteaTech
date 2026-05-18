import { NextResponse } from 'next/server';
import { z } from 'zod';

const payloadSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  company: z.string().optional(),
  serviceInterest: z.string().min(1),
  budgetRange: z.string().optional(),
  message: z.string().min(20),
  website: z.string().optional(),
});

export async function POST(request: Request) {
  const body = await request.json();
  const parsed = payloadSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json({ message: 'Validation failed' }, { status: 400 });
  }

  if (parsed.data.website && parsed.data.website.trim().length > 0) {
    return NextResponse.json({ status: 'accepted' }, { status: 200 });
  }

  const apiUrl = process.env.GO_API_INTERNAL_URL ?? process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;
  if (!apiUrl) {
    return NextResponse.json({ status: 'ok' }, { status: 202 });
  }

  try {
    const response = await fetch(`${apiUrl}/api/contact`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(parsed.data),
    });

    if (!response.ok) {
      return NextResponse.json({ message: 'Upstream error' }, { status: 502 });
    }

    return NextResponse.json({ status: 'ok' }, { status: 201 });
  } catch {
    return NextResponse.json({ message: 'Service unavailable' }, { status: 503 });
  }
}
