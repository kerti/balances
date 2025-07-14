<script setup>
import "chartjs-adapter-luxon"
import { Chart, registerables } from "chart.js"
import { nextTick, onActivated, onMounted, ref, watch } from "vue"
import { useAssetsStore } from "@/stores/assetsStore"

Chart.register(...registerables)

const assetsStore = useAssetsStore()

const canvas = ref(null)
let chartInstance = null

function renderChart() {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "line",
      data: {
        datasets: assetsStore.assetValueHistory,
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
          y: {
            stacked: true,
            min: 0,
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

watch(
  () => assetsStore.assetValueHistory,
  (newAssetHistory) => {
    if (newAssetHistory.length > 0) {
      destroyChart()
      renderChart()
    }
  }
)

onMounted(async () => {
  await nextTick()
  if (assetsStore.assetValueHistory.length > 0) {
    renderChart()
    destroyChart()
  }
})

onActivated(async () => {
  await nextTick()
  if (assetsStore.assetValueHistory.length > 0) {
    renderChart()
    destroyChart()
  }
})
</script>

<template>
  <div class="flex-1 h-full">
    <canvas ref="canvas"></canvas>
  </div>
</template>