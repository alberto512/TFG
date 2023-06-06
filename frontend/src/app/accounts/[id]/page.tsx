'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Default from '@layouts/Default/Default';
import TransactionsPage from '@templates/Transactions/TransactionsPage';

type Props = {
  params: {
    id: string;
  };
};

export default function Transactions({ params }: Props) {
  const router = useRouter();

  useEffect(() => {
    if (localStorage.getItem('jwt') === null) {
      router.push('/');
    }
  }, []);

  return (
    <Default>
      <TransactionsPage id={params.id} />
    </Default>
  );
}
