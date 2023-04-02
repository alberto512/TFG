import React, { useState, useEffect, useContext, createContext } from 'react';
import { Navigate, useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';

const AuthContext = createContext(null);
const path = process.env.REACT_APP_BACKEND_URL;

export const useAuth = () => {
  return useContext(AuthContext);
};

export const ProtectedRoute = ({ children }) => {
  const { token } = useAuth();
  const location = useLocation();

  if (!token) {
    return <Navigate to='/' replace state={{ from: location }} />;
  }

  return children;
};

function getJwt() {
  const saved = localStorage.getItem('jwt');
  const initial = JSON.parse(saved);
  return initial || null;
}

const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(() => {
    return getJwt();
  });
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    localStorage.setItem('jwt', JSON.stringify(token));
  }, [token]);

  const login = async (username, password) => {
    try {
      const response = await axios.post(path, {
        query: `mutation Login($username: String!, $password: String!) {
          login(username: $username, password: $password)
        }`,
        variables: {
          username,
          password,
        },
      });

      if (response.data.data.login === '') {
        return;
      }

      setToken(response.data.data.login);

      const origin = location.state?.from?.pathname || '/dashboard';
      navigate(origin);
    } catch (error) {
      console.log(error);
    }
  };

  const register = async (username, password, role) => {
    try {
      const response = await axios.post(path, {
        query: `mutation Login($username: String!, $password: String!) {
          login(username: $username, password: $password)
        }`,
        variables: {
          username,
          password,
        },
      });

      if (response.data.data.login === '') {
        return;
      }

      setToken(response.data.data.login);

      const origin = location.state?.from?.pathname || '/dashboard';
      navigate(origin);
    } catch (error) {
      console.log(error);
    }
  };

  const logout = () => {
    setToken(null);
  };

  const value = {
    token,
    onLogin: login,
    onRegister: register,
    onLogout: logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export default AuthProvider;
