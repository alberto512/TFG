import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';
import axios from 'axios';
import './SantanderLogin.css';

const SantanderLogin = () => {
  const { token } = useAuth();
  const [queryParameters] = useSearchParams();
  const backendUrl = process.env.REACT_APP_BACKEND_URL;
  const santanderUrl =
    process.env.REACT_APP_SANTANDER_URL + 'prestep-authorize';
  const redirectUri = process.env.REACT_APP_FRONTEND_URL + 'santanderLogin/';
  const [isToken, setIsToken] = useState(false);

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
          },
        }
      );
      console.log(response);
      setIsToken(true);
    };

    if (queryParameters.get('code') == null) {
      window.location.href =
        santanderUrl +
        '?redirect_uri=' +
        redirectUri +
        '&response_type=code&client_id=' +
        process.env.REACT_APP_SANTANDER_ID;
    } else {
      if (!isToken) {
        getToken();
      }
    }
  }, [queryParameters, token, isToken, backendUrl, redirectUri, santanderUrl]);

  return <></>;
};

export default SantanderLogin;
