import React from 'react';
import { useNavigate } from 'react-router-dom';
import './RegisterBank.css';

const RegisterBank = () => {
  const navigate = useNavigate();

  const santanderBank = () => {
    // If tokens are correct redirect to /dashboard
    // In /dashboard you should see a boton of update and the information of the accounts
    // Order files and folders
    // Review imports
    navigate('/santanderLogin');
  };

  return (
    <>
      <button onClick={santanderBank}>Register Santander user</button>
    </>
  );
};

export default RegisterBank;
