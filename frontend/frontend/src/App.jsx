import { useState, useEffect, lazy, Suspense } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import './App.css';

import Login from './pages/login';

const Signup = lazy(() => import('./pages/signup'));
const Home = lazy(() => import('./pages/home'));
const AdminPanel = lazy(() => import('./pages/admin_panel'));

const LoadingFallback = () => <div className="loading-container">Loading...</div>;

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  
  const isAdminPanel = location.pathname === '/admin';

  useEffect(() => {
    const checkSession = async () => {
      try {
        // Check if we have a token in localStorage
        const token = localStorage.getItem('token');
        
        if (!token) {
          setLoading(false);
          return;
        }
        
        const res = await fetch('http://localhost:8080/validate_token', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
        });

        if (res.ok) {
          const data = await res.json();
          
          // Ensure we include the token in the user object
          setUser({
            email: data.email,
            name: data.name,
            role: data.role,
            token: token
          });
        } else {
          console.error('Token validation failed with status:', res.status);
          localStorage.removeItem('token');
        }
      } catch (err) {
        console.error('Token validation failed:', err);
        localStorage.removeItem('token');
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, []);

  const handleLogin = (userData) => {
    // Ensure we're saving the token to localStorage here too
    if (userData.token) {
      localStorage.setItem('token', userData.token);
    }
    setUser(userData);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  const handleSignupSuccess = () => {
    console.log("Signup successful!");
  };

  if (loading) {
    return <div className="loading-container">Loading...</div>;
  }

  if (isAdminPanel && user && user.role === 'admin') {
    return (
      <Suspense fallback={<LoadingFallback />}>
        <AdminPanel user={user} onLogout={handleLogout} />
      </Suspense>
    );
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
              <Suspense fallback={<LoadingFallback />}>
                <Signup onSignupSuccess={handleSignupSuccess} />
              </Suspense>
            </div>
          )} 
        />
        <Route 
          path="/admin" 
          element={
            user && user.role === 'admin' ? (
              <Suspense fallback={<LoadingFallback />}>
                <AdminPanel user={user} onLogout={handleLogout} />
              </Suspense>
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
              <Suspense fallback={<LoadingFallback />}>
                <Home user={user} onLogout={handleLogout} />
              </Suspense>
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