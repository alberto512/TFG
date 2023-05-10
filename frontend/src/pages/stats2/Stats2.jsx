import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Sankey, Layer, Rectangle } from 'recharts';
import './Stats2.css';

const CustomNode = ({
  x,
  y,
  width,
  height,
  index,
  payload,
  containerWidth,
}) => {
  const isOut = x + width + 6 > containerWidth;
  return (
    <Layer key={`CustomNode${index}`}>
      <Rectangle x={x} y={y} width={width} height={height} fill='#8884d8' />
      <text
        textAnchor={isOut ? 'end' : 'start'}
        x={isOut ? x - 6 : x + width + 6}
        y={y + height / 2}
        fontSize='30'
        fill='#fff'
      >
        {payload.name}
      </text>
      <text
        textAnchor={isOut ? 'end' : 'start'}
        x={isOut ? x - 6 : x + width + 6}
        y={y + height / 2 + 30}
        fontSize='20'
        fill='#fff'
      >
        {`${payload.value}`}
      </text>
    </Layer>
  );
};

const Stats2 = () => {
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [loadingAccounts, setLoadingAccounts] = useState(false);
  const [accounts, setAccounts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedAccount, setSelectedAccount] = useState(null);
  const [error, setError] = useState(false);
  const [balancesPositive, setBalancesPositive] = useState(null);
  const [balancesNegative, setBalancesNegative] = useState(null);
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

      setCategories(
        response.data.data.categories.map((category) => category.id)
      );
    };

    getData();
  }, [token, backendUrl]);

  useEffect(() => {
    setLoadingAccounts(true);
    getAcounts();
    getCategories();
  }, [getAcounts, getCategories]);

  const handleAccount = (id) => {
    if (id === selectedAccount) {
      setSelectedAccount(null);
    } else {
      setSelectedAccount(id);
    }
  };

  const getBalances = async () => {
    var d = new Date();
    d.setDate(d.getDate() - 180);

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
          accountIds: [selectedAccount],
          categoryIds: categories,
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

    let balancesPositive = {
      nodes: [
        {
          name: accounts.find((account) => account.id === selectedAccount).iban,
        },
      ],
      links: [],
    };

    let balancesNegative = {
      nodes: [
        {
          name: accounts.find((account) => account.id === selectedAccount).iban,
        },
      ],
      links: [],
    };

    for (const i in response.data.data.balances) {
      if (response.data.data.balances[i].amount > 0) {
        const nodes = balancesPositive.nodes;
        const links = balancesPositive.links;
        if (
          nodes.some(
            (node) => node.name === response.data.data.balances[i].category.name
          )
        ) {
          links.push({
            source: 0,
            target: nodes.indexOf(
              nodes.find(
                (node) =>
                  node.name === response.data.data.balances[i].category.name
              )
            ),
            value: response.data.data.balances[i].amount,
          });
        } else {
          nodes.push({
            name: response.data.data.balances[i].category.name,
          });
          links.push({
            source: 0,
            target: nodes.length - 1,
            value: response.data.data.balances[i].amount,
          });
        }

        balancesPositive = {
          nodes: nodes,
          links: links,
        };
      } else if (response.data.data.balances[i].amount < 0) {
        const nodes = balancesNegative.nodes;
        const links = balancesNegative.links;
        if (
          nodes.some(
            (node) => node.name === response.data.data.balances[i].category.name
          )
        ) {
          links.push({
            source: 0,
            target: nodes.indexOf(
              nodes.find(
                (node) =>
                  node.name === response.data.data.balances[i].category.name
              )
            ),
            value: Math.abs(response.data.data.balances[i].amount),
          });
        } else {
          nodes.push({
            name: response.data.data.balances[i].category.name,
          });
          links.push({
            source: 0,
            target: nodes.length - 1,
            value: Math.abs(response.data.data.balances[i].amount),
          });
        }

        balancesNegative = {
          nodes: nodes,
          links: links,
        };
      }
    }

    setBalancesPositive(balancesPositive);

    setBalancesNegative(balancesNegative);

    setLoadingData(false);
  };

  const handleSubmit = () => {
    if (!selectedAccount) {
      setError(true);
      return;
    }
    setError(false);
    setLoadingData(true);
    getBalances();
  };

  return (
    <div className='wrapper'>
      {!balancesPositive && !balancesNegative ? (
        <>
          <div className='title-stats-wrapper'>
            <span className='title-stats'>Stats</span>
            <span
              className={`title-description ${
                error ? 'title-description-error' : ''
              }`}
            >
              Select an account
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
                      selectedAccount && selectedAccount === account.id
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
          {balancesNegative ? (
            <div className='chart'>
              <span className='title-chart'>Negative balances</span>
              <Sankey
                width={1000}
                height={400}
                data={balancesNegative}
                nodePadding={50}
                node={<CustomNode containerWidth={100} />}
                link={{ stroke: '#dd002b' }}
              ></Sankey>
            </div>
          ) : null}
          {balancesPositive ? (
            <div className='chart'>
              <span className='title-chart'>Postive balances</span>
              <Sankey
                width={1000}
                height={400}
                data={balancesPositive}
                nodePadding={50}
                node={<CustomNode containerWidth={100} />}
                link={{ stroke: '#00ddb2' }}
              ></Sankey>
            </div>
          ) : null}
        </div>
      )}
    </div>
  );
};

export default Stats2;
