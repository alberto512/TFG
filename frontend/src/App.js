import { Routes, Route } from 'react-router-dom';
import AuthProvider, {
  ProtectedRoute,
} from 'components/authProvider/AuthProvider';
import Navigation from 'components/navigation/Navigation';
import Home from 'pages/home/Home';
import Login from 'pages/login/Login';
import Register from 'pages/register/Register';
import Accounts from 'pages/accounts/Accounts';
import RegisterBank from 'pages/registerBank/RegisterBank';
import SantanderLogin from 'pages/santander/SantanderLogin';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faBars, faXmark, faSpinner } from '@fortawesome/free-solid-svg-icons';
import './App.css';
import Account from 'pages/account/Account';

library.add(faBars, faXmark, faSpinner);

function App() {
  return (
    <AuthProvider>
      <Navigation />
      <Routes>
        <Route path='/' element={<Home />} />
        <Route path='login' element={<Login />} />
        <Route path='register' element={<Register />} />
        <Route
          path='registerBank'
          element={
            <ProtectedRoute>
              <RegisterBank />
            </ProtectedRoute>
          }
        />
        <Route
          path='santanderLogin'
          element={
            <ProtectedRoute>
              <SantanderLogin />
            </ProtectedRoute>
          }
        />
        <Route
          path='accounts'
          element={
            <ProtectedRoute>
              <Accounts />
            </ProtectedRoute>
          }
        />
        <Route
          path='account/:id'
          element={
            <ProtectedRoute>
              <Account />
            </ProtectedRoute>
          }
        />
      </Routes>
    </AuthProvider>
  );
}

export default App;
