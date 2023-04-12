import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import './Dashboard.css';

const Dashboard = () => {
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const [accounts, setAccounts] = useState([]);

  const getAccounts = async () => {
    const response = await axios.post(
      backendUrl,
      {
        query: `query { accountsByToken }`,
      },
      {
        headers: {
          Authorization: token,
          withCredentials: true,
        },
      }
    );

    const responseDecoded = JSON.parse(response.data.data.accountsByToken);

    const accounts = [];

    console.log(responseDecoded);

    for (const account in responseDecoded) {
      accounts.push({
        iban: responseDecoded[account][0].account.iban,
        type: responseDecoded[account][0].account.name,
        currency: responseDecoded[account][0].account.currency,
        ammount: responseDecoded[account][0].account.balance.amount,
        transactions: responseDecoded[account][1].account.transactions,
      });
    }

    setAccounts(accounts);
  };

  useEffect(() => {
    getAccounts();
  }, []);

  return (
    <>
      <h1>Dashboard</h1>

      <div>Authenticated</div>
      <ul>
        {accounts.map((account) => (
          <li key={account.iban}>
            <span>
              {account.iban} | {account.type} | {account.currency} |{' '}
              {account.ammount}
            </span>
            <ul>
              {account.transactions.map((transaction) => (
                <li>
                  <span>
                    {transaction.valueDate} | {transaction.description} |{' '}
                    {transaction.amount.content}
                  </span>
                </li>
              ))}
            </ul>
          </li>
        ))}
      </ul>
    </>
  );
};

export default Dashboard;
