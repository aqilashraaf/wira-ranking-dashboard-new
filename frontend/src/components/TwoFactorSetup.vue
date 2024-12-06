<template>
  <div class="two-factor-setup">
    <div v-if="!isEnabled" class="setup-container">
      <h2>Setup Two-Factor Authentication</h2>
      <div v-if="qrCode" class="qr-section">
        <p>Scan this QR code with your authenticator app:</p>
        <div class="qr-wrapper">
          <img :src="qrCodeDataUrl" alt="2FA QR Code" class="qr-code" />
        </div>
        <p class="secret-key">
          Or manually enter this code: <code>{{ secret }}</code>
        </p>
        <div class="verify-section">
          <input
            v-model="verificationCode"
            type="text"
            placeholder="Enter 6-digit code"
            maxlength="6"
            class="verification-input"
          />
          <button @click="verify2FA" :disabled="!verificationCode" class="verify-button">
            Verify and Enable 2FA
          </button>
        </div>
      </div>
      <button v-else @click="setup2FA" class="setup-button">
        Set up Two-Factor Authentication
      </button>
    </div>
    <div v-else class="enabled-container">
      <h2>Two-Factor Authentication is Enabled</h2>
      <p>Your account is protected with 2FA.</p>
      <button @click="disable2FA" class="disable-button">
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
.two-factor-setup {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.setup-container, .enabled-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.qr-section {
  text-align: center;
}

.qr-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 20px 0;
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.qr-code {
  width: 200px;
  height: 200px;
}

.secret-key {
  background: #f8f9fa;
  padding: 10px;
  border-radius: 4px;
  margin: 20px 0;
  color: #2c3e50;
}

.secret-key code {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  color: #2c3e50;
}

.verify-section {
  margin-top: 20px;
}

.verification-input {
  width: 150px;
  padding: 8px;
  font-size: 16px;
  text-align: center;
  margin-right: 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-family: monospace;
  letter-spacing: 2px;
}

button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

button:hover {
  background: #45a049;
}

button:disabled {
  background: #cccccc;
  cursor: not-allowed;
}

.disable-button {
  background: #dc3545;
}

.disable-button:hover {
  background: #c82333;
}
</style>
