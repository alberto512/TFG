import { Routes, Route } from 'react-router-dom';
import AuthProvider, {
  ProtectedRoute,
} from 'components/authProvider/AuthProvider';
import Navigation from 'components/navigation/Navigation';
import Home from 'pages/home/Home';
import Login from 'pages/login/Login';
import Register from 'pages/register/Register';
import Dashboard from 'pages/dashboard/Dashboard';
import RegisterBank from 'pages/registerBank/RegisterBank';
import SantanderLogin from 'pages/santander/SantanderLogin';
import './App.css';

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
          path='dashboard'
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
      </Routes>
    </AuthProvider>
  );
}

export default App;
