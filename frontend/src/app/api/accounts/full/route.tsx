import { NextResponse } from 'next/server';
import { fullAccounts } from '@/lib/Accounts';

export async function GET(request: Request) {
  const token = request.headers.get('authorization') || '';
  const response = await fullAccounts(token);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(response.data.data.accounts);
}
