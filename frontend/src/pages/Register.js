import { useState } from "react";
import { register } from "../services/authService";
import { useNavigate } from "react-router-dom";

function Register() {
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    setErrors({});

    try {
      const res = await register({ username, email, password });
      localStorage.setItem("token", res.data.token);
      navigate("/dashboard");
    } catch (err) {
      const resError = err.response?.data;
      if (resError?.errors) {
        setErrors(resError.errors);
      } else if (resError?.error) {
        setErrors({ general: resError.error });
      } else {
        setErrors({ general: "Noma'lum xatolik yuz berdi" });
      }
    }
  };

  return (
    <div className="container">
      <h2>Ro‘yxatdan o‘tish</h2>
      <form onSubmit={handleRegister}>
        <input
          type="text"
          placeholder="Foydalanuvchi nomi"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          className="input"
        />
        {errors.username && <p className="error">{errors.username}</p>}

        <input
          type="email"
          placeholder="Email manzili"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="input"
        />
        {errors.email && <p className="error">{errors.email}</p>}

        <input
          type="password"
          placeholder="Parol"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="input"
        />
        {errors.password && <p className="error">{errors.password}</p>}

        {errors.general && <p className="error">{errors.general}</p>}

        <button type="submit">Ro‘yxatdan o‘tish</button>
      </form>
      <p>
        Allaqachon hisobingiz bormi? <a href="/">Kirish</a>
      </p>
    </div>
  );
}

export default Register;
