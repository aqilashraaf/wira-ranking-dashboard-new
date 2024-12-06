<template>
  <div class="h-full">
    <Bar v-if="chartData" :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Bar } from 'vue-chartjs'
import axios from 'axios'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
)

const api = axios.create({
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

const chartData = ref(null)
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
      grid: {
        color: 'rgba(205, 133, 63, 0.1)'
      },
      ticks: {
        color: '#ffffff'
      }
    },
    x: {
      grid: {
        color: 'rgba(205, 133, 63, 0.1)'
      },
      ticks: {
        color: '#ffffff'
      }
    }
  },
  plugins: {
    legend: {
      display: false
    }
  }
}

onMounted(async () => {
  try {
    const response = await api.get('/api/rankings/stats')
    const data = response.data
    
    chartData.value = {
      labels: data.data.map(stat => `Class ${stat.class_id}`),
      datasets: [
        {
          data: data.data.map(stat => stat.player_count),
          backgroundColor: '#DAA520',
          borderColor: '#8B4513',
          borderWidth: 1
        }
      ]
    }
  } catch (error) {
    console.error('Error fetching class distribution:', error)
  }
})
</script>
