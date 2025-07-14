<script setup>
import "chartjs-adapter-luxon"
import { Chart, registerables } from "chart.js"
import { onActivated, onMounted, ref } from "vue"

Chart.register(...registerables)

const canvas = ref(null)
let chartInstance = null

function renderChart() {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "doughnut",
      data: {
        labels: ["Cash", "Vehicles", "Properties"],
        datasets: [{ data: [123, 456, 789] }],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: "top",
          },
        },
      },
    })
  }
}

function destroyChart() {
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
}

onMounted(() => {
  destroyChart()
  renderChart()
})

onActivated(() => {
  destroyChart()
  renderChart()
})
</script>

<template>
  <div class="flex justify-center items-center w-full h-full">
    <canvas ref="canvas" class="max-w-full max-h-full"></canvas>
  </div>
</template>