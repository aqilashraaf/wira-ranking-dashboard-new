<template>
  <div class="profile-container">
    <div class="profile-section">
      <h2>Profile Information</h2>
      <div class="profile-card">
        <h1>Profile Settings</h1>
        <div class="profile-form">
          <div class="form-group">
            <label for="username">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              :disabled="true"
            />
          </div>
          <div class="form-group">
            <label for="email">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              :disabled="true"
            />
          </div>
          <div class="form-group">
            <label for="currentPassword">Current Password</label>
            <input
              id="currentPassword"
              v-model="currentPassword"
              type="password"
              placeholder="Enter current password"
            />
          </div>
          <div class="form-group">
            <label for="newPassword">New Password</label>
            <input
              id="newPassword"
              v-model="newPassword"
              type="password"
              placeholder="Enter new password"
            />
          </div>
          <div class="form-group">
            <label for="confirmNewPassword">Confirm New Password</label>
            <input
              id="confirmNewPassword"
              v-model="confirmNewPassword"
              type="password"
              placeholder="Confirm new password"
            />
          </div>
          <button @click="updatePassword" :disabled="!isValid">
            Update Password
          </button>
        </div>
      </div>
    </div>

    <div class="twofa-section">
      <h2>Two-Factor Authentication</h2>
      <div class="twofa-card">
        <TwoFactorSetup />
      </div>
    </div>

    <div class="activities-section">
      <h2>Recent Activities</h2>
      <div class="activities-list" v-if="activities.length">
        <div v-for="activity in activities" :key="activity.timestamp" class="activity-item">
          <div class="activity-type">{{ activity.type }}</div>
          <div class="activity-description">{{ activity.description }}</div>
          <div class="activity-time">{{ formatDate(activity.timestamp) }}</div>
        </div>
      </div>
      <div v-else class="no-activities">
        No recent activities found
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import axios from 'axios';
import { useRouter } from 'vue-router';
import TwoFactorSetup from '@/components/TwoFactorSetup.vue';

export default {
  name: 'Profile',
  components: {
    TwoFactorSetup
  },
  setup() {
    const toast = useToast();
    const router = useRouter();
    
    const username = ref('');
    const email = ref('');
    const currentPassword = ref('');
    const newPassword = ref('');
    const confirmNewPassword = ref('');
    const activities = ref([]);
    const twoFactorEnabled = ref(false);

    const isValid = computed(() => {
      return (
        currentPassword.value.length > 0 &&
        newPassword.value.length >= 8 &&
        newPassword.value === confirmNewPassword.value
      );
    });

    const fetchProfile = async () => {
      try {
        console.log('Fetching profile...');
        const token = localStorage.getItem('access_token');
        console.log('Token:', token);
        
        const response = await axios.get('/api/user/profile', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        console.log('Profile response:', response.data);
        
        if (response.data) {
          username.value = response.data.username || '';
          email.value = response.data.email || '';
          twoFactorEnabled.value = response.data.two_factor_enabled || false;
        } else {
          throw new Error('No profile data received');
        }
      } catch (error) {
        console.error('Profile fetch error:', error);
        console.error('Error response:', error.response);
        if (error.response?.status === 401) {
          toast.error('Session expired. Please login again.');
          router.push('/login');
        } else {
          toast.error(error.response?.data?.error || 'Failed to load profile');
        }
      }
    };

    const fetchActivities = async () => {
      try {
        const response = await axios.get('/api/user/activities')
        activities.value = response.data.activities
      } catch (error) {
        console.error('Error fetching activities:', error)
        toast.error('Failed to load activities')
      }
    }

    const updatePassword = async () => {
      if (!isValid.value) return;

      try {
        await axios.post('/api/user/change-password', {
          current_password: currentPassword.value,
          new_password: newPassword.value
        });

        toast.success('Password updated successfully');
        currentPassword.value = '';
        newPassword.value = '';
        confirmNewPassword.value = '';
      } catch (error) {
        toast.error(error.response?.data?.error || 'Failed to update password');
      }
    };

    const formatDate = (timestamp) => {
      return new Date(timestamp).toLocaleString()
    }

    onMounted(async () => {
      await Promise.all([fetchProfile(), fetchActivities()])
    });

    return {
      username,
      email,
      currentPassword,
      newPassword,
      confirmNewPassword,
      activities,
      twoFactorEnabled,
      formatDate,
      isValid,
      updatePassword
    };
  }
};
</script>

<style scoped>
.profile-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 60px);
  padding: 20px;
  background-color: #f5f5f5;
}

.profile-section {
  margin-bottom: 2rem;
}

.profile-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
}

.profile-card h1 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 2rem;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #2c3e50;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.form-group input {
  padding: 0.75rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.2s;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.form-group input:not([disabled]):hover {
  border-color: #409eff;
}

.form-group input:not([disabled]):focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

/* Style for disabled inputs */
.form-group input:disabled {
  background-color: #f8f9fa;
  border-color: #e9ecef;
  color: #2c3e50;
  font-weight: 500;
  cursor: not-allowed;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 1.1rem;
  letter-spacing: 0.5px;
}

button {
  background-color: #409eff;
  color: white;
  padding: 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

button:hover:not(:disabled) {
  background-color: #66b1ff;
}

button:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}

.activities-section {
  margin-top: 2rem;
  padding: 1rem;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.activities-list {
  margin-top: 1rem;
}

.activity-item {
  padding: 0.75rem;
  border-bottom: 1px solid #eee;
  display: grid;
  grid-template-columns: 120px 1fr 180px;
  gap: 1rem;
  align-items: center;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-type {
  font-weight: bold;
  color: #2c3e50;
}

.activity-description {
  color: #666;
}

.activity-time {
  color: #999;
  font-size: 0.9em;
  text-align: right;
}

.no-activities {
  text-align: center;
  color: #666;
  padding: 2rem;
}

.twofa-section {
  margin: 2rem 0;
}

.twofa-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.twofa-setup, .twofa-enabled {
  text-align: center;
}

.twofa-setup p, .twofa-enabled p {
  color: #666;
  margin-bottom: 1rem;
}

.setup-button, .disable-button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.setup-button {
  background-color: #4CAF50;
  color: white;
}

.setup-button:hover {
  background-color: #45a049;
}

.disable-button {
  background-color: #f44336;
  color: white;
}

.disable-button:hover {
  background-color: #d32f2f;
}
</style>
