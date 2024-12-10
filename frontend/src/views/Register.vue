<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md">
      <h1 class="text-3xl font-bold text-center text-gray-900 dark:text-white">WIRA Dashboard Register</h1>
      <div class="mt-8 space-y-6">
        <div class="space-y-2">
          <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Choose a username"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
            @blur="validateUsername"
          />
          <span class="form-help error" v-if="errors.username">{{ errors.username }}</span>
        </div>

        <div class="space-y-2">
          <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="Enter your email"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
            @blur="validateEmail"
          />
          <span class="form-help error" v-if="errors.email">{{ errors.email }}</span>
        </div>

        <div class="space-y-2">
          <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
          <div class="relative">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Choose a password"
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
              @blur="validatePassword"
            />
            <button 
              type="button" 
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-600 dark:text-gray-400"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '🔒' : '👁️' }}
            </button>
          </div>
          <div v-if="password" class="password-requirements">
            <div class="requirement" :class="{ met: hasMinLength }">
              <span class="icon">{{ hasMinLength ? '✅' : '❌' }}</span>
              At least 8 characters
            </div>
            <div class="requirement" :class="{ met: hasUpperCase }">
              <span class="icon">{{ hasUpperCase ? '✅' : '❌' }}</span>
              At least one uppercase letter
            </div>
            <div class="requirement" :class="{ met: hasLowerCase }">
              <span class="icon">{{ hasLowerCase ? '✅' : '❌' }}</span>
              At least one lowercase letter
            </div>
            <div class="requirement" :class="{ met: hasNumber }">
              <span class="icon">{{ hasNumber ? '✅' : '❌' }}</span>
              At least one number
            </div>
            <div class="requirement" :class="{ met: hasSpecialChar }">
              <span class="icon">{{ hasSpecialChar ? '✅' : '❌' }}</span>
              At least one special character (!@#$%^&*(),.?":{}|)
            </div>
          </div>
          <span class="form-help error" v-if="errors.password">{{ errors.password }}</span>
        </div>

        <div class="space-y-2">
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Confirm Password</label>
          <div class="relative">
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="Confirm your password"
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
              @input="validateConfirmPassword"
            />
            <button 
              type="button" 
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-600 dark:text-gray-400"
              @click="showConfirmPassword = !showConfirmPassword"
            >
              {{ showConfirmPassword ? '🔒' : '👁️' }}
            </button>
          </div>
          <div v-if="confirmPassword" class="password-match" :class="{ matched: passwordsMatch }">
            <span class="icon">{{ passwordsMatch ? '✅' : '❌' }}</span>
            {{ passwordsMatch ? 'Passwords match' : 'Passwords do not match' }}
          </div>
          <span class="form-help error" v-if="errors.confirmPassword">{{ errors.confirmPassword }}</span>
        </div>

        <button 
          @click="handleRegister" 
          :class="{ 'button-disabled': !isFormValid }"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Register
        </button>

        <p class="text-center text-sm text-gray-600 dark:text-gray-400">
          Already have an account? 
          <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300">
            Login
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import axios from 'axios';

export default {
  name: 'Register',
  setup() {
    const router = useRouter();
    const toast = useToast();
    
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const showPassword = ref(false);
    const showConfirmPassword = ref(false);
    const errors = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    });

    // Password requirements
    const hasMinLength = computed(() => password.value.length >= 8);
    const hasUpperCase = computed(() => /[A-Z]/.test(password.value));
    const hasLowerCase = computed(() => /[a-z]/.test(password.value));
    const hasNumber = computed(() => /[0-9]/.test(password.value));
    const hasSpecialChar = computed(() => /[!@#$%^&*(),.?":{}|]/.test(password.value));

    // Password matching
    const passwordsMatch = computed(() => 
      confirmPassword.value && password.value === confirmPassword.value
    );

    const validateUsername = () => {
      if (username.value.length < 3) {
        errors.value.username = 'Username must be at least 3 characters long';
      } else if (!/^[a-zA-Z0-9_]+$/.test(username.value)) {
        errors.value.username = 'Username can only contain letters, numbers, and underscores';
      } else {
        errors.value.username = '';
      }
    };

    const validateEmail = () => {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(email.value)) {
        errors.value.email = 'Please enter a valid email address';
      } else {
        errors.value.email = '';
      }
    };

    const validatePassword = () => {
      if (!hasMinLength.value) {
        errors.value.password = 'Password must be at least 8 characters long';
      } else if (!hasUpperCase.value) {
        errors.value.password = 'Password must contain at least one uppercase letter';
      } else if (!hasLowerCase.value) {
        errors.value.password = 'Password must contain at least one lowercase letter';
      } else if (!hasNumber.value) {
        errors.value.password = 'Password must contain at least one number';
      } else if (!hasSpecialChar.value) {
        errors.value.password = 'Password must contain at least one special character';
      } else {
        errors.value.password = '';
      }
      if (confirmPassword.value) {
        validateConfirmPassword();
      }
    };

    const validateConfirmPassword = () => {
      if (!confirmPassword.value) {
        errors.value.confirmPassword = 'Please confirm your password';
      } else if (!passwordsMatch.value) {
        errors.value.confirmPassword = 'Passwords do not match';
      } else {
        errors.value.confirmPassword = '';
      }
    };

    const isFormValid = computed(() => {
      return (
        username.value.length >= 3 &&
        /^[a-zA-Z0-9_]+$/.test(username.value) &&
        /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value) &&
        password.value.length >= 8 &&
        /[A-Z]/.test(password.value) &&
        /[a-z]/.test(password.value) &&
        /[0-9]/.test(password.value) &&
        /[!@#$%^&*(),.?":{}|]/.test(password.value) &&
        password.value === confirmPassword.value &&
        !Object.values(errors.value).some(error => error !== '')
      );
    });

    const handleRegister = async () => {
      // Validate all fields
      validateUsername();
      validateEmail();
      validatePassword();
      validateConfirmPassword();

      if (!isFormValid.value) {
        toast.error('Please fix the form errors before submitting');
        return;
      }

      try {
        await axios.post('/api/auth/register', {
          username: username.value,
          email: email.value,
          password: password.value
        });

        toast.success('Registration successful! Please login.');
        router.push('/login');
      } catch (error) {
        toast.error(error.response?.data?.error || 'Registration failed');
      }
    };

    return {
      username,
      email,
      password,
      confirmPassword,
      showPassword,
      showConfirmPassword,
      errors,
      isFormValid,
      handleRegister,
      hasMinLength,
      hasUpperCase,
      hasLowerCase,
      hasNumber,
      hasSpecialChar,
      passwordsMatch
    };
  }
};
</script>

