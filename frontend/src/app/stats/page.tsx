'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import StatsPage from '@templates/Stats/StatsPage';

export default function Stats() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') === null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <StatsPage />
    </Default>
  );
}
