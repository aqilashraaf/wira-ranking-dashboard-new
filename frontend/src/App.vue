<template>
  <div class="min-h-screen flex flex-col bg-wira-background text-wira-text font-cinzel">
    <header class="bg-gradient-to-r from-wira-primary to-wira-secondary p-4 md:p-6 shadow-lg">
      <div class="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center gap-4">
        <div class="flex items-center gap-4">
          <img src="./assets/wira-logo.jpg" alt="WIRA Logo" class="h-12 w-auto" />
          <h1 class="text-2xl md:text-3xl font-bold text-wira-text drop-shadow-lg">WIRA Rankings</h1>
        </div>
        <nav class="flex flex-col md:flex-row gap-4 md:gap-8">
          <router-link 
            v-for="link in ['Dashboard', 'Rankings', 'Statistics']" 
            :key="link"
            :to="link === 'Dashboard' ? '/' : `/${link.toLowerCase()}`"
            class="px-4 py-2 text-lg text-wira-text hover:bg-wira-accent rounded-md transition-colors duration-300"
          >
            {{ link }}
          </router-link>
          <div class="dropdown">
            <button class="dropbtn">Account</button>
            <div class="dropdown-content">
              <router-link to="/profile">Profile</router-link>
              <router-link to="/profile/security">Security Settings</router-link>
              <a href="#" @click.prevent="handleLogout">Logout</a>
            </div>
          </div>
        </nav>
      </div>
    </header>

    <main class="flex-1 p-4 md:p-8 max-w-7xl mx-auto w-full">
      <router-view></router-view>
    </main>

    <footer class="bg-wira-primary text-center p-4 mt-auto">
      <p class="text-sm md:text-base"> 2024 WIRA Rankings - Developed by Aqash</p>
    </footer>
  </div>
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import axios from 'axios'

const router = useRouter()
const toast = useToast()

const isAuthenticated = computed(() => {
  return !!localStorage.getItem('access_token')
})

const handleLogout = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  delete axios.defaults.headers.common['Authorization']
  toast.success('Logged out successfully')
  router.push('/login')
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&display=swap');
@tailwind base;
@tailwind components;
@tailwind utilities;

.dropdown {
  position: relative;
  display: inline-block;
}

.dropbtn {
  background-color: transparent;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.dropdown-content {
  display: none;
  position: absolute;
  right: 0;
  background-color: white;
  min-width: 160px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  z-index: 1;
  border-radius: 4px;
}

.dropdown-content a {
  color: #2c3e50 !important;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
}

.dropdown-content a:hover {
  background-color: #f1f1f1;
}

.dropdown:hover .dropdown-content {
  display: block;
}
</style>
