import axios from 'axios';

const api = axios.create({
  baseURL: "http://localhost:8080/task",
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

export const create = (data) => api.post('/create', data);
export const getTask = (id) => api.get(`/get/${id}`);
export const list = () => api.get('/list');
export const update = (id, data) => api.put(`/update/${id}`, data);
export const remove = (id) => api.delete(`/delete/${id}`);

export default api;
