<script setup>
import { onMounted, ref, watch } from "vue"
import { Chart, registerables } from "chart.js"
import "chartjs-adapter-luxon"
Chart.register(...registerables)

const props = defineProps({ chartData: Object })
const canvas = ref(null)
let chartInstance = null

onMounted(() => {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: props.chartData,
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          intersect: false,
          mode: "index",
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
  () => props.chartData,
  (newData) => {
    if (chartInstance) {
      chartInstance.data = newData
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