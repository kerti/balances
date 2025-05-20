<script setup>
import { onMounted, ref, watch } from "vue"
import { Chart, registerables } from "chart.js"
Chart.register(...registerables)

const props = defineProps({ chartData: Object })
const canvas = ref(null)
let chartInstance = null

onMounted(() => {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: props.chartData,
      options: { responsive: true, maintainAspectRatio: false },
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