import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';
const TOKEN_KEY = 'taskboard_auth_token';
const USER_KEY = 'taskboard_user';

export interface User {
  id: number;
  email: string;
  username: string;
  first_name?: string;
  last_name?: string;
}

export interface LoginResponse {
  message: string;
  token: string;
  user: User;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  first_name?: string;
  last_name?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

class AuthService {
  // Token management
  getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
  }

  setToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(TOKEN_KEY, token);
  }

  removeToken(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  // User management
  getUser(): User | null {
    if (typeof window === 'undefined') return null;
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }

  setUser(user: User): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.getToken() !== null;
  }

  // Register new user
  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await axios.post<{ message: string; user: User }>(
      `${API_BASE_URL}/auth/register`,
      data
    );
    
    // After registration, automatically login
    const loginResponse = await this.login({
      email: data.email,
      password: data.password,
    });

    return loginResponse;
  }

  // Login user
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await axios.post<LoginResponse>(
      `${API_BASE_URL}/auth/login`,
      data
    );

    // Store token and user
    this.setToken(response.data.token);
    this.setUser(response.data.user);

    return response.data;
  }

  // Logout user
  logout(): void {
    this.removeToken();
    // Optionally redirect to login page
    window.location.href = '/';
  }

  // Get current user profile
  async getProfile(): Promise<User> {
    const response = await axios.get<User>(`${API_BASE_URL}/auth/profile`);
    this.setUser(response.data);
    return response.data;
  }

  // Update user profile
  async updateProfile(data: { first_name?: string; last_name?: string }): Promise<User> {
    const response = await axios.put<{ message: string; user: User }>(
      `${API_BASE_URL}/auth/profile`,
      data
    );
    this.setUser(response.data.user);
    return response.data.user;
  }
}

export const authService = new AuthService();
export default authService;

