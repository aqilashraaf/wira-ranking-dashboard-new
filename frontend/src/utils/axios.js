import axios from 'axios';
import router from '../router';

// Configure axios defaults
const isDevelopment = import.meta.env.DEV;
const apiURL = isDevelopment
  ? 'http://localhost:8081'
  : 'https://ricrym.aqash.xyz/api';

console.log('Environment:', isDevelopment ? 'Development' : 'Production');
console.log('API URL:', apiURL);

axios.defaults.baseURL = apiURL;
axios.defaults.headers.common['Content-Type'] = 'application/json';

// Add request interceptor to include token
axios.interceptors.request.use(
  (config) => {
    // Log the request details for debugging
    const fullUrl = config.baseURL + config.url;
    console.log('Request:', {
      method: config.method?.toUpperCase(),
      url: config.url,
      fullUrl,
      data: config.data,
      headers: config.headers
    });
    
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
    console.error('Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
axios.interceptors.response.use(
  (response) => {
    console.log('Response:', {
      status: response.status,
      data: response.data
    });
    return response;
  },
  async (error) => {
    console.error('API Error:', {
      message: error.message,
      response: error.response ? {
        data: error.response.data,
        status: error.response.status,
        headers: error.response.headers
      } : null
    });

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

    const response = await axios.post('/auth/refresh', {
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

export default axios;
