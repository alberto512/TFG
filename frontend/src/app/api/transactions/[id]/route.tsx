import { NextResponse } from 'next/server';
import { transactions } from '@/lib/Transactions';

type Props = {
  params: {
    id: string;
  };
};

export async function GET(request: Request, { params }: Props) {
  const token = request.headers.get('authorization') || '';
  const response = await transactions(token, params.id);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(response.data.data.accountById);
}
