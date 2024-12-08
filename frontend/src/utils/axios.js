import axios from 'axios';
import router from '../router';

// Configure axios defaults
const baseURL = import.meta.env.PROD 
  ? '/api'  // Production URL (relative path)
  : 'http://localhost:3000';  // Development URL

axios.defaults.baseURL = baseURL;

// Add request interceptor to include token
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      // Check token expiration before adding it
      try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const payload = JSON.parse(window.atob(base64));
        
        if (payload.exp && payload.exp * 1000 < Date.now()) {
          // Token is expired, attempt refresh
          const refreshToken = localStorage.getItem('refresh_token');
          if (refreshToken) {
            return refreshAccessToken()
              .then(newToken => {
                config.headers.Authorization = `Bearer ${newToken}`;
                return config;
              })
              .catch(() => {
                // If refresh fails, redirect to login
                handleAuthError();
                return Promise.reject('Token expired');
              });
          } else {
            handleAuthError();
            return Promise.reject('Token expired');
          }
        }
      } catch (e) {
        // If token is invalid, clear it
        handleAuthError();
        return Promise.reject('Invalid token');
      }
      
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle token refresh
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If error is 401 and we haven't tried to refresh token yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const newToken = await refreshAccessToken();
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return axios(originalRequest);
      } catch (refreshError) {
        handleAuthError();
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

async function refreshAccessToken() {
  try {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await axios.post('/api/auth/refresh', {
      refresh_token: refreshToken
    });

    const { access_token } = response.data;
    localStorage.setItem('access_token', access_token);
    return access_token;
  } catch (error) {
    throw error;
  }
}

function handleAuthError() {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  
  // Only redirect to login if we're not already there
  if (router.currentRoute.value.name !== 'Login') {
    router.push({
      name: 'Login',
      query: { redirect: router.currentRoute.value.fullPath }
    });
  }
}
