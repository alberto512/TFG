import { NextResponse } from 'next/server';
import { signUp } from '@/lib/SignUp';

export async function POST(request: Request) {
  const data = await request.json();
  const response = await signUp(data.username, data.password);
  if (response.data == null) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
  if (response.data.data.createUser.id === '') {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }
  return NextResponse.json(response.data.data.createUser.id);
}
