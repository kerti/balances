<script setup>
import "chartjs-adapter-luxon"
import { Chart, registerables } from "chart.js"
import { computed, nextTick, onActivated, onMounted, ref, watch } from "vue"
import { useAssetsStore } from "@/stores/assetsStore"

Chart.register(...registerables)

const assetsStore = useAssetsStore()

const sumOfCash = computed(() => {
  return assetsStore.assets.reduce((sum, val) => {
    if (val.class === "cash") {
      return sum + val.value
    } else {
      return sum
    }
  }, 0)
})

const sumOfVehicles = computed(() => {
  return assetsStore.assets.reduce((sum, val) => {
    if (val.class === "vehicle") {
      return sum + val.value
    } else {
      return sum
    }
  }, 0)
})

const sumOfProperties = computed(() => {
  return assetsStore.assets.reduce((sum, val) => {
    if (val.class === "property") {
      return sum + val.value
    } else {
      return sum
    }
  }, 0)
})

const canvas = ref(null)
let chartInstance = null

function renderChart() {
  if (canvas.value) {
    chartInstance = new Chart(canvas.value, {
      type: "doughnut",
      data: {
        labels: ["Cash", "Vehicles", "Properties"],
        datasets: [{ data: [sumOfCash, sumOfVehicles, sumOfProperties] }],
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

watch(
  () => assetsStore.assets,
  (newAssets) => {
    if (newAssets.length > 0) {
      destroyChart()
      renderChart()
    }
  }
)

onMounted(async () => {
  await nextTick()
  if (assetsStore.assets.length > 0) {
    destroyChart()
    renderChart()
  }
})

onActivated(async () => {
  await nextTick()
  if (assetsStore.assets.length > 0) {
    destroyChart()
    renderChart()
  }
})
</script>

<template>
  <div class="flex justify-center items-center w-full h-full">
    <canvas ref="canvas" class="max-w-full max-h-full"></canvas>
  </div>
</template>