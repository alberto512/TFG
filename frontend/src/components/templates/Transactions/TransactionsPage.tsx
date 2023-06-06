'use client';
import { use, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import InputSelect from '@elements/InputSelect/InputSelect';
import styles from './TransactionsPage.module.css';
import Loading from '@modules/Loading/Loading';

type Props = {
  id: string;
};

const TransactionsPage = ({ id }: Props) => {
  const [loading, setLoading] = useState(false);
  const [iban, setIban] = useState('');
  const [currency, setCurrency] = useState('');
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const router = useRouter();

  const sortByDate = () => {
    setTransactions((prevState) => {
      prevState.sort((a: Transaction, b: Transaction) => {
        return new Date(b.date).getTime() - new Date(a.date).getTime();
      });
      return [...prevState];
    });
  };

  const sortByDescription = () => {
    setTransactions((prevState) => {
      prevState.sort((a: Transaction, b: Transaction) => {
        return a.description.localeCompare(b.description);
      });
      return [...prevState];
    });
  };

  const sortByCategory = () => {
    setTransactions((prevState) => {
      prevState.sort((a: Transaction, b: Transaction) => {
        return a.category.name.localeCompare(b.category.name);
      });
      return [...prevState];
    });
  };

  const sortByAmount = () => {
    setTransactions((prevState) => {
      prevState.sort((a: Transaction, b: Transaction) => {
        return b.amount - a.amount;
      });
      return [...prevState];
    });
  };

  const handleChangeCategory = (transaction: Transaction, categoryId: string) => {
    axios
      .put(
        '/api/transactions',
        {
          id: transaction.id,
          categoryId: categoryId,
        },
        {
          headers: {
            Authorization: localStorage.getItem('jwt'),
          },
        }
      )
      .then((_response) => {})
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  };

  useEffect(() => {
    setLoading(true);
    axios
      .get('/api/transactions/' + id, {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setIban(response.data.iban);
        setCurrency(response.data.currency);
        const transactions = response.data.transactions.sort((a: Transaction, b: Transaction) => {
          return new Date(b.date).getTime() - new Date(a.date).getTime();
        });
        setTransactions(transactions);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });

    axios
      .get('/api/categories', {
        headers: {
          Authorization: localStorage.getItem('jwt'),
        },
      })
      .then((response) => {
        setCategories(response.data);
      })
      .catch((error) => {
        console.log(error);
        if (error.response.status === 401) {
          localStorage.removeItem('jwt');
          router.push('/login');
        }
      });
  }, []);

  useEffect(() => {
    if (transactions.length > 0 && categories.length > 0) {
      setLoading(false);
    }
  }, [transactions, categories]);

  return (
    <div className={styles.transactionsPage}>
      {loading ? (
        <Loading />
      ) : (
        <>
          <div className={styles.titleWrapper}>
            <span className={styles.title}>Transactions of {iban}</span>
            <div className={styles.keysWrapper}>
              <span className={styles.key} onClick={sortByDate}>
                Date
              </span>
              <span className={styles.key} onClick={sortByDescription}>
                Description
              </span>
              <span className={styles.key} onClick={sortByCategory}>
                Category
              </span>
              <span className={styles.key} onClick={sortByAmount}>
                Amount
              </span>
            </div>
          </div>
          <div className={styles.transactionsWrapper}>
            {transactions.map((transaction, index, transactions) => {
              const style = {
                style: 'currency',
                currency: currency,
              };
              const date = new Date(transaction.date);
              return (
                <div
                  className={`${styles.transaction} ${index + 1 === transactions.length ? styles.lastTransaction : ''}`}
                  key={transaction.id}
                >
                  <span className={`${styles.item} ${styles.date}`}>{date.toLocaleDateString('es')}</span>
                  <span className={styles.item}>{transaction.description}</span>
                  <span className={styles.item}>
                    <InputSelect
                      placeholder='Select a category'
                      options={categories}
                      onChange={(categoryId) => handleChangeCategory(transaction, categoryId)}
                      valueDefault={transaction.category.name}
                    />
                  </span>
                  <span className={`${styles.item} ${transaction.amount >= 0 ? styles.positive : styles.negative}`}>
                    {transaction.amount.toLocaleString('de-DE', style)}
                  </span>
                </div>
              );
            })}
          </div>
        </>
      )}
    </div>
  );
};

export default TransactionsPage;
