import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './Accounts.css';

const Accounts = () => {
  const navigate = useNavigate();
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [refreshLoading, setRefreshLoading] = useState(false);
  const [loading, setLoading] = useState(false);
  const [accounts, setAccounts] = useState([]);

  const getAccounts = async () => {
    const response = await axios.post(
      backendUrl,
      {
        query: `query { accounts() {
          id,
          iban,
          name,
          currency,
          amount,
          bank,
        } 
      }`,
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    console.log(response.data.accounts);

    const responseDecoded = JSON.parse(response.data.accounts);

    const accounts = [];

    console.log(responseDecoded);

    for (const account in responseDecoded) {
      accounts.push({
        iban: responseDecoded[account][0].account.iban,
        type: responseDecoded[account][0].account.name,
        currency: responseDecoded[account][0].account.currency,
        amount: responseDecoded[account][0].account.balance.amount,
        transactions: responseDecoded[account][1].account.transactions,
      });
    }

    setAccounts(accounts);
  };

  useEffect(() => {
    setLoading(true);
    getAccounts();
    setLoading(false);
  }, []);

  const refreshData = async () => {
    setRefreshLoading(true);
    const response = await axios.post(
      backendUrl,
      {
        query: `query { refreshBankData }`,
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

    setRefreshLoading(false);
  };

  return (
    <div className='wrapper'>
      <div className='title-accounts-wrapper'>
        <span className='title-accounts'>All accounts</span>
        {!refreshLoading ? (
          <span className='btn-refresh' onClick={refreshData}>
            Refresh data
          </span>
        ) : (
          <FontAwesomeIcon
            className='spinner-accounts'
            icon='fa-solid fa-spinner'
            spin
          />
        )}
      </div>
      <div className='scroller'>
        {loading ? (
          <FontAwesomeIcon
            className='spinner'
            icon='fa-solid fa-spinner'
            spin
          />
        ) : (
          accounts
            .sort((a, b) => a.bank.localeCompare(b.bank))
            .map((account) => (
              <div
                key={account.id}
                className={`account-wrapper ${
                  account.amount <= 0 ? 'account-wrapper-negative' : ''
                }`}
                onClick={() => navigate('/account/' + account.id)}
              >
                <span className='iban'>{account.iban}</span>
                <div className='account-info'>
                  <span>{account.bank}</span>
                  <span>{account.type}</span>
                </div>
                <div className='account-balance'>
                  <span>
                    {account.amount} {account.currency}
                  </span>
                </div>
              </div>
            ))
        )}
      </div>
    </div>
  );
};

export default Accounts;
