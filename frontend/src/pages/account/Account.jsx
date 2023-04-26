import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
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
        query: `query { accountById }`,
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    const responseDecoded = JSON.parse(response.data.data.accountsByToken);

    console.log(responseDecoded);

    setAccount({
      iban: responseDecoded[account][0].account.iban,
      currency: responseDecoded[account][0].account.currency,
      amount: responseDecoded[account][0].account.balance.amount,
      transactions: responseDecoded[account][1].account.transactions,
    });
  };

  useEffect(() => {
    //getAccount();
  }, []);

  return (
    <div className='wrapper'>
      <div className='title-account-wrapper'>
        <span className='title-account'>{account.iban}</span>
        <span
          className={
            account.amount <= 0 ? 'amount-account-negative' : 'amount-account'
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
              <span>{transaction.valueDate}</span>
              <span>
                {transaction.amount} {transaction.currency}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Account;
