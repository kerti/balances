<script setup>
import { useDateUtils } from "@/composables/useDateUtils"
import { useNumUtils } from "@/composables/useNumUtils"
import PieChart from "@/components/assets/AssetProportionsPieChart.vue"
import LineChart from "@/components/assets/AssetValuesStackedLineChart.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
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
                <tr>
                  <td>Shopping Account</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">cash</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(1230) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
                <tr>
                  <td>John's Car</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">vehicle</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(123000) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
                <tr>
                  <td>Jane's House</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">property</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(1850000) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
                <tr>
                  <td>Mike's Bike</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">vehicle</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(17430) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
                <tr>
                  <td>Margaret's Yacht</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">vehicle</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(2957000) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
                <tr>
                  <td>Margaret's Yacht</td>
                  <td>
                    <span class="badge badge-sm badge-neutral">vehicle</span>
                  </td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(2957000) }}
                  </td>
                  <td>
                    {{
                      dateUtils.epochToLocalDate(dateUtils.getEpochOneYearAgo())
                    }}
                  </td>
                </tr>
              </tbody>
              <tfoot class="bg-base-200 font-bold bottom-0 z-10">
                <tr>
                  <td colspan="2">Total Value of Assets</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(3492340) }}
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
            <pie-chart class="w-full h-full max-w-md max-h-md"></pie-chart>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Half: Value History of Assets -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body">
        <h2 class="card-title">Value History</h2>
        <line-chart />
      </div>
    </div>
  </div>
</template>