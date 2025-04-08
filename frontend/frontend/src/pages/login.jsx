import { useState } from "react";

function Login() {
  const [form, setForm] = useState({ email: "", password: "" });
  const [msg, setMsg] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  const handleChange = e => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async e => {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/login/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: new URLSearchParams(form),
    });

    if (res.ok) {
      const text = await res.text();
      setMsg(text);
      setIsSuccess(true);
    } else {
      setMsg("Login failed. Please check your credentials.");
      setIsSuccess(false);
    }
  };

  return (
    <div>
      <h2>Welcome Back</h2>
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
        <button type="submit" className="submit-button">Login</button>
      </form>
      {msg && (
        <div className={`message ${isSuccess ? 'success-message' : 'error-message'}`}>
          {msg}
        </div>
      )}
    </div>
  );
}

export default Login;