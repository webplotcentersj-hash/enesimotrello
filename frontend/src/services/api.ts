import axios from 'axios';
import { ensureAnonymousUserId } from './anonymous.ts';
import authService from './auth.ts';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor: Add authentication token and anonymous user ID
api.interceptors.request.use(
  (config) => {
    // Add JWT token if authenticated
    const token = authService.getToken();
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    } else {
      // Fallback to anonymous user ID if not authenticated
      const anonymousUserId = ensureAnonymousUserId();
      if (anonymousUserId) {
        config.headers['X-Anonymous-User-Id'] = anonymousUserId;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor: Handle authentication errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, logout user
      authService.logout();
      // Optionally redirect to login
      if (window.location.pathname !== '/') {
        window.location.href = '/';
      }
    }
    return Promise.reject(error);
  }
);

// Board API
export const boardAPI = {
  getBoards: () => api.get('/boards'),
  createBoard: (data: { title: string; description?: string }) =>
    api.post('/boards', data),
  getBoard: (id: number) => api.get(`/boards/${id}`),
  updateBoard: (id: number, data: { title: string; description?: string }) =>
    api.put(`/boards/${id}`, data),
  deleteBoard: (id: number) => api.delete(`/boards/${id}`),
};

// Task API
export const taskAPI = {
  getTasks: (boardId: number) => api.get(`/tasks/board/${boardId}`),
  createTask: (boardId: number, data: {
    title: string;
    description?: string;
    priority?: 'low' | 'medium' | 'high';
    assignee_id?: number;
    due_date?: string;
  }) => api.post(`/tasks/board/${boardId}`, data),
  getTask: (boardId: number, taskId: number) => api.get(`/tasks/${taskId}`),
  updateTask: (boardId: number, taskId: number, data: {
    title: string;
    description?: string;
    status?: 'todo' | 'in_progress' | 'done';
    priority?: 'low' | 'medium' | 'high';
    assignee_id?: number;
    due_date?: string;
  }) => api.put(`/tasks/${taskId}`, data),
  deleteTask: (boardId: number, taskId: number) => api.delete(`/tasks/${taskId}`),
};

export default api;

