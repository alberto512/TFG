import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { PieChart, Pie, Cell, Legend } from 'recharts';
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
  const [balancesPositive, setBalancesPositive] = useState([]);
  const [balancesNegative, setBalancesNegative] = useState([]);
  const [loadingData, setLoadingData] = useState(false);

  const COLORS_NEGATIVE = [
    '#F1C40F',
    '#27AE60',
    '#9B59B6',
    '#2ECC71',
    '#8E44AD',
    '#F44336',
    '#16A085',
    '#FF9800',
    '#2980B9',
    '#3498DB',
  ];

  const COLORS_POSITIVE = [
    '#3366FF',
    '#FF69B4',
    '#FF6347',
    '#FFA07A',
    '#1E90FF',
    '#FF5733',
    '#7B68EE',
    '#8B008B',
    '#00FF7F',
    '#FFD700',
  ];

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

    setBalancesNegative([]);
    setBalancesPositive([]);

    for (const i in response.data.data.balances) {
      if (response.data.data.balances[i].amount >= 0) {
        setBalancesPositive((prev) => [
          ...prev,
          {
            amount: response.data.data.balances[i].amount,
            category: response.data.data.balances[i].category.name,
          },
        ]);
      } else {
        setBalancesNegative((prev) => [
          ...prev,
          {
            amount: Math.abs(response.data.data.balances[i].amount),
            category: response.data.data.balances[i].category.name,
          },
        ]);
      }
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

  return (
    <div className='wrapper'>
      {balancesPositive.length === 0 && balancesNegative.length === 0 ? (
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
          {balancesNegative.length !== 0 ? (
            <div className='chart'>
              <span className='title-chart'>Negative balances</span>
              <PieChart width={400} height={400}>
                <Pie
                  data={balancesNegative}
                  dataKey='amount'
                  nameKey='category'
                  innerRadius={50}
                  outerRadius={80}
                  label
                >
                  {balancesNegative.map((_, index) => (
                    <Cell key={`${index}`} fill={COLORS_NEGATIVE[index % 10]} />
                  ))}
                </Pie>
                <Legend />
              </PieChart>
            </div>
          ) : null}
          {balancesPositive.length !== 0 ? (
            <div className='chart'>
              <span className='title-chart'>Postive balances</span>
              <PieChart width={400} height={400}>
                <Pie
                  data={balancesPositive}
                  dataKey='amount'
                  nameKey='category'
                  innerRadius={50}
                  outerRadius={80}
                  label
                >
                  {balancesPositive.map((_, index) => (
                    <Cell key={`${index}`} fill={COLORS_POSITIVE[index % 10]} />
                  ))}
                </Pie>
                <Legend />
              </PieChart>
            </div>
          ) : null}
        </div>
      )}
    </div>
  );
};

export default Stats;
