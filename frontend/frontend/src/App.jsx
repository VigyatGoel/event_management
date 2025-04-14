import { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import Login from './pages/login';
import Signup from './pages/signup';
import Home from './pages/home';
import AdminPanel from './pages/admin_panel';
import './App.css';
import './pages/home.css';

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  
  // Check if we're on the admin panel page
  const isAdminPanel = location.pathname === '/admin';

  useEffect(() => {
    const checkSession = async () => {
      try {
        const res = await fetch('http://localhost:8080/session', {
          method: 'GET',
          credentials: 'include',
        });

        if (res.ok) {
          const data = await res.json();
          
          // Make sure we have role information from the session
          setUser({
            email: data.email,
            name: data.name,
            role: data.role
          });
        }
      } catch (err) {
        console.error('Session check failed:', err);
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, []);

  const handleLogin = (userData) => {
    setUser(userData);
  };

  const handleLogout = () => {
    setUser(null);
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  // Return admin panel directly without wrapping it in app container when on admin route
  if (isAdminPanel && user && user.role === 'admin') {
    return <AdminPanel user={user} onLogout={handleLogout} />;
  }

  return (
    <div className="app-container">
      <h1>Event Management</h1>

      <Routes>
        <Route 
          path="/login" 
          element={user ? (
            <Navigate to={user.role === 'admin' ? '/admin' : '/'} replace />
          ) : (
            <div className="form-container">
              <Login onLoginSuccess={handleLogin} />
            </div>
          )} 
        />
        <Route 
          path="/signup" 
          element={user ? (
            <Navigate to={user.role === 'admin' ? '/admin' : '/'} replace />
          ) : (
            <div className="form-container">
              <Signup onSignupSuccess={() => navigate('/login')} />
            </div>
          )} 
        />
        <Route 
          path="/admin" 
          element={
            user && user.role === 'admin' ? (
              <AdminPanel user={user} onLogout={handleLogout} />
            ) : (
              <Navigate to="/login" replace />
            )
          } 
        />
        <Route 
          path="/" 
          element={user ? (
            user.role === 'admin' ? (
              <Navigate to="/admin" replace />
            ) : (
              <Home user={user} onLogout={handleLogout} />
            )
          ) : (
            <Navigate to="/login" replace />
          )} 
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </div>
  );
}

export default App;