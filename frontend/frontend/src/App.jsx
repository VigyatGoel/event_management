import { useState, useEffect } from 'react';
import Login from './pages/login';
import Signup from './pages/signup';
import Homepage from './pages/home';
import './App.css';
import './pages/home.css';

function App() {
  const [page, setPage] = useState('login');
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true); // show loading while checking session

  useEffect(() => {
    const checkSession = async () => {
      try {
        const res = await fetch('http://localhost:8080/session', {
          method: 'GET',
          credentials: 'include',
        });

        if (res.ok) {
          const data = await res.json();
          setUser(data);
          setPage('homepage');
        }
      } catch (err) {
        console.error('Session check failed:', err);
      } finally {
        setLoading(false); // either way, done loading
      }
    };

    checkSession();
  }, []);

  const handleLogin = (userData) => {
    setUser(userData);
    setPage('homepage');
  };

  const handleLogout = () => {
    setUser(null);
    setPage('login');
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="app-container">
      <h1>Event Management</h1>

      {user ? (
        <Homepage user={user} onLogout={handleLogout} />
      ) : (
        <>
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
        </>
      )}
    </div>
  );
}

export default App;