<template>
  <div class="max-w-lg mx-auto p-5">
    <div v-if="!isEnabled" class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-800 dark:text-white mb-4">Setup Two-Factor Authentication</h2>
      <div v-if="qrCode" class="space-y-4">
        <p class="text-gray-600 dark:text-gray-300">Scan this QR code with your authenticator app:</p>
        <div class="flex justify-center">
          <img :src="qrCodeDataUrl" alt="2FA QR Code" class="w-48 h-48" />
        </div>
        <p class="text-gray-600 dark:text-gray-300">
          Or manually enter this code: <code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">{{ secret }}</code>
        </p>
        <div class="space-y-3">
          <input
            v-model="verificationCode"
            type="text"
            placeholder="Enter 6-digit code"
            maxlength="6"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-gray-100 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm bg-white dark:bg-gray-700"
          />
          <button 
            @click="verify2FA" 
            :disabled="!verificationCode"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition duration-200"
          >
            Verify and Enable 2FA
          </button>
        </div>
      </div>
      <button 
        v-else 
        @click="setup2FA"
        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition duration-200"
      >
        Set up Two-Factor Authentication
      </button>
    </div>
    <div v-else class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-800 dark:text-white mb-4">Two-Factor Authentication is Enabled</h2>
      <p class="text-gray-600 dark:text-gray-300 mb-4">Your account is protected with 2FA.</p>
      <button 
        @click="disable2FA"
        class="w-full bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md transition duration-200"
      >
        Disable 2FA
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { ref, onMounted, watch } from 'vue';
import { useToast } from 'vue-toastification';
import QRCode from 'qrcode';

export default {
  name: 'TwoFactorSetup',
  setup() {
    const toast = useToast();
    const isEnabled = ref(false);
    const qrCode = ref('');
    const qrCodeDataUrl = ref('');
    const secret = ref('');
    const verificationCode = ref('');

    const generateQRCode = async (url) => {
      try {
        qrCodeDataUrl.value = await QRCode.toDataURL(url, {
          width: 200,
          margin: 2,
          color: {
            dark: '#2c3e50',
            light: '#ffffff'
          }
        });
      } catch (error) {
        console.error('Failed to generate QR code:', error);
        toast.error('Failed to generate QR code');
      }
    };

    watch(qrCode, async (newValue) => {
      if (newValue) {
        await generateQRCode(newValue);
      }
    });

    const checkStatus = async () => {
      try {
        const response = await axios.get('/api/2fa/status');
        isEnabled.value = response.data.two_factor_enabled;
      } catch (error) {
        toast.error('Failed to check 2FA status');
      }
    };

    const setup2FA = async () => {
      try {
        const response = await axios.post('/api/2fa/setup');
        qrCode.value = response.data.qr_url;
        secret.value = response.data.secret;
      } catch (error) {
        toast.error('Failed to setup 2FA');
      }
    };

    const verify2FA = async () => {
      try {
        await axios.post('/api/2fa/enable', {
          totp_code: verificationCode.value
        });
        toast.success('2FA enabled successfully');
        isEnabled.value = true;
        qrCode.value = '';
        secret.value = '';
        verificationCode.value = '';
      } catch (error) {
        toast.error(error.response?.data?.error || 'Failed to verify 2FA code');
      }
    };

    const disable2FA = async () => {
      try {
        await axios.post('/api/2fa/disable');
        toast.success('2FA disabled successfully');
        isEnabled.value = false;
      } catch (error) {
        toast.error('Failed to disable 2FA');
      }
    };

    onMounted(checkStatus);

    return {
      isEnabled,
      qrCode,
      qrCodeDataUrl,
      secret,
      verificationCode,
      setup2FA,
      verify2FA,
      disable2FA
    };
  }
};
</script>

<style scoped>
.form-help.error {
  @apply text-red-600 dark:text-red-400 text-sm mt-1;
}

.text-content {
  @apply text-gray-900 dark:text-gray-100;
}

.qr-code {
  @apply bg-white dark:bg-gray-200 p-4 rounded-lg inline-block;
}
</style>
