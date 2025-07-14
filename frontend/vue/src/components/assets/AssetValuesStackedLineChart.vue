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
      type: "line",
      data: {
        datasets: [
          {
            label: "Cash",
            fill: true,
            data: [
              { x: 1722445199000, y: 200000000 },
              { x: 1725123599000, y: 220000000 },
              { x: 1727715599000, y: 230000000 },
            ],
          },
          {
            label: "Vehicles",
            fill: true,
            data: [
              { x: 1722445199000, y: 400000000 },
              { x: 1725123599000, y: 380000000 },
              { x: 1727715599000, y: 365000000 },
            ],
          },
          {
            label: "Properties",
            fill: true,
            data: [
              { x: 1722445199000, y: 1200000000 },
              { x: 1725123599000, y: 1250000000 },
              { x: 1727715599000, y: 1320000000 },
            ],
          },
        ],
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
  <div class="flex-1 h-full">
    <canvas ref="canvas"></canvas>
  </div>
</template>