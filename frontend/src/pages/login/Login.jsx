import React, { useState } from 'react';
import { useAuth } from 'components/authProvider/AuthProvider';
import './Login.css';

const Login = () => {
  const { onLogin } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const loginHandler = () => {
    onLogin(username, password);
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
      <div className='btn-submit' onClick={loginHandler}>
        <span>Sign In</span>
      </div>
    </div>
  );
};

export default Login;
