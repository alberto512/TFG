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

    for (const account in responseDecoded) {
      console.log(responseDecoded[account]);
      accounts.push({
        iban: responseDecoded[account][0].iban,
        type: responseDecoded[account][0].name,
        currency: responseDecoded[account][0].currency,
        ammount: responseDecoded[account][0].balance.amount,
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
          <li>
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
