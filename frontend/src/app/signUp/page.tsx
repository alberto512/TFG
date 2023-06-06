'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import SignUpPage from '@templates/SignUp/SignUpPage';

export default function SignUp() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') !== null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <SignUpPage />
    </Default>
  );
}
