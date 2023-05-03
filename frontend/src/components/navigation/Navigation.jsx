import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from 'components/authProvider/AuthProvider';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './Navigation.css';

const Navigation = () => {
  const { token, onLogout } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);

  const handleLogout = () => {
    onLogout();
    setIsMenuOpen(false);
  };

  return (
    <>
      {!isMenuOpen && (
        <FontAwesomeIcon
          className='icon-menu'
          icon='fa-solid fa-bars'
          onClick={() => setIsMenuOpen(true)}
        />
      )}
      {isMenuOpen && (
        <div className='menu-wrapper' onMouseLeave={() => setIsMenuOpen(false)}>
          <FontAwesomeIcon
            className='close-icon'
            icon='fa-solid fa-xmark'
            onClick={() => setIsMenuOpen(false)}
          />
          <NavLink
            className='nav-link'
            to='/'
            onClick={() => setIsMenuOpen(false)}
          >
            Home
          </NavLink>
          {!token && (
            <>
              <NavLink
                className='nav-link'
                to='/login'
                onClick={() => setIsMenuOpen(false)}
              >
                Login
              </NavLink>
              <NavLink
                className='nav-link'
                to='/register'
                onClick={() => setIsMenuOpen(false)}
              >
                Register
              </NavLink>
            </>
          )}
          {token && (
            <>
              <NavLink
                className='nav-link'
                to='/accounts'
                onClick={() => setIsMenuOpen(false)}
              >
                Accounts
              </NavLink>
              <NavLink
                className='nav-link'
                to='/stats'
                onClick={() => setIsMenuOpen(false)}
              >
                Stats
              </NavLink>
              <span className='nav-link' onClick={handleLogout}>
                Sign Out
              </span>
            </>
          )}
        </div>
      )}
    </>
  );
};

export default Navigation;
