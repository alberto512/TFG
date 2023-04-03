import React from 'react';
import { useSearchParams } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';
import axios from 'axios';
import './SantanderLogin.css';

const SantanderLogin = async () => {
  const { token } = useAuth();
  const [queryParameters] = useSearchParams();
  const backendUrl = process.env.BACKEND_URL;
  const santanderUrl = process.env.SANTANDER_URL + 'prestep-authorize';
  const redirectUri = process.env.FRONTEND_URL + 'santanderLogin/';

  if (queryParameters.get('code') == null) {
    window.location.href =
      santanderUrl +
      '?redirect_uri=' +
      redirectUri +
      '&response_type=code&client_id=' +
      process.env.SANTANDER_ID;
  } else {
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
  }

  return <></>;
};

export default SantanderLogin;
