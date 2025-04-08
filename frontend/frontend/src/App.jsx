import { useState, useEffect } from 'react';
import Login from './pages/login';
import Signup from './pages/signup';
import Homepage from './pages/home';
import './App.css';
import './pages/home.css';

function App() {
  const [page, setPage] = useState('login');
  const [user, setUser] = useState(null);
  
  // Check if user is logged in (from localStorage)
  useEffect(() => {
    const loggedInUser = localStorage.getItem('user');
    if (loggedInUser) {
      setUser(JSON.parse(loggedInUser));
      setPage('homepage');
    }
  }, []);

  const handleLogin = (userData) => {
    setUser(userData);
    // Save user data to localStorage
    localStorage.setItem('user', JSON.stringify(userData));
    setPage('homepage');
  };
  
  const handleLogout = () => {
    setUser(null);
    // Remove user data from localStorage
    localStorage.removeItem('user');
    setPage('login');
  };

  // If user is logged in, show homepage
  if (user) {
    return <Homepage user={user} onLogout={handleLogout} />;
  }

  // Otherwise show login/signup
  return (
    <div className="app-container">
      <h1>Event Management</h1>
      <div className="tab-buttons">
        <button 
          className={`tab-button ${page === 'login' ? 'active' : ''}`} 
          onClick={() => setPage('login')}
        >
          Login
        </button>
        <button 
          className={`tab-button ${page === 'signup' ? 'active' : ''}`} 
          onClick={() => setPage('signup')}
        >
          Signup
        </button>
      </div>
      <div className="form-container">
        {page === 'login' ? (
          <Login onLoginSuccess={handleLogin} />
        ) : (
          <Signup onSignupSuccess={() => setPage('login')} />
        )}
      </div>
    </div>
  );
}

export default App;