import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function Homepage({ user, onLogout }) {
  const navigate = useNavigate();

  // Check JWT token on component mount
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/login");
    }
  }, [navigate]);

  // Clear token and logout
  const handleLogout = () => {
    localStorage.removeItem("token");
    onLogout(); // optional: can redirect or reset app state
  };

  return (
    <div className="homepage">
      <div className="homepage-header">
        <h1>Event Management System</h1>
        <div className="user-info">
          <span>Welcome, {user.name}</span>
          <button className="logout-button" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>

      <div className="welcome-container">
        <h2>Welcome to the Event Management System</h2>
        <p>
          Thank you for logging in to our platform. This system helps you manage and 
          participate in various events.
        </p>
        <p>
          Stay tuned for upcoming features including:
        </p>
        <ul>
          <li>Creating new events</li>
          <li>Registering for events</li>
          <li>Managing your event calendar</li>
          <li>Connecting with other participants</li>
        </ul>
      </div>
    </div>
  );
}

export default Homepage;