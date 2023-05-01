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
import {
  faBars,
  faXmark,
  faSpinner,
  faPenToSquare,
} from '@fortawesome/free-solid-svg-icons';
import './App.css';
import Account from 'pages/account/Account';
import EditCategory from 'pages/editCategory/EditCategory';

library.add(faBars, faXmark, faSpinner, faPenToSquare);

function App() {
  return (
    <AuthProvider>
      <Navigation />
      <Routes>
        <Route path='/' element={<Home />} />
        <Route path='login' element={<Login />} />
        <Route path='register' element={<Register />} />
        <Route path='registerBank' element={<RegisterBank />} />
        <Route path='santanderLogin' element={<SantanderLogin />} />
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
        <Route
          path='editCategory/:id'
          element={
            <ProtectedRoute>
              <EditCategory />
            </ProtectedRoute>
          }
        />
      </Routes>
    </AuthProvider>
  );
}

export default App;
