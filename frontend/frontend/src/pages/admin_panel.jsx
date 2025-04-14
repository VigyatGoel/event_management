import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './admin_panel.css';

function AdminPanel({ user, onLogout }) {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('dashboard');

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

  // Mock data for demonstration
  const stats = {
    totalEvents: 24,
    totalUsers: 156,
    totalRegistrations: 342,
  };

  const recentUsers = [
    { id: 1, name: "John Doe", email: "john@example.com", role: "attendee", joinDate: "2025-04-10" },
    { id: 2, name: "Jane Smith", email: "jane@example.com", role: "organiser", joinDate: "2025-04-09" },
    { id: 3, name: "Mike Johnson", email: "mike@example.com", role: "attendee", joinDate: "2025-04-08" },
  ];

  if (!user || user.role !== 'admin') {
    navigate('/login');
    return <div>Redirecting...</div>;
  }

  return (
    <div className="admin-panel">
      <div className="admin-header">
        <div className="admin-title">
          <h1>Event Management System</h1>
          <span className="admin-badge">Admin Panel</span>
        </div>
        <div className="admin-controls">
          <span>Welcome, {user.name}</span>
          <button className="logout-button" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>

      {/* Single admin panel content with one sidebar */}
      <div className="admin-panel-content">
        <div className="admin-sidebar">
          <ul className="sidebar-menu">
            <li>
              <a 
                href="#" 
                className={activeTab === 'dashboard' ? 'active' : ''} 
                onClick={(e) => {e.preventDefault(); setActiveTab('dashboard')}}
              >
                Dashboard
              </a>
            </li>
            <li>
              <a 
                href="#" 
                className={activeTab === 'users' ? 'active' : ''} 
                onClick={(e) => {e.preventDefault(); setActiveTab('users')}}
              >
                User Management
              </a>
            </li>
            <li>
              <a 
                href="#" 
                className={activeTab === 'events' ? 'active' : ''} 
                onClick={(e) => {e.preventDefault(); setActiveTab('events')}}
              >
                Event Management
              </a>
            </li>
            <li>
              <a 
                href="#" 
                className={activeTab === 'registrations' ? 'active' : ''} 
                onClick={(e) => {e.preventDefault(); setActiveTab('registrations')}}
              >
                Registrations
              </a>
            </li>
            <li>
              <a 
                href="#" 
                className={activeTab === 'settings' ? 'active' : ''} 
                onClick={(e) => {e.preventDefault(); setActiveTab('settings')}}
              >
                Settings
              </a>
            </li>
          </ul>
        </div>

        <div className="admin-main">
          {activeTab === 'dashboard' && (
            <>
              <h2>Dashboard</h2>
              
              <div className="dashboard-stats">
                <div className="stat-card">
                  <h3>Total Events</h3>
                  <div className="stat-value">{stats.totalEvents}</div>
                </div>
                <div className="stat-card">
                  <h3>Total Users</h3>
                  <div className="stat-value">{stats.totalUsers}</div>
                </div>
                <div className="stat-card">
                  <h3>Total Registrations</h3>
                  <div className="stat-value">{stats.totalRegistrations}</div>
                </div>
              </div>

              <h3>Recent Users</h3>
              <div className="search-bar">
                <input type="text" placeholder="Search users..." />
                <span className="search-icon">🔍</span>
              </div>
              
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Join Date</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {recentUsers.map(user => (
                    <tr key={user.id}>
                      <td>{user.name}</td>
                      <td>{user.email}</td>
                      <td>{user.role}</td>
                      <td>{user.joinDate}</td>
                      <td>
                        <div className="action-buttons">
                          <button className="btn btn-secondary">Edit</button>
                          <button className="btn btn-danger">Delete</button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}

          {activeTab === 'users' && (
            <div>
              <h2>User Management</h2>
              <p>User management interface will be implemented here.</p>
            </div>
          )}

          {activeTab === 'events' && (
            <div>
              <h2>Event Management</h2>
              <p>Event management interface will be implemented here.</p>
            </div>
          )}

          {activeTab === 'registrations' && (
            <div>
              <h2>Registration Management</h2>
              <p>Registration management interface will be implemented here.</p>
            </div>
          )}

          {activeTab === 'settings' && (
            <div>
              <h2>Settings</h2>
              <p>Settings interface will be implemented here.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default AdminPanel;