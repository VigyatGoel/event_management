import { useState } from "react";

function Signup() {
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [msg, setMsg] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  const handleChange = e => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async e => {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/signup/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: new URLSearchParams(form),
    });

    if (res.ok) {
      setMsg("Signup successful! Go to login.");
      setIsSuccess(true);
    } else {
      setMsg("Signup failed: Email might already exist.");
      setIsSuccess(false);
    }
  };

  return (
    <div>
      <h2>Create Account</h2>
      <p>Register to manage and join events</p>
      <form onSubmit={handleSubmit}>
        <input 
          name="name" 
          placeholder="Full Name" 
          className="form-input"
          onChange={handleChange} 
          required 
        />
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
        <button type="submit" className="submit-button">Create Account</button>
      </form>
      {msg && (
        <div className={`message ${isSuccess ? 'success-message' : 'error-message'}`}>
          {msg}
        </div>
      )}
    </div>
  );
}

export default Signup;