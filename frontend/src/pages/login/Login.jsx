import React, { useState } from 'react';
import './Login.css';
import { useAuth } from 'components/authProvider/AuthProvider';

const Login = () => {
  const { token, onLogin } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const loginHandler = () => {
    onLogin(username, password);
  };

  return (
    <>
      <h1>Login</h1>
      {!token && (
        <>
          <input
            type='text'
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type='password'
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button type='button' onClick={loginHandler}>
            Sign In
          </button>
        </>
      )}
    </>
  );
};

export default Login;
