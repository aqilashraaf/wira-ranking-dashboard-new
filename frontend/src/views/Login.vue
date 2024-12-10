<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md">
      <h1 class="text-3xl font-bold text-center text-gray-900 dark:text-white">WIRA Dashboard Login</h1>
      <div v-if="!requires2FA" class="mt-8 space-y-6">
        <div class="space-y-2">
          <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Enter username"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
            @keyup.enter="handleLogin"
          />
        </div>
        <div class="space-y-2">
          <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
          <div class="relative">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter password"
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
              @keyup.enter="handleLogin"
            />
            <button 
              type="button" 
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-600 dark:text-gray-400"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '🔒' : '👁️' }}
            </button>
          </div>
        </div>
        <button 
          @click="handleLogin" 
          :disabled="!isValid"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Login
        </button>
        <p class="text-center text-sm text-gray-600 dark:text-gray-400">
          Don't have an account? 
          <router-link to="/register" class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300">
            Register
          </router-link>
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
        
        // Check if 2FA is required
        if (response.data.requires_2fa) {
          console.log('2FA required, showing prompt...');
          requires2FA.value = true;
          loginCredentials.value = {
            username: username.value,
            password: password.value
          };
          toast.info('Please enter your 2FA code');
          return;
        }
        
        // If no 2FA required, proceed with login
        localStorage.setItem('access_token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
        axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.access_token}`;
        toast.success('Login successful');
        router.push(router.currentRoute.value.query.redirect || '/dashboard');
        
      } catch (error) {
        console.log('Login error response:', error.response?.data);
        toast.error(error.response?.data?.error || 'Login failed');
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
/* Removed styles */
</style>
