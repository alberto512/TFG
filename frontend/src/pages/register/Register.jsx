import React, { useState } from 'react';
import { useAuth } from 'components/authProvider/AuthProvider';
import './Register.css';

const Register = () => {
  const { token, onRegister } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('');

  const registerHandler = () => {
    console.log(username, password, role);
    //onRegister(username, password, role);
  };

  return (
    <>
      <h1>Register</h1>
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
          <select value={role} onChange={(e) => setRole(e.target.value)}>
            <option value='ADMIN'>Admin</option>
            <option value='USER'>User</option>
          </select>
          <button type='button' onClick={registerHandler}>
            Sign In
          </button>
        </>
      )}
    </>
  );
};

export default Login;
