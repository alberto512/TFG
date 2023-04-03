import React from 'react';
import { redirect, useSearchParams } from 'react-router-dom';
import axios from 'axios';
import { Buffer } from 'buffer';
import './RegisterBank.css';

const RegisterBank = () => {
  const [queryParameters] = useSearchParams();
  const authorization = Buffer.from(
    'bc75ee49-9924-4160-904e-6b246d751e2c:U8jK1aX4jI8fU0hO7oX8oA4uC6pJ4nC8cY7aN4cN8pI3nT6lU4'
  ).toString('base64');
  const body2 = {
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

  const auxFunc = () => {
    window.location.href =
      'https://apis-sandbox.bancosantander.es/canales-digitales/sb/prestep-authorize?redirect_uri=https://tfg-app.netlify.app/&response_type=code&client_id=bc75ee49-9924-4160-904e-6b246d751e2c';
  };

  const auxFunc2 = () => {
    const body1 = {
      grant_type: 'authorization_code',
      redirect_uri: 'https://tfg-app.netlify.app/',
      code: queryParameters.get('code'),
    };

    axios
      .post(
        'https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/token',
        body1,
        {
          headers: {
            'X-IBM-Client-Id': 'bc75ee49-9924-4160-904e-6b246d751e2c',
            'X-IBM-Client-Secret':
              'U8jK1aX4jI8fU0hO7oX8oA4uC6pJ4nC8cY7aN4cN8pI3nT6lU4',
            Authorization: authorization,
            'content-type': 'application/x-www-form-urlencoded',
            accept: 'application/json',
          },
        }
      )
      .then((response) => {
        axios
          .post(
            'https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/authorize/?client_id=bc75ee49-9924-4160-904e-6b246d751e2c&redirect_uri=https://tfg-app.netlify.app/&response_type=code',
            body2,
            {
              headers: {
                Authorization: 'Bearer ' + response.data.access_token,
                'content-type': 'application/json',
                accept: 'application/json',
              },
            }
          )
          .then((_response) => {})
          .catch((error) => {
            if (error.response.status === 403) {
              console.log(error.response.data);
              window.location.href = error.response.data.redirect_uri;
            }
          });
      })
      .catch((error) => console.log(error));
  };

  const auxFunc3 = () => {
    const body1 = {
      grant_type: 'authorization_code',
      redirect_uri: 'https://tfg-app.netlify.app/',
      code: queryParameters.get('code'),
    };

    axios
      .post(
        'https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/token',
        body1,
        {
          headers: {
            'X-IBM-Client-Id': 'bc75ee49-9924-4160-904e-6b246d751e2c',
            'X-IBM-Client-Secret':
              'U8jK1aX4jI8fU0hO7oX8oA4uC6pJ4nC8cY7aN4cN8pI3nT6lU4',
            Authorization: authorization,
            'content-type': 'application/x-www-form-urlencoded',
            accept: 'application/json',
          },
        }
      )
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => console.log(error));
  };

  const santanderBank = () => {
    // Redirect to login Santander url (/santanderLogin)
    // If no code redirect to Santander bank
    // If code ask server for token
    // If token is correct redirect to permissions url(/santanderPermissions)
    // If no code, ask server for token (Server should check if token is active and ask for another if not) and redirect to permissions page of santander
    // If code get definitivive tokens
    // If tokens are correct redirect to /dashboard
    // In /dashboard you should see a boton of update and the information of the accounts
    // Order files and folders
    // Create new app of santander and change client id and client secret
    redirect('/santanderLogin');
  };

  return (
    <>
      <button onClick={auxFunc}>Hola</button>
      <button onClick={auxFunc2}>Hola 2</button>
      <button onClick={auxFunc3}>Hola 3</button>
      <button onClick={santanderBank}>Register Santander user</button>
    </>
  );
};

export default RegisterBank;
