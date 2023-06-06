'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import LoginPage from '@templates/Login/LoginPage';

export default function Login() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') !== null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <LoginPage />
    </Default>
  );
}
