'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import styles from './AccountsPage.module.css';
import Button from '@elements/Button/Button';

const AccountsPage = () => {
  const router = useRouter();
  const [accounts, setAccounts] = useState<Account[]>([]);

  const refreshAccounts = () => {
    axios
      .post('/api/accounts', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((_response) => {
        axios
          .get('/api/accounts', {
            headers: {
              Authorization: localStorage.getItem('jwt'),
            },
          })
          .then((response) => {
            setAccounts(response.data);
          })
          .catch((error) => {
            console.log(error);
            if (error.response.status === 401) {
              localStorage.removeItem('jwt');
              router.push('/login');
            }
          });
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
    axios
      .get('/api/accounts', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setAccounts(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  }, []);

  return (
    <div className={styles.accountsPage}>
      <div className={styles.titleWrapper}>
        <div className={styles.mainTitleWrapper}>
          <span className={styles.title}>Accounts</span>
          <div className={styles.buttonWrapper}>
            <Button label={'Refresh'} onClick={refreshAccounts} />
          </div>
        </div>
        <div className={styles.keysWrapper}>
          <span className={styles.key}>IBAN</span>
          <span className={styles.key}>Type</span>
          <span className={styles.key}>Bank</span>
          <span className={styles.key}>Amount</span>
        </div>
      </div>
      <div className={styles.accountsWrapper}>
        {accounts.map((account, index, accounts) => {
          const style = {
            style: 'currency',
            currency: account.currency,
          };
          return (
            <div
              key={account.id}
              className={`${styles.account} ${index + 1 === accounts.length ? styles.lastAccount : ''}`}
              onClick={() => router.push('/accounts/' + account.id)}
            >
              <span className={`${styles.item} ${styles.iban}`}>{account.iban}</span>
              <span className={styles.item}>{account.name}</span>
              <span className={styles.item}>{account.bank}</span>
              <span className={`${styles.item} ${account.amount >= 0 ? styles.positive : styles.negative}`}>
                {account.amount.toLocaleString('de-DE', style)}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default AccountsPage;
