import { useNavigate } from 'react-router-dom';

function Home({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    fetch("http://localhost:8080/logout", {
      method: "POST",
      credentials: "include",
    })
      .then(() => {
        onLogout();
        navigate("/login");
      })
      .catch((err) => {
        console.error("Logout failed:", err);
        onLogout();
        navigate("/login");
      });
  };

  if (!user) return <div>Redirecting...</div>;

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
        <p>Thank you for logging in to our platform. This system helps you manage and participate in various events.</p>
        <p>Stay tuned for upcoming features including:</p>
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