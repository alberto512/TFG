import React, { useEffect } from 'react';
import axios from 'axios';
import { useAuth } from 'components/authProvider/AuthProvider';
import './Dashboard.css';

const Dashboard = () => {
  const { token } = useAuth();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;

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

    console.log(response);
  };

  useEffect(() => {
    getAccounts();
  }, []);

  return (
    <>
      <h1>Dashboard</h1>

      <div>Authenticated</div>
    </>
  );
};

export default Dashboard;
