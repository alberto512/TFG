import { NextResponse } from 'next/server';
import { categories, create } from '@/lib/Categories';

export async function POST(request: Request) {
  const token = request.headers.get('authorization') || '';
  const data = await request.json();
  const response = await create(token, data.name);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(null);
}

export async function GET(request: Request) {
  const token = request.headers.get('authorization') || '';
  const response = await categories(token);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(response.data.data.categories);
}
