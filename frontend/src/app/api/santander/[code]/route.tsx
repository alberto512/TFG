import { NextResponse } from 'next/server';
import { tokenCode } from '@/lib/Santander';

type Props = {
  params: {
    code: string;
  };
};

export async function GET(request: Request, { params }: Props) {
  const token = request.headers.get('authorization') || '';
  const response = await tokenCode(token, params.code);
  if (response.code === 401) {
    return NextResponse.json({ error: 'Invalid token' }, { status: 401 });
  }
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  return NextResponse.json(response.data.tokenWithCode);
}
