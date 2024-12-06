<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-wira-secondary">Statistics</h2>

    <!-- Class Statistics -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div v-for="stat in classStats" :key="stat.class_id" 
           class="bg-wira-primary/10 p-6 rounded-lg shadow-lg">
        <h3 class="text-xl font-bold text-wira-secondary mb-4">Class {{ stat.class_id }}</h3>
        <div class="space-y-4">
          <div class="flex justify-between items-center">
            <span>Total Players:</span>
            <span class="font-bold">{{ stat.player_count.toLocaleString() }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span>Average Score:</span>
            <span class="font-bold">{{ Math.round(stat.average_score).toLocaleString() }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span>Highest Score:</span>
            <span class="font-bold text-wira-secondary">{{ stat.highest_score.toLocaleString() }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span>Lowest Score:</span>
            <span class="font-bold">{{ stat.lowest_score.toLocaleString() }}</span>
          </div>
          
          <!-- Score Distribution Bar -->
          <div class="mt-4">
            <div class="h-2 bg-wira-background rounded-full overflow-hidden">
              <div 
                class="h-full bg-wira-secondary"
                :style="{ width: ((stat.average_score / stat.highest_score) * 100) + '%' }"
              ></div>
            </div>
            <div class="flex justify-between text-sm mt-1">
              <span>{{ stat.lowest_score }}</span>
              <span>{{ stat.highest_score }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const api = axios.create({
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  withCredentials: false,
  mode: 'no-cors'
})

const classStats = ref([])

onMounted(async () => {
  try {
    const response = await api.get('/api/rankings/stats')
    classStats.value = response.data.data
  } catch (error) {
    console.error('Error fetching class statistics:', error)
  }
})
</script>
