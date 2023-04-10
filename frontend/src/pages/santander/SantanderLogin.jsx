import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';
import axios from 'axios';
import './SantanderLogin.css';

const SantanderLogin = () => {
  const { token } = useAuth();
  const navigate = useNavigate();
  const [queryParameters] = useSearchParams();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const santanderUrlPrestep =
    process.env.REACT_APP_SANTANDER_URL + 'prestep-authorize';
  const santanderUrlAuthorize =
    process.env.REACT_APP_SANTANDER_URL + 'v2/authorize';
  const redirectUri = process.env.REACT_APP_FRONTEND_URL + 'santanderLogin/';
  const [isToken, setIsToken] = useState(false);

  const body = {
    access: {
      accounts: [],
      balances: [],
      transactions: [],
      cards_accounts: [],
      cards_balances: [],
      cards_transactions: [],
    },
    recurringIndicator: true,
  };

  useEffect(() => {
    const getToken = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query GetTokenWithCode($code: String!) {
            getTokenWithCode(code: $code)
          }`,
          variables: {
            code: queryParameters.get('code'),
          },
        },
        {
          headers: {
            Authorization: token,
            withCredentials: true,
          },
        }
      );
      if (localStorage.getItem('Authorize')) {
        localStorage.setItem('Authorize', false);
        navigate('/dashboard');
      } else {
        localStorage.setItem('Authorize', true);
        await axios
          .post(
            santanderUrlAuthorize +
              '?redirect_uri=' +
              redirectUri +
              '&response_type=code&client_id=' +
              process.env.REACT_APP_SANTANDER_ID,
            body,
            {
              headers: {
                Authorization: 'Bearer ' + response.data.data.getTokenWithCode,
                'content-type': 'application/json',
                accept: 'application/json',
              },
            }
          )
          .then((_response) => {})
          .catch((error) => {
            if (error.response.status === 403) {
              window.location.href = error.response.data.redirect_uri;
            }
          });
      }
      setIsToken(true);
    };

    if (queryParameters.get('code') == null) {
      window.location.href =
        santanderUrlPrestep +
        '?redirect_uri=' +
        redirectUri +
        '&response_type=code&client_id=' +
        process.env.REACT_APP_SANTANDER_ID;
    } else {
      if (!isToken) {
        getToken();
      }
    }
  }, [
    queryParameters,
    token,
    isToken,
    backendUrl,
    redirectUri,
    santanderUrlPrestep,
    santanderUrlAuthorize,
  ]);

  return <></>;
};

export default SantanderLogin;
