import { useState } from 'react';

function Homepage({ user, onLogout }) {
  return (
    <div className="homepage">
      <div className="homepage-header">
        <h1>Event Management System</h1>
        <div className="user-info">
          <span>Welcome, {user.name}</span>
          <button className="logout-button" onClick={onLogout}>
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