<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-3xl mx-auto space-y-8">
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Profile Information</h2>
        <div class="space-y-6">
          <div class="space-y-2">
            <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              :disabled="true"
              class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white disabled:opacity-75 disabled:cursor-not-allowed"
            />
          </div>

          <div class="space-y-2">
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              :disabled="true"
              class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white disabled:opacity-75 disabled:cursor-not-allowed"
            />
          </div>

          <div class="space-y-2">
            <label for="currentPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Current Password</label>
            <input
              id="currentPassword"
              v-model="currentPassword"
              type="password"
              placeholder="Enter current password"
              class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            />
          </div>

          <div class="space-y-2">
            <label for="newPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300">New Password</label>
            <input
              id="newPassword"
              v-model="newPassword"
              type="password"
              placeholder="Enter new password"
              class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            />
          </div>

          <div class="space-y-2">
            <label for="confirmNewPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Confirm New Password</label>
            <input
              id="confirmNewPassword"
              v-model="confirmNewPassword"
              type="password"
              placeholder="Confirm new password"
              class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            />
          </div>

          <div class="flex justify-end space-x-4">
            <button 
              @click="updatePassword"
              :disabled="!isValid"
              class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed dark:disabled:opacity-40"
            >
              Update Password
            </button>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Recent Activities</h2>
        <div class="space-y-4">
          <div v-if="activities.length === 0" class="text-gray-500 dark:text-gray-400 text-center py-4">
            No recent activities
          </div>
          <div v-else class="divide-y divide-gray-200 dark:divide-gray-700">
            <div v-for="activity in activities" :key="activity.timestamp" class="py-4">
              <div class="flex items-center space-x-4">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ activity.type }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(activity.timestamp) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
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
/* Add your custom styles here */
</style>
