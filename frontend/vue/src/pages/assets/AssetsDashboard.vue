<script setup>
import { useDateUtils } from "@/composables/useDateUtils"
import { useNumUtils } from "@/composables/useNumUtils"
import AssetProportionsPieChart from "@/components/assets/AssetProportionsPieChart.vue"
import AssetValuesLineChart from "@/components/assets/AssetValuesStackedLineChart.vue"
import { computed, onMounted, onUnmounted } from "vue"
import { useAssetsStore } from "@/stores/assetsStore"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const assetsStore = useAssetsStore()

const totalAssetValue = computed(() => {
  return assetsStore.assets.reduce((sum, val) => {
    return sum + val.value
  }, 0)
})

onMounted(() => {
  assetsStore.hydrate()
})

onUnmounted(() => {
  assetsStore.dehydrate()
})
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: Asset List and Pie Chart -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Left: Asset List -->
      <div class="card bg-base-100 shadow-md md:col-span-2">
        <div class="card-body">
          <h2 class="card-title">List of Assets</h2>
          <div class="overflow-auto h-92">
            <table class="table table-zebra w-full table-pin-rows">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Class</th>
                  <th class="text-right">Value</th>
                  <th>Last Updated</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(asset, index) in assetsStore.assets" :key="index">
                  <td>{{ asset.name }}</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">{{
                      asset.class
                    }}</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(asset.value) }}
                  </td>
                  <td>
                    {{ dateUtils.epochToLocalDate(asset.lastUpdated) }}
                  </td>
                </tr>
              </tbody>
              <tfoot class="bg-base-200 font-bold bottom-0 z-10">
                <tr>
                  <td colspan="2">Total Value of Assets</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(totalAssetValue) }}
                  </td>
                  <td></td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>

      <!-- Right: Asset Pie Chart -->
      <div class="card bg-base-100 shadow-md">
        <div class="card-body">
          <h2 class="card-title">Asset Proportions</h2>
          <div class="flex justify-center items-center h-92">
            <AssetProportionsPieChart
              class="w-full h-full max-w-md max-h-md"
            ></AssetProportionsPieChart>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Half: Value History of Assets -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body">
        <h2 class="card-title">Value History</h2>
        <AssetValuesLineChart />
      </div>
    </div>
  </div>
</template>