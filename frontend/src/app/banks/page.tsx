'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import BanksPage from '@templates/Banks/BanksPage';

export default function Banks() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') === null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <BanksPage />
    </Default>
  );
}
