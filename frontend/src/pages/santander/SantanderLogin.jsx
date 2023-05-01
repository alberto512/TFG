import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
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

  useEffect(() => {
    const getToken = async () => {
      const response = await axios.post(
        backendUrl,
        {
          query: `query TokenWithCode($code: String!) {
            tokenWithCode(code: $code)
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

      if (JSON.parse(localStorage.getItem('Authorize'))) {
        localStorage.setItem('Authorize', false);
        navigate('/accounts');
      } else {
        localStorage.setItem('Authorize', true);

        const body = {
          access: {
            accounts: [],
            balances: [],
            transactions: [],
          },
          recurringIndicator: true,
        };

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
                Authorization: 'Bearer ' + response.data.data.tokenWithCode,
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
    navigate,
  ]);

  return (
    <div className='wrapper'>
      <FontAwesomeIcon className='spinner' icon='fa-solid fa-spinner' spin />
    </div>
  );
};

export default SantanderLogin;
