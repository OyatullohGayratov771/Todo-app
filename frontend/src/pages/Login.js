import { useState } from "react";
import { login } from "../services/authService";
import { useNavigate } from "react-router-dom";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setErrors({});

    try {
      const res = await login({ username, password });
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
      <h2>Kirish</h2>
      <form onSubmit={handleLogin}>
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
          type="password"
          placeholder="Parol"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="input"
        />
        {errors.password && <p className="error">{errors.password}</p>}

        {errors.general && <p className="error">{errors.general}</p>}

        <button type="submit">Kirish</button>
      </form>
      <p>Hisobingiz yo‘qmi? <a href="/register">Ro‘yxatdan o‘tish</a></p>
    </div>
  );
}

export default Login;
