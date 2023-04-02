import { Routes, Route } from 'react-router-dom';
import AuthProvider from 'components/authProvider/AuthProvider';
import Navigation from 'components/navigation/Navigation';
import Home from 'pages/home/Home';
import Login from 'pages/login/Login';
import './App.css';

function App() {
  return (
    <AuthProvider>
      <Navigation />
      <Routes>
        <Route path='/' element={<Home />} />
        <Route path='login' element={<Login />} />
      </Routes>
    </AuthProvider>
  );
}

export default App;
