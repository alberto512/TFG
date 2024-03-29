'use client';
import styles from './SantanderPage.module.css';
import Loading from '@modules/Loading/Loading';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { useSearchParams, useRouter } from 'next/navigation';

const SantanderPage = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isToken, setIsToken] = useState(false);

  const getToken = () => {
    axios
      .get('/api/santander/' + searchParams.get('code') || '', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        if (localStorage.getItem('Authorize') === 'true') {
          localStorage.setItem('Authorize', 'false');
          router.push('/accounts');
        } else {
          localStorage.setItem('Authorize', 'true');

          axios
            .post(
              '/api/santander',
              {},
              {
                headers: {
                  Authorization: 'Bearer ' + response.data,
                  'content-type': 'application/json',
                  accept: 'application/json',
                },
              }
            )
            .then((response) => {
              window.location.href = response.data;
            })
            .catch((error) => {
              console.log(error);
              if (error.response.status === 401) {
                localStorage.removeItem('jwt');
                router.push('/login');
              }
            });
        }
        setIsToken(true);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  useEffect(() => {
    if (searchParams.get('code') == null) {
      window.location.href = process.env.NEXT_PUBLIC_SANTANDER_PRESTEP || '';
    } else {
      if (!isToken) {
        getToken();
      }
    }
  }, []);

  return (
    <div className={styles.santanderPage}>
      <Loading />
    </div>
  );
};

export default SantanderPage;
