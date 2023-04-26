import React, { useState } from 'react';
import { useAuth } from 'components/authProvider/AuthProvider';
import './Register.css';

const Register = () => {
  const { onRegister } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const registerHandler = () => {
    onRegister(username, password, 'USER');
  };

  return (
    <div className='wrapper'>
      <div className='form-container'>
        <div className='form-wrapper'>
          <span className='form-label'>Username</span>
          <input
            className='input-label'
            type='text'
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className='form-wrapper'>
          <span className='form-label'>Password</span>
          <input
            className='input-label'
            type='password'
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
      </div>
      <div className='btn-submit' onClick={registerHandler}>
        <span>Register</span>
      </div>
    </div>
  );
};

export default Register;
