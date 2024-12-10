<template>
  <div class="max-w-md mx-auto p-5">
    <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-md text-center">
      <h2 class="text-2xl font-semibold text-gray-800 dark:text-white mb-4">Two-Factor Authentication Required</h2>
      <p class="text-gray-600 dark:text-gray-300 mb-6">Please enter the 6-digit code from your authenticator app.</p>
      <div class="space-y-4">
        <input
          v-model="code"
          type="text"
          placeholder="Enter 6-digit code"
          maxlength="6"
          class="w-40 px-4 py-2 text-center text-lg tracking-widest border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-800 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          @keyup.enter="verify"
        />
        <button 
          @click="verify" 
          :disabled="!isValidCode"
          class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition duration-200"
        >
          Verify
        </button>
      </div>
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
</style>
