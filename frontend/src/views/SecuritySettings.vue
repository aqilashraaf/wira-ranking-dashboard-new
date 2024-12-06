<template>
  <div class="security-settings">
    <div class="security-card">
      <h1>Security Settings</h1>
      <div class="settings-section">
        <h2>Two-Factor Authentication</h2>
        <p class="status-text">
          Status: <span :class="{ 'enabled': is2FAEnabled, 'disabled': !is2FAEnabled }">
            {{ is2FAEnabled ? 'Enabled' : 'Disabled' }}
          </span>
        </p>
        <TwoFactorSetup v-if="showSetup" />
        <div v-else class="action-buttons">
          <button v-if="!is2FAEnabled" @click="showSetup = true" class="setup-button">
            Set up 2FA
          </button>
          <button v-else @click="disable2FA" class="disable-button">
            Disable 2FA
          </button>
        </div>
      </div>
      <div class="settings-section">
        <h2>Recent Activity</h2>
        <div class="activity-list" v-if="activities.length">
          <div v-for="activity in activities" :key="activity.id" class="activity-item">
            <span class="activity-type">{{ activity.type }}</span>
            <span class="activity-time">{{ formatDate(activity.timestamp) }}</span>
          </div>
        </div>
        <p v-else class="no-activity">No recent activity</p>
      </div>
      <div class="section">
        <h3>Change Password</h3>
        <form @submit.prevent="changePassword" class="password-form">
          <div class="form-group">
            <label for="currentPassword">Current Password</label>
            <input
              id="currentPassword"
              v-model="currentPassword"
              type="password"
              required
            />
          </div>
          <div class="form-group">
            <label for="newPassword">New Password</label>
            <input
              id="newPassword"
              v-model="newPassword"
              type="password"
              required
            />
          </div>
          <div class="form-group">
            <label for="confirmNewPassword">Confirm New Password</label>
            <input
              id="confirmNewPassword"
              v-model="confirmNewPassword"
              type="password"
              required
            />
          </div>
          <button type="submit" :disabled="!isPasswordValid">
            Change Password
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import axios from 'axios';
import TwoFactorSetup from '@/components/TwoFactorSetup.vue';

export default {
  name: 'SecuritySettings',
  components: {
    TwoFactorSetup
  },
  setup() {
    const toast = useToast();
    const is2FAEnabled = ref(false);
    const showSetup = ref(false);
    const activities = ref([]);
    const currentPassword = ref('');
    const newPassword = ref('');
    const confirmNewPassword = ref('');

    const isPasswordValid = computed(() => {
      return (
        currentPassword.value.length > 0 &&
        newPassword.value.length >= 8 &&
        newPassword.value === confirmNewPassword.value
      );
    });

    const check2FAStatus = async () => {
      try {
        const response = await axios.get('/api/2fa/status');
        is2FAEnabled.value = response.data.enabled;
      } catch (error) {
        toast.error('Failed to check 2FA status');
      }
    };

    const disable2FA = async () => {
      try {
        await axios.post('/api/2fa/disable');
        is2FAEnabled.value = false;
        toast.success('2FA has been disabled');
      } catch (error) {
        toast.error('Failed to disable 2FA');
      }
    };

    const fetchActivities = async () => {
      try {
        const response = await axios.get('/api/user/activities');
        activities.value = response.data.activities;
      } catch (error) {
        toast.error('Failed to load activities');
      }
    };

    const formatDate = (timestamp) => {
      return new Date(timestamp).toLocaleString();
    };

    const changePassword = async () => {
      try {
        await axios.post('/api/auth/change-password', {
          current_password: currentPassword.value,
          new_password: newPassword.value
        });
        
        toast.success('Password changed successfully');
        currentPassword.value = '';
        newPassword.value = '';
        confirmNewPassword.value = '';
      } catch (error) {
        toast.error(error.response?.data?.error || 'Failed to change password');
      }
    };

    onMounted(() => {
      check2FAStatus();
      fetchActivities();
    });

    return {
      is2FAEnabled,
      showSetup,
      activities,
      disable2FA,
      formatDate,
      currentPassword,
      newPassword,
      confirmNewPassword,
      isPasswordValid,
      changePassword
    };
  }
};
</script>

<style scoped>
.security-settings {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  background: #f5f5f5;
  padding: 2rem;
}

.security-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 600px;
}

h1 {
  color: #2c3e50;
  margin-bottom: 2rem;
  text-align: center;
}

.settings-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #e2e8f0;
}

.settings-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

h2 {
  color: #2c3e50;
  margin-bottom: 1rem;
  font-size: 1.25rem;
}

.status-text {
  margin-bottom: 1rem;
  font-size: 1rem;
  color: #4a5568;
}

.status-text span {
  font-weight: bold;
}

.enabled {
  color: #4CAF50;
}

.disabled {
  color: #dc3545;
}

.action-buttons {
  display: flex;
  gap: 1rem;
}

button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.setup-button {
  background: #4CAF50;
  color: white;
}

.setup-button:hover {
  background: #45a049;
}

.disable-button {
  background: #dc3545;
  color: white;
}

.disable-button:hover {
  background: #c82333;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.activity-type {
  color: #2c3e50;
  font-weight: 500;
}

.activity-time {
  color: #6c757d;
  font-size: 0.875rem;
}

.no-activity {
  color: #6c757d;
  text-align: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.section {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
}

h3 {
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

.password-form {
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
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 1rem;
}

input:focus {
  outline: none;
  border-color: #4CAF50;
  box-shadow: 0 0 0 1px #4CAF50;
}

button[type="submit"] {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-top: 0.5rem;
}

button[type="submit"]:hover {
  background: #45a049;
}

button[type="submit"]:disabled {
  background: #cccccc;
  cursor: not-allowed;
}
</style>
