import { balances } from '@/lib/Stats';
import { NextResponse } from 'next/server';

export async function POST(request: Request) {
  const data = await request.json();
  const token = request.headers.get('authorization') || '';
  const response = await balances(token, data.selectedAccounts, data.selectedCategories, data.initDate, data.endDate);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(response.data.data.balances);
}
