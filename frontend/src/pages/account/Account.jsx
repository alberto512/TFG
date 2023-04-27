import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './Account.css';

const Account = () => {
  const { id } = useParams();
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [account, setAccount] = useState({});

  const getAccount = async () => {
    const response = await axios.post(
      backendUrl,
      {
        query: `query AccountById($id: ID!) {
          accountById(id: $id) {
            id,
            iban,
            amount,
            currency,
            transactions {
              id,
              description,
              amount,
              date,
            },
          }
        }`,
        variables: {
          id,
        },
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    setAccount(response.data.data.accountById);
    console.log('setAccount');
  };

  useEffect(() => {
    getAccount();
  }, []);

  return (
    <div className='wrapper'>
      {Object.keys(account).length === 0 ? (
        <FontAwesomeIcon className='spinner' icon='fa-solid fa-spinner' spin />
      ) : (
        <>
          <div className='title-account-wrapper'>
            <span className='title-account'>{account.iban}</span>
            <span
              className={
                account.amount <= 0
                  ? 'amount-account-negative'
                  : 'amount-account'
              }
            >
              {account.amount} {account.currency}
            </span>
          </div>
          <div className='scroller'>
            {account.transactions.map((transaction) => (
              <div
                className={`transaction-wrapper ${
                  transaction.amount <= 0 ? 'transaction-wrapper-negative' : ''
                }`}
              >
                <span className='description'>{transaction.description}</span>
                <div className='transaction-info'>
                  <span>{transaction.date}</span>
                  <span>
                    {transaction.amount} {account.currency}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default Account;
