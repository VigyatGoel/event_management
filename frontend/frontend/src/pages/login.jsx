import { useState } from "react";

function Login({ onLoginSuccess }) {
  const [form, setForm] = useState({ email: "", password: "" });
  const [msg, setMsg] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
  
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: new URLSearchParams(form),
      });
  
      let data = {};
      
      data = await res.json();
      
      if (res.ok) {
        setMsg(data.message);
        setIsSuccess(true);
  
        localStorage.setItem("token", data.token);
  
        setTimeout(() => {
          onLoginSuccess({
            email: data.email,
            name: data.name,
            token: data.token,
          });
        }, 1000);
      } else {
        setMsg(data.message || "Login failed. Please check your credentials.");
        setIsSuccess(false);
      }
    } catch (error) {
      console.error("Login error:", error);
      setMsg(error.message || "Connection error. Please try again.");
      setIsSuccess(false);
    }
  };

  return (
    <div>
      <h2>Welcome</h2>
      <p>Enter your credentials to access your account</p>
      <form onSubmit={handleSubmit}>
        <input
          name="email"
          type="email"
          placeholder="Email"
          className="form-input"
          onChange={handleChange}
          required
        />
        <input
          name="password"
          type="password"
          placeholder="Password"
          className="form-input"
          onChange={handleChange}
          required
        />
        <button type="submit" className="submit-button">
          Login
        </button>
      </form>
      {msg && (
        <div className={`message ${isSuccess ? "success-message" : "error-message"}`}>
          {msg}
        </div>
      )}
    </div>
  );
}

export default Login;