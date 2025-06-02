<script setup>
import "chartjs-adapter-luxon"
import { onMounted, ref, watch } from "vue"
import { Chart, registerables } from "chart.js"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"

Chart.register(...registerables)

const bankAccountsStore = useBankAccountsStore()
const canvas = ref(null)
let chartInstance = null

onMounted(() => {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: bankAccountsStore.chartData,
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
  { deep: true }
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