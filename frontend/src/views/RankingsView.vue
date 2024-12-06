<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-wira-secondary">Rankings</h2>

    <!-- Filters -->
    <div class="flex flex-wrap gap-4 bg-wira-primary/10 p-4 rounded-lg">
      <div class="flex-1 min-w-[200px]">
        <label class="block text-sm font-medium mb-2">Class</label>
        <select 
          v-model="selectedClass" 
          class="w-full bg-wira-background border border-wira-accent rounded px-3 py-2"
        >
          <option value="0">All Classes</option>
          <option v-for="n in 8" :key="n" :value="n">Class {{ n }}</option>
        </select>
      </div>
      <div class="flex-1 min-w-[200px]">
        <label class="block text-sm font-medium mb-2">Search Player</label>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Enter username..."
          class="w-full bg-wira-background border border-wira-accent rounded px-3 py-2"
        >
      </div>
    </div>

    <!-- Rankings Table -->
    <div class="bg-wira-primary/10 rounded-lg shadow-lg overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-wira-primary/30">
              <th class="text-left py-3 px-4">Rank</th>
              <th class="text-left py-3 px-4">Username</th>
              <th class="text-left py-3 px-4">Class</th>
              <th class="text-right py-3 px-4">Score</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="player in players" 
                :key="player.username" 
                class="border-b border-wira-accent/20 hover:bg-wira-accent/10">
              <td class="py-3 px-4">{{ player.rank }}</td>
              <td class="py-3 px-4">{{ player.username }}</td>
              <td class="py-3 px-4">Class {{ player.class_id }}</td>
              <td class="py-3 px-4 text-right">{{ player.highest_score.toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-wira-secondary border-t-transparent mx-auto"></div>
      </div>

      <!-- Pagination -->
      <div class="p-4 flex justify-between items-center bg-wira-primary/20">
        <div class="text-sm">
          Showing {{ players.length }} of {{ totalPlayers }} players
        </div>
        <div class="flex gap-2">
          <button 
            @click="prevPage" 
            :disabled="currentPage === 1"
            class="px-4 py-2 rounded bg-wira-accent disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span class="px-4 py-2">Page {{ currentPage }}</span>
          <button 
            @click="nextPage"
            :disabled="currentPage * perPage >= totalPlayers"
            class="px-4 py-2 rounded bg-wira-accent disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import axios from 'axios'

const players = ref([])
const totalPlayers = ref(0)
const currentPage = ref(1)
const perPage = ref(20)
const selectedClass = ref(0)
const searchQuery = ref('')
const loading = ref(false)

// Create axios instance with default config
const api = axios.create({
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

const fetchRankings = async () => {
  loading.value = true
  try {
    let url = `/api/rankings?page=${currentPage.value}&per_page=${perPage.value}`
    if (selectedClass.value && selectedClass.value > 0) {
      url += `&class_id=${selectedClass.value}`
    }
    const response = await api.get(url)
    if (response.data && response.data.data) {
      players.value = response.data.data
      totalPlayers.value = response.data.total || 0
    } else {
      console.error('Invalid data format received:', response.data)
      players.value = []
      totalPlayers.value = 0
    }
  } catch (error) {
    console.error('Error fetching rankings:', error)
    players.value = []
    totalPlayers.value = 0
  } finally {
    loading.value = false
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value * perPage.value < totalPlayers.value) {
    currentPage.value++
  }
}

watch([currentPage, selectedClass], () => {
  fetchRankings()
})

watch(searchQuery, async (newQuery) => {
  if (newQuery.length > 2) {
    loading.value = true
    try {
      const response = await api.get(`/api/rankings/search?username=${newQuery}`)
      players.value = response.data.data
      totalPlayers.value = response.data.total
    } catch (error) {
      console.error('Error searching players:', error)
    } finally {
      loading.value = false
    }
  } else if (newQuery.length === 0) {
    fetchRankings()
  }
})

onMounted(() => {
  fetchRankings()
})
</script>
