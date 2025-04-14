import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Home({ user, setUser, onLogout }) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkSession = async () => {
      try {
        const response = await fetch('http://localhost:8080/session', {
          method: 'GET',
          credentials: 'include', // Important to send cookies with the request
        });

        if (response.ok) {
          // Get user info from the session response
          const userData = await response.json();
          setUser(userData); // Set user from session
        } else {
          // Redirect to login if session is invalid
          navigate("/login");
        }
      } catch (error) {
        console.error("Session check failed:", error);
        navigate("/login");
      } finally {
        setLoading(false);
      }
    };

    checkSession();
  }, [navigate, setUser]);

  const handleLogout = () => {
    // Send logout request to backend to clear the session
    fetch("http://localhost:8080/logout", {
      method: "POST",
      credentials: "include",
    })
      .then(() => {
        onLogout(); // Call logout handler passed as prop
        navigate("/login"); // Redirect to login page
      })
      .catch((err) => {
        console.error("Logout failed:", err);
        onLogout();
        navigate("/login");
      });
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!user) return null;

  return (
    <div className="homepage">
      <div className="homepage-header">
        <div className="system-title">
          <h1>Event Management System</h1>
        </div>
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

export default Home;