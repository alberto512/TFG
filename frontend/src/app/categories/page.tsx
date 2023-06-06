'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import CategoriesPage from '@templates/Categories/CategoriesPage';
import Default from '@layouts/Default/Default';

export default function Categories() {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') === null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <CategoriesPage />
    </Default>
  );
}