<style scoped>
.min-h-screen {
  min-height: 100vh;
}

.bg-gray-100 {
  background-color: #f7fafc;
}

.dark\:bg-gray-900 {
  background-color: #1a1d23;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.justify-center {
  justify-content: center;
}

.py-12 {
  padding-top: 3rem;
  padding-bottom: 3rem;
}

.px-4 {
  padding-left: 1rem;
  padding-right: 1rem;
}

.sm\:px-6 {
  padding-left: 1.5rem;
  padding-right: 1.5rem;
}

.lg\:px-8 {
  padding-left: 2rem;
  padding-right: 2rem;
}

.max-w-md {
  max-width: 720px;
}

.w-full {
  width: 100%;
}

.space-y-8 {
  margin-top: 2rem;
  margin-bottom: 2rem;
}

.bg-white {
  background-color: #ffffff;
}

.dark\:bg-gray-800 {
  background-color: #2d3748;
}

.p-8 {
  padding: 2rem;
}

.rounded-lg {
  border-radius: 0.5rem;
}

.shadow-md {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.text-3xl {
  font-size: 1.875rem;
}

.font-bold {
  font-weight: bold;
}

.text-center {
  text-align: center;
}

.text-gray-900 {
  color: #1a1d23;
}

.dark\:text-white {
  color: #ffffff;
}

.mt-8 {
  margin-top: 2rem;
}

.space-y-6 {
  margin-top: 1.5rem;
  margin-bottom: 1.5rem;
}

.block {
  display: block;
}

.text-sm {
  font-size: 0.875rem;
}

.font-medium {
  font-weight: medium;
}

.text-gray-700 {
  color: #4a5568;
}

.dark\:text-gray-300 {
  color: #9ca3af;
}

.appearance-none {
  appearance: none;
}

.relative {
  position: relative;
}

.block {
  display: block;
}

.w-full {
  width: 100%;
}

.px-3 {
  padding-left: 0.75rem;
  padding-right: 0.75rem;
}

.py-2 {
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}

.border {
  border-width: 1px;
  border-style: solid;
}

.border-gray-300 {
  border-color: #d1d5db;
}

.dark\:border-gray-600 {
  border-color: #4a5568;
}

.placeholder-gray-500 {
  color: #a0aec0;
}

.dark\:placeholder-gray-400 {
  color: #9ca3af;
}

.text-gray-900 {
  color: #1a1d23;
}

.dark\:text-white {
  color: #ffffff;
}

.rounded-md {
  border-radius: 0.375rem;
}

.focus\:outline-none {
  outline: none;
}

.focus\:ring-indigo-500 {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.5);
}

