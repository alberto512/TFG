import { NextResponse } from 'next/server';
import axios from 'axios';

export async function POST(request: Request) {
  console.log(request.headers.get('authorization'));
  const body = {
    access: {
      accounts: [],
      balances: [],
      transactions: [],
    },
    recurringIndicator: true,
  };

  return axios
    .post(process.env.NEXT_PUBLIC_SANTANDER_AUTHORIZE || '', body, {
      headers: {
        Authorization: request.headers.get('authorization') || '',
        'content-type': 'application/json',
        accept: 'application/json',
      },
    })
    .then((_response) => {})
    .catch((error) => {
      if (error.response.status === 403) {
        window.location.href = error.response.data.redirect_uri;
      } else {
        return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
      }
    });
}
