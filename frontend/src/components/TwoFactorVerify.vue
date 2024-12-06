<template>
  <div class="two-factor-verify">
    <h2>Two-Factor Authentication Required</h2>
    <p>Please enter the 6-digit code from your authenticator app.</p>
    <div class="verify-form">
      <input
        v-model="code"
        type="text"
        placeholder="Enter 6-digit code"
        maxlength="6"
        class="code-input"
        @keyup.enter="verify"
      />
      <button @click="verify" :disabled="!isValidCode" class="verify-button">
        Verify
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

export default {
  name: 'TwoFactorVerify',
  props: {
    onVerify: {
      type: Function,
      required: true
    }
  },
  setup(props) {
    const code = ref('');

    const isValidCode = computed(() => {
      return code.value.length === 6 && /^\d+$/.test(code.value);
    });

    const verify = () => {
      if (isValidCode.value) {
        props.onVerify(code.value);
      }
    };

    return {
      code,
      isValidCode,
      verify
    };
  }
};
</script>

<style scoped>
.two-factor-verify {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

h2 {
  color: #2c3e50;
  margin-bottom: 15px;
}

p {
  color: #6c757d;
  margin-bottom: 20px;
}

.verify-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
  align-items: center;
}

.code-input {
  width: 150px;
  padding: 10px;
  font-size: 18px;
  text-align: center;
  border: 1px solid #ced4da;
  border-radius: 4px;
  letter-spacing: 2px;
}

.verify-button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 10px 30px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.verify-button:hover {
  background: #45a049;
}

.verify-button:disabled {
  background: #cccccc;
  cursor: not-allowed;
}
</style>
