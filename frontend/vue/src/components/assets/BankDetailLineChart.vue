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
        datasets: bankAccountsStore.detailChartData,
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
  () => bankAccountsStore.detailChartData,
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
  <canvas ref="canvas"></canvas>
</template>

<style scoped>
canvas {
  width: 100% !important;
  height: 300px !important;
}
</style>