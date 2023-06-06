'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import SantanderPage from '@templates/Santander/SantanderPage';

export default function Santander() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') === null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <SantanderPage />
    </Default>
  );
}
