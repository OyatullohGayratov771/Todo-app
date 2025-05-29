import { Box, Typography, TextField, Button, Stack } from "@mui/material";
import { useState } from "react";
import { updateName, updateEmail, updatePassword } from "../services/authService";

function Profile() {
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: ""
  });

  const [messages, setMessages] = useState({
    name: "",
    email: "",
    password: ""
  });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
    setMessages({ ...messages, [e.target.name]: "" }); // eski xabarni tozalash
  };

  const handleUpdateField = async (field) => {
    try {
      if (!form[field]) {
        setMessages((prev) => ({
          ...prev,
          [field]: "Maydon bo‘sh bo‘lishi mumkin emas."
        }));
        return;
      }

      if (field === "name") await updateName({ name: form.name });
      if (field === "email") await updateEmail({ email: form.email });
      if (field === "password") await updatePassword({ password: form.password });

      setMessages((prev) => ({
        ...prev,
        [field]: "Muvaffaqiyatli yangilandi."
      }));

      setForm((prev) => ({ ...prev, [field]: "" })); // inputni tozalash
    } catch (err) {
      console.error(err);
      setMessages((prev) => ({
        ...prev,
        [field]: "Xatolik yuz berdi."
      }));
    }
  };

  return (
    <Box sx={{ p: 3, maxWidth: 400, mx: "auto" }}>
      <Typography variant="h5" gutterBottom>Profilni yangilash</Typography>

      {/* Name update */}
      <Stack direction="row" spacing={1} alignItems="center">
        <TextField
          label="Ism"
          name="name"
          fullWidth
          value={form.name}
          onChange={handleChange}
          InputProps={{ sx: { backgroundColor: "white" } }}
        />
        <Button variant="contained" onClick={() => handleUpdateField("name")}>OK</Button>
      </Stack>
      {messages.name && (
        <Typography variant="body2" sx={{ mt: 0.5 }} color="success.main">{messages.name}</Typography>
      )}

      {/* Email update */}
      <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 2 }}>
        <TextField
          label="Email"
          name="email"
          fullWidth
          value={form.email}
          onChange={handleChange}
          InputProps={{ sx: { backgroundColor: "white" } }}
        />
        <Button variant="contained" onClick={() => handleUpdateField("email")}>OK</Button>
      </Stack>
      {messages.email && (
        <Typography variant="body2" sx={{ mt: 0.5 }} color="success.main">{messages.email}</Typography>
      )}

      {/* Password update */}
      <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 2 }}>
        <TextField
          label="Yangi parol"
          name="password"
          type="password"
          fullWidth
          value={form.password}
          onChange={handleChange}
          InputProps={{ sx: { backgroundColor: "white" } }}
        />
        <Button variant="contained" onClick={() => handleUpdateField("password")}>OK</Button>
      </Stack>
      {messages.password && (
        <Typography variant="body2" sx={{ mt: 0.5 }} color="success.main">{messages.password}</Typography>
      )}

      {/* Back button */}
      <Button
        variant="outlined"
        fullWidth
        sx={{ mt: 4 }}
        onClick={() => window.history.back()}
      >
        Back
      </Button>
    </Box>
  );
}

export default Profile;
