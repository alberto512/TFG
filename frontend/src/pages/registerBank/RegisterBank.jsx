import React from 'react';
import { useNavigate } from 'react-router-dom';
import './RegisterBank.css';

const RegisterBank = () => {
  const navigate = useNavigate();

  const santanderBank = () => {
    navigate('/santanderLogin');
  };

  return (
    <div className='wrapper'>
      <div className='item' onClick={santanderBank}>
        Register Santander user
      </div>
    </div>
  );
};

export default RegisterBank;
