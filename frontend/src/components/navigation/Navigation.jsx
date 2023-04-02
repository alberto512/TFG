import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';

const Navigation = () => {
  const { token, onLogout } = useAuth();

  return (
    <nav>
      <NavLink to='/'>Home</NavLink>
      {!token && (
        <>
          <NavLink to='/login'>Login</NavLink>
          <NavLink to='/register'>Register</NavLink>
        </>
      )}
      {token && (
        <>
          <NavLink to='/dashboard'>Dashboard</NavLink>
          <button type='button' onClick={onLogout}>
            Sign Out
          </button>
        </>
      )}
    </nav>
  );
};

export default Navigation;
