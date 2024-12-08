<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-wira-secondary">Dashboard</h2>
    
    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-wira-primary/20 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-bold text-wira-secondary mb-2">Total Players</h3>
        <p class="text-4xl font-bold">27,401</p>
      </div>
      <div class="bg-wira-primary/20 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-bold text-wira-secondary mb-2">Total Scores</h3>
        <p class="text-4xl font-bold">205,337</p>
      </div>
      <div class="bg-wira-primary/20 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-bold text-wira-secondary mb-2">Classes</h3>
        <p class="text-4xl font-bold">8</p>
      </div>
    </div>

    <!-- Class Distribution -->
    <div class="bg-wira-primary/10 p-6 rounded-lg shadow-lg">
      <h3 class="text-2xl font-bold text-wira-secondary mb-4">Class Distribution</h3>
      <div class="h-[400px]">
        <ClassDistributionChart />
      </div>
    </div>

    <!-- Recent Top Scores -->
    <div class="bg-wira-primary/10 p-6 rounded-lg shadow-lg">
      <h3 class="text-2xl font-bold text-wira-secondary mb-4">Top Players</h3>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-wira-accent/30">
              <th class="text-left py-3 px-4">Rank</th>
              <th class="text-left py-3 px-4">Username</th>
              <th class="text-left py-3 px-4">Class</th>
              <th class="text-right py-3 px-4">Score</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(player, index) in topPlayers" :key="index" 
                class="border-b border-wira-accent/20 hover:bg-wira-accent/10">
              <td class="py-3 px-4">{{ index + 1 }}</td>
              <td class="py-3 px-4">{{ player.username }}</td>
              <td class="py-3 px-4">Class {{ player.classId }}</td>
              <td class="py-3 px-4 text-right">{{ player.score.toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import { useRouter } from 'vue-router'
import ClassDistributionChart from '../components/ClassDistributionChart.vue'

const toast = useToast()
const router = useRouter()
const topPlayers = ref([])

onMounted(async () => {
  try {
    const token = localStorage.getItem('access_token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await fetch('/api/rankings?page=1&per_page=10', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (response.status === 401) {
      toast.error('Session expired. Please login again.')
      router.push('/login')
      return
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    if (data && data.data) {
      topPlayers.value = data.data.map(player => ({
        username: player.username,
        classId: player.class_id,
        score: player.highest_score
      }))
    } else {
      console.error('Invalid data format received:', data)
      toast.error('Failed to load rankings data')
    }
  } catch (error) {
    console.error('Error fetching top players:', error)
    toast.error('Failed to load rankings')
  }
})
</script>
