import React, { useCallback, useEffect, useState } from 'react';
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

  const getAccounts = useCallback(() => {
    const getData = async () => {
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

      setAccounts(response.data.data.accounts);
      setLoading(false);
    };

    getData();
  }, [token, backendUrl]);

  useEffect(() => {
    setLoading(true);
    getAccounts();
  }, [getAccounts]);

  const refreshData = async () => {
    setRefreshLoading(true);
    await axios.post(
      backendUrl,
      {
        query: `mutation { refreshBankData }`,
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    setRefreshLoading(false);
    setLoading(true);
    getAccounts();
    setLoading(false);
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
                  <span>{account.name}</span>
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