.focus\:border-indigo-500 {
  border-color: #6366f1;
}

.focus\:z-10 {
  z-index: 10;
}

.sm\:text-sm {
  font-size: 0.875rem;
}

.bg-white {
  background-color: #ffffff;
}

.dark\:bg-gray-700 {
  background-color: #2d3748;
}

.absolute {
  position: absolute;
}

.inset-y-0 {
  top: 0;
  bottom: 0;
}

.right-0 {
  right: 0;
}

.pr-3 {
  padding-right: 0.75rem;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.text-gray-600 {
  color: #718096;
}

.dark\:text-gray-400 {
  color: #9ca3af;
}

.mt-1 {
  margin-top: 0.25rem;
}

.text-sm {
  font-size: 0.875rem;
}

.text-red-600 {
  color: #e53e3e;
}

.dark\:text-red-400 {
  color: #dc3545;
}

.password-requirements {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: #666;
  padding: 0.5rem;
  border-radius: 0.25rem;
  background: #f8f9fa;
}

.requirement {
  display: flex;
  align-items: center;
  margin: 0.25rem 0;
  opacity: 0.7;
  transition: all 0.2s;
}

.requirement.met {
  opacity: 1;
  color: #4CAF50;
}

.requirement .icon {
  margin-right: 0.5rem;
  font-size: 1rem;
  min-width: 1rem;
  text-align: center;
}

.password-match {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: #dc3545;
  display: flex;
  align-items: center;
  transition: all 0.2s;
}

.password-match.matched {
  color: #4CAF50;
}

.password-match .icon {
  margin-right: 0.5rem;
  font-size: 1rem;
  min-width: 1rem;
  text-align: center;
}

.group {
  display: inline-block;
  vertical-align: middle;
}

.relative {
  position: relative;
}

.w-full {
  width: 100%;
}

.flex {
  display: flex;
}

.justify-center {
  justify-content: center;
}

.py-2 {
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}

.px-4 {
  padding-left: 1rem;
  padding-right: 1rem;
}

.border {
  border-width: 1px;
  border-style: solid;
}

.border-transparent {
  border-color: transparent;
}

.text-sm {
  font-size: 0.875rem;
}

.font-medium {
  font-weight: medium;
}

.text-white {
  color: #ffffff;
}

.bg-indigo-600 {
  background-color: #5a67d8;
}

.hover\:bg-indigo-700 {
  background-color: #4f46e5;
}

.focus\:outline-none {
  outline: none;
}

.focus\:ring-2 {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
}

.focus\:ring-offset-2 {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
}

.focus\:ring-indigo-500 {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.5);
}

.disabled\:opacity-50 {
  opacity: 0.5;
}

.disabled\:cursor-not-allowed {
  cursor: not-allowed;
}

.button-disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: #ccc !important;
}

.text-center {
  text-align: center;
}

.text-sm {
  font-size: 0.875rem;
}

.text-gray-600 {
  color: #718096;
}

.dark\:text-gray-400 {
  color: #9ca3af;
}

.font-medium {
  font-weight: medium;
}

.text-indigo-600 {
  color: #5a67d8;
}

.hover\:text-indigo-500 {
  color: #4f46e5;
}

.dark\:text-indigo-400 {
  color: #6574cd;
}

.dark\:hover\:text-indigo-300 {
  color: #7a84dc;
}
</style>
