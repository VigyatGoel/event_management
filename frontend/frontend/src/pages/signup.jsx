import { useState } from "react";

function Signup() {
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [msg, setMsg] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const res = await fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: new URLSearchParams(form),
      });

      const contentType = res.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        throw new Error("Invalid server response format.");
      }

      const data = await res.json();

      if (res.ok) {
        setMsg(data.message || "Signup successful!");
        setIsSuccess(true);
      } else {
        setMsg(data.message || "Signup failed. Please try again.");
        setIsSuccess(false);
      }
    } catch (error) {
      console.error("Signup error:", error);
      setMsg("Server connection failed. Please try again later.");
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
        <button type="submit" className="submit-button">
          Create Account
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

export default Signup;