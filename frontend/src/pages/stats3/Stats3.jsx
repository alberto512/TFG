import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  AreaChart,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Area,
} from 'recharts';
import './Stats3.css';

const Stats3 = () => {
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [accounts, setAccounts] = useState([]);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [categories, setCategories] = useState([]);
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [error, setError] = useState(false);
  const [balances, setBalances] = useState([]);

  const getAcounts = useCallback(() => {
    const getData = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query { accounts() {
            id,
            iban,
            transactions {
              id,
              date,
              amount,
              category {
                id,
                name,
              },
            },
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
    setLoadingCategories(true);
    getAcounts();
    getCategories();
  }, [getAcounts, getCategories]);

  const handleCategory = (id) => {
    if (selectedCategories.includes(id)) {
      setSelectedCategories((prev) => prev.filter((item) => item !== id));
    } else {
      setSelectedCategories((prev) => [...prev, id]);
    }
  };

  const getBalances = async () => {
    var d = new Date();
    const initDate = d.setDate(d.getDate() - 90);
    const endDate = new Date().getTime();

    const transactions = [];

    for (let i in accounts) {
      for (let j in accounts[i].transactions) {
        if (
          accounts[i].transactions[j].category.id === selectedCategories[0] ||
          accounts[i].transactions[j].category.id === selectedCategories[1]
        ) {
          transactions.push(accounts[i].transactions[j]);
        }
      }
    }

    const interval = Math.floor((endDate - initDate) / 4);

    const dateValues = [
      initDate,
      initDate + interval,
      initDate + interval * 2,
      initDate + interval * 3,
      endDate,
    ];

    let data = dateValues.map((value) => {
      const dataObj = {
        date: value,
      };
      dataObj[
        categories.find(
          (category) => category.id === selectedCategories[0]
        ).name
      ] = 0;
      dataObj[
        categories.find(
          (category) => category.id === selectedCategories[1]
        ).name
      ] = 0;
      return dataObj;
    });

    for (let i = 1; i < data.length; i++) {
      data[i][
        categories.find(
          (category) => category.id === selectedCategories[0]
        ).name
      ] =
        data[i - 1][
          categories.find(
            (category) => category.id === selectedCategories[0]
          ).name
        ];
      data[i][
        categories.find(
          (category) => category.id === selectedCategories[1]
        ).name
      ] =
        data[i - 1][
          categories.find(
            (category) => category.id === selectedCategories[1]
          ).name
        ];
      for (let j in transactions) {
        if (
          transactions[j].date >= data[i - 1].date &&
          transactions[j].date < data[i].date
        ) {
          data[i][transactions[j].category.name] += transactions[j].amount;
        }
      }
    }

    data = data.map((item) => ({
      ...item,
      date: new Date(item.date)
        .toLocaleDateString('es-ES', {
          year: '2-digit',
          month: '2-digit',
          day: '2-digit',
        })
        .replace(/\//g, '-'),
    }));

    setBalances(data);
  };

  const handleSubmit = () => {
    if (selectedCategories.length !== 2) {
      setError(true);
      return;
    }
    setError(false);
    getBalances();
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
              Select two categories
            </span>
          </div>
          <div className='scrollers'>
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
            <span className='btn-refresh' onClick={handleSubmit}>
              Next
            </span>
          </div>
        </>
      ) : (
        <div className='charts'>
          <AreaChart width={1000} height={250} data={balances}>
            <defs>
              <linearGradient id='colorCategory1' x1='0' y1='0' x2='0' y2='1'>
                <stop offset='5%' stopColor='#ff4136' stopOpacity={0.8} />
                <stop offset='95%' stopColor='#ff4136' stopOpacity={0} />
              </linearGradient>
              <linearGradient id='colorCategory2' x1='0' y1='0' x2='0' y2='1'>
                <stop offset='5%' stopColor='#82ca9d' stopOpacity={0.8} />
                <stop offset='95%' stopColor='#82ca9d' stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis
              dataKey='date'
              tick={{ stroke: '#6b5b95' }}
              stroke={'#6b5b95'}
            />
            <YAxis
              padding={{ bottom: 15 }}
              tick={{ stroke: '#6b5b95' }}
              stroke={'#6b5b95'}
            />
            <Tooltip />
            <Area
              type='monotone'
              dataKey={
                categories.find(
                  (category) => category.id === selectedCategories[0]
                ).name
              }
              stroke='#ff4136'
              fillOpacity={1}
              fill='url(#colorCategory1)'
            />
            <Area
              type='monotone'
              dataKey={
                categories.find(
                  (category) => category.id === selectedCategories[1]
                ).name
              }
              stroke='#82ca9d'
              fillOpacity={1}
              fill='url(#colorCategory2)'
            />
          </AreaChart>
        </div>
      )}
    </div>
  );
};

export default Stats3;
