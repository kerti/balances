<script setup>
import "chartjs-adapter-luxon"
import { onActivated, onMounted, ref, watch } from "vue"
import { Chart, registerables } from "chart.js"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"

Chart.register(...registerables)

const bankAccountsStore = useBankAccountsStore()
const canvas = ref(null)
let chartInstance = null

function renderChart() {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: {
        datasets: bankAccountsStore.chartData,
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
  () => bankAccountsStore.chartData,
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