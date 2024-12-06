<template>
  <div class="register-container">
    <div class="register-card">
      <h1>WIRA Dashboard Register</h1>
      <div class="register-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="Choose a username"
            @blur="validateUsername"
          />
          <span class="form-help error" v-if="errors.username">{{ errors.username }}</span>
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="Enter your email"
            @blur="validateEmail"
          />
          <span class="form-help error" v-if="errors.email">{{ errors.email }}</span>
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <div class="password-input-container">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Choose a password"
              @blur="validatePassword"
            />
            <button 
              type="button" 
              class="toggle-password"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '🔒' : '👁️' }}
            </button>
          </div>
          <div v-if="password" class="password-info">
            Characters: {{ password.split('').map(char => 
              char === char.toUpperCase() ? '⬆️' : '⬇️'
            ).join(' ') }}
          </div>
          <span class="form-help error" v-if="errors.password">{{ errors.password }}</span>
        </div>
        <div class="form-group">
          <label for="confirmPassword">Confirm Password</label>
          <div class="password-input-container">
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="Confirm your password"
              @blur="validateConfirmPassword"
            />
            <button 
              type="button" 
              class="toggle-password"
              @click="showConfirmPassword = !showConfirmPassword"
            >
              {{ showConfirmPassword ? '🔒' : '👁️' }}
            </button>
          </div>
          <span class="form-help error" v-if="errors.confirmPassword">{{ errors.confirmPassword }}</span>
        </div>
        <button @click="handleRegister" :class="{ 'button-disabled': !isFormValid }">
          Register
        </button>
        <p class="login-link">
          Already have an account? <router-link to="/login">Login</router-link>
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
      if (password.value.length < 8) {
        errors.value.password = 'Password must be at least 8 characters long';
      } else if (!/[A-Z]/.test(password.value)) {
        errors.value.password = 'Password must contain at least one uppercase letter';
      } else if (!/[a-z]/.test(password.value)) {
        errors.value.password = 'Password must contain at least one lowercase letter';
      } else if (!/[0-9]/.test(password.value)) {
        errors.value.password = 'Password must contain at least one number';
      } else {
        errors.value.password = '';
      }
      if (confirmPassword.value) {
        validateConfirmPassword();
      }
    };

    const validateConfirmPassword = () => {
      if (password.value !== confirmPassword.value) {
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
      handleRegister
    };
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f5f5;
}

.register-card {
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

.register-form {
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

#password, #confirmPassword {
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
  font-size: 1.2em;
}

.password-info {
  margin-top: 5px;
  font-family: monospace;
  font-size: 14px;
  color: #666;
  letter-spacing: 2px;
}

.form-help {
  display: block;
  margin-top: 5px;
  font-size: 14px;
}

.error {
  color: #dc3545;
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

.button-disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: #ccc !important;
}

button:not(.button-disabled):hover {
  background: #45a049;
}

.login-link {
  text-align: center;
  margin-top: 1rem;
  color: #4a5568;
}

.login-link a {
  color: #4CAF50;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>
