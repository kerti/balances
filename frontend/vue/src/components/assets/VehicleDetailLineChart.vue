<script setup>
import "chartjs-adapter-luxon"
import { onActivated, onMounted, ref, watch } from "vue"
import { Chart, registerables } from "chart.js"
import { useVehiclesStore } from "@/stores/vehiclesStore"

Chart.register(...registerables)

const vehiclesStore = useVehiclesStore()
const canvas = ref(null)
let chartInstance = null

function renderChart() {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: {
        datasets: vehiclesStore.dvChartData,
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          intersect: false,
          mode: "nearest",
          axis: "xy",
        },
        scales: {
          x: {
            type: "time",
            time: {
              tooltipFormat: "dd LLL yy",
            },
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

watch(
  () => vehiclesStore.dvChartData,
  (newData) => {
    if (chartInstance) {
      chartInstance.data = {
        datasets: newData,
      }
      chartInstance.update()
    }
  },
  { deep: true, immediate: true }
)
</script>

<template>
  <div class="flex-1 h-full">
    <canvas ref="canvas"></canvas>
  </div>
</template>