// src/services/authService.js
import axios from 'axios';

const api = axios.create({
  baseURL: "http://localhost:8080/user",
});

export const login = (data) => api.post('/login', data);
export const register = (data) => api.post('/register', data);
export const updateName = (data) => api.put('/updatename', data);
export const updatePassword = (data) => api.put('/updatepassword', data);
export const updateEmail = (data) => api.put('/updateemail', data);

export default api;