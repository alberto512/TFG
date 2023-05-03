import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { BarChart, Bar, XAxis, YAxis, Tooltip } from 'recharts';
import './Stats.css';

const Stats = () => {
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [loadingAccounts, setLoadingAccounts] = useState(false);
  const [accounts, setAccounts] = useState([]);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [categories, setCategories] = useState([]);
  const [selectedAccounts, setSelectedAccounts] = useState([]);
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [error, setError] = useState(false);
  const [balances, setBalances] = useState([]);
  const [loadingData, setLoadingData] = useState(false);

  const getAcounts = useCallback(() => {
    const getData = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query { accounts() {
            id,
            iban,
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
      setLoadingAccounts(false);
    };

    getData();
  }, [token, backendUrl]);

  const getCategories = useCallback(() => {
    const getData = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query { categories() {
            id,
            name,
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

      setCategories(response.data.data.categories);
      setLoadingCategories(false);
    };

    getData();
  }, [token, backendUrl]);

  useEffect(() => {
    setLoadingAccounts(true);
    setLoadingCategories(true);
    getAcounts();
    getCategories();
  }, [getAcounts, getCategories]);

  const handleAccount = (id) => {
    if (selectedAccounts.includes(id)) {
      setSelectedAccounts((prev) => prev.filter((item) => item !== id));
    } else {
      setSelectedAccounts((prev) => [...prev, id]);
    }
  };

  const handleCategory = (id) => {
    if (selectedCategories.includes(id)) {
      setSelectedCategories((prev) => prev.filter((item) => item !== id));
    } else {
      setSelectedCategories((prev) => [...prev, id]);
    }
  };

  const getBalances = async () => {
    var d = new Date();
    d.setDate(d.getDate() - 90);

    const response = await axios.post(
      backendUrl,
      {
        query: `query Balances($accountIds: [ID!]!, $categoryIds: [ID!]!, $initDate: Int!, $endDate: Int!) { 
        balances(accountIds: $accountIds, categoryIds: $categoryIds, initDate: $initDate, endDate: $endDate) {
            amount,
            category {
                name,
            },
        }
      }`,
        variables: {
          accountIds: selectedAccounts,
          categoryIds: selectedCategories,
          initDate: d.getTime(),
          endDate: new Date().getTime(),
        },
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    setBalances([]);

    for (const i in response.data.data.balances) {
      setBalances((prev) => [
        ...prev,
        {
          amount: response.data.data.balances[i].amount,
          category: response.data.data.balances[i].category.name,
        },
      ]);
    }
    setLoadingData(false);
  };

  const handleSubmit = () => {
    if (selectedAccounts.length === 0 || selectedCategories.length === 0) {
      setError(true);
      return;
    }
    setError(false);
    setLoadingData(true);
    getBalances();
  };

  const renderCustomAxisTick = ({ x, y, payload }) => {
    let icon = '';
    switch (payload.value) {
      case 'Food':
        icon = 'fa-solid fa-bowl-food';
        break;
      case 'Home':
        icon = 'fa-solid fa-house';
        break;
      case 'Lifestyle':
        icon = 'fa-solid fa-heart';
        break;
      case 'Health':
        icon = 'fa-solid fa-heart-pulse';
        break;
      case 'Shopping':
        icon = 'fa-solid fa-cart-shopping';
        break;
      case 'Children':
        icon = 'fa-solid fa-child';
        break;
      case 'Vacation':
        icon = 'fa-solid fa-plane';
        break;
      case 'Education':
        icon = 'fa-solid fa-school';
        break;
      case 'Salary':
        icon = 'fa-solid fa-sack-dollar';
        break;
      default:
        icon = 'fa-solid fa-network-wired';
        break;
    }

    return (
      <FontAwesomeIcon
        x={x - 12}
        y={y + 4}
        width={24}
        height={24}
        icon={icon}
        className='icon-stats'
      />
    );
  };

  return (
    <div className='wrapper'>
      {balances.length === 0 ? (
        <>
          <div className='title-stats-wrapper'>
            <span className='title-stats'>Stats</span>
            <span
              className={`title-description ${
                error ? 'title-description-error' : ''
              }`}
            >
              Select at least one account and one category
            </span>
          </div>
          <div className='scrollers'>
            <div className='scroller'>
              {loadingAccounts ? (
                <FontAwesomeIcon
                  className='spinner'
                  icon='fa-solid fa-spinner'
                  spin
                />
              ) : (
                accounts.map((account) => (
                  <div
                    key={account.id}
                    className={`wrapper-item ${
                      selectedAccounts.includes(account.id)
                        ? 'wrapper-item-selected'
                        : ''
                    }`}
                    onClick={() => handleAccount(account.id)}
                  >
                    <span>{account.iban}</span>
                  </div>
                ))
              )}
            </div>
            <div className='scroller'>
              {loadingCategories ? (
                <FontAwesomeIcon
                  className='spinner'
                  icon='fa-solid fa-spinner'
                  spin
                />
              ) : (
                categories.map((category) => (
                  <div
                    key={category.id}
                    className={`wrapper-item ${
                      selectedCategories.includes(category.id)
                        ? 'wrapper-item-selected'
                        : ''
                    }`}
                    onClick={() => handleCategory(category.id)}
                  >
                    <span>{category.name}</span>
                  </div>
                ))
              )}
            </div>
          </div>
          <div className='footer'>
            {loadingData ? (
              <FontAwesomeIcon
                className='spinner-footer'
                icon='fa-solid fa-spinner'
                spin
              />
            ) : (
              <span className='btn-refresh' onClick={handleSubmit}>
                Next
              </span>
            )}
          </div>
        </>
      ) : (
        <div className='charts'>
          <BarChart width={800} height={800} data={balances}>
            <XAxis
              dataKey='category'
              tick={renderCustomAxisTick}
              stroke='#8884d8'
            />
            <YAxis stroke='#8884d8' />
            <Tooltip />
            <Bar dataKey='amount' barSize={30} fill='#8884d8' />
          </BarChart>
        </div>
      )}
    </div>
  );
};

export default Stats;
