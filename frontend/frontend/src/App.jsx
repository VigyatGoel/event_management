import { useState, useEffect, lazy, Suspense } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import './App.css';

// Import only the components needed for initial load
import Login from './pages/login';

// Lazy load other components
const Signup = lazy(() => import('./pages/signup'));
const Home = lazy(() => import('./pages/home'));
const AdminPanel = lazy(() => import('./pages/admin_panel'));

// Loading component to show while lazy components load
const LoadingFallback = () => <div className="loading-container">Loading...</div>;

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

  const handleSignupSuccess = () => {
    console.log("Signup successful!");
  };

  if (loading) {
    return <div className="loading-container">Loading...</div>;
  }

  // Return admin panel directly without wrapping it in app container when on admin route
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