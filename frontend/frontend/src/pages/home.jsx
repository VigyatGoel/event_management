import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function Homepage({ user, onLogout }) {
  const navigate = useNavigate();

  // Check if the user is logged in, and redirect accordingly
  useEffect(() => {
    const checkSession = async () => {
      try {
        const response = await fetch('http://localhost:8080/redirect-home', {
          method: 'GET',
          credentials: 'include', // Important to send cookies with the request
        });

        if (response.ok) {
          // Session is valid, stay on the homepage
          return;
        } else {
          // Redirect to login if session is not valid
          navigate("/login");
        }
      } catch (error) {
        // Redirect to login on error
        navigate("/login");
      }
    };

    checkSession(); // Call the check session on page load
  }, [navigate]);

  const handleLogout = () => {
    onLogout(); // Log the user out and possibly redirect in the parent component
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
