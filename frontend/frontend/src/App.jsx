import { useState } from 'react';
import Login from './pages/login';
import Signup from './pages/signup';
import './App.css';

function App() {
  const [page, setPage] = useState('login');

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
        {page === 'login' ? <Login /> : <Signup />}
      </div>
    </div>
  );
}

export default App;