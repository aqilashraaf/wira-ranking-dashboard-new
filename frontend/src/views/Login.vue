<template>
  <div class="login-container">
    <div class="login-card">
      <h1>WIRA Dashboard Login</h1>
      <div v-if="!requires2FA" class="login-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Enter username"
            @keyup.enter="handleLogin"
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <div class="password-input-container">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter password"
              @keyup.enter="handleLogin"
            />
            <button 
              type="button" 
              class="toggle-password"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '🔒' : '👁️' }}
            </button>
          </div>
        </div>
        <button @click="handleLogin" :disabled="!isValid">
          Login
        </button>
        <p class="register-link">
          Don't have an account? <router-link to="/register">Register</router-link>
        </p>
      </div>
      <div v-else>
        <TwoFactorVerify
          :onVerify="handle2FAVerification"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import axios from 'axios';
import TwoFactorVerify from '../components/TwoFactorVerify.vue';

export default {
  name: 'Login',
  components: {
    TwoFactorVerify
  },
  setup() {
    const router = useRouter();
    const toast = useToast();
    
    const username = ref('');
    const password = ref('');
    const requires2FA = ref(false);
    const loginCredentials = ref(null);
    const showPassword = ref(false);

    const isValid = computed(() => {
      return username.value.length > 0 && password.value.length > 0;
    });

    const handleLogin = async () => {
      try {
        console.log('Attempting login...');
        const response = await axios.post('/api/auth/login', {
          username: username.value,
          password: password.value,
          totp_code: null
        });
        console.log('Login response:', response.data);
        
        // Store tokens and redirect
        localStorage.setItem('access_token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
        axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.access_token}`;
        toast.success('Login successful');
        router.push(router.currentRoute.value.query.redirect || '/dashboard');
        
      } catch (error) {
        console.log('Login error response:', error.response?.data);
        
        // Check if 2FA is required
        if (error.response?.data?.requires_2fa) {
          console.log('2FA required, showing prompt...');
          requires2FA.value = true;
          loginCredentials.value = {
            username: username.value,
            password: password.value
          };
          toast.info('Please enter your 2FA code');
        } else {
          // Handle other errors
          toast.error(error.response?.data?.error || 'Login failed');
        }
      }
    };

    const handle2FAVerification = async (code) => {
      try {
        console.log('Verifying 2FA code...');
        const response = await axios.post('/api/auth/login', {
          username: loginCredentials.value.username,
          password: loginCredentials.value.password,
          totp_code: code
        });
        
        console.log('2FA verification response:', response.data);
        
        // Store tokens
        localStorage.setItem('access_token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
        axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.access_token}`;
        
        toast.success('Login successful');
        router.push('/dashboard');
      } catch (error) {
        console.error('2FA verification error:', error.response?.data || error);
        toast.error(error.response?.data?.error || 'Invalid 2FA code');
      }
    };

    return {
      username,
      password,
      requires2FA,
      isValid,
      handleLogin,
      handle2FAVerification,
      showPassword
    };
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f5f5;
}

.login-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

h1 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  color: #4a5568;
  font-size: 0.875rem;
}

input {
  width: 100%;
  padding: 10px;
  margin-top: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  transition: box-shadow 0.3s ease;
}

#password {
  letter-spacing: 1px;
}

.password-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.toggle-password {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  font-size: 1.2rem;
}

button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background: #45a049;
}

button:disabled {
  background: #cccccc;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 1rem;
  color: #4a5568;
}

.register-link a {
  color: #4CAF50;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
