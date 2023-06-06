import { NextResponse } from 'next/server';
import { updateTransaction } from '@/lib/Transactions';

export async function PUT(request: Request) {
  const data = await request.json();
  const token = request.headers.get('authorization') || '';
  const response = await updateTransaction(token, data.id, data.categoryId);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json({ data: null });
}
