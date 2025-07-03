<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"

import { useVehiclesStore } from "@/stores/vehiclesStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/VehicleDetailLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const vehiclesStore = useVehiclesStore()
const defaultPageSize = ev.getDefaultPageSize() * 31 // assume maxiumum of 31 values per month

const debouncedGet = debounce(() => {
  vehiclesStore.getVehicleForDV()
}, 300)

// TODO: add controls to leverage this
watch(
  [
    () => vehiclesStore.dvValuesStartDate,
    () => vehiclesStore.dvValuesEndDate,
    () => vehiclesStore.dvPageSize,
  ],
  ([newValuesStartDate, newValuesEndDate, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultDVValuesStartDate = dateUtils.getEpochOneYearAgo()
    const dvValuesStartDateParam =
      Number.isInteger(newValuesStartDate) &&
      newValuesStartDate !== defaultDVValuesStartDate
        ? newValuesStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        dvValuesStartDate: dvValuesStartDateParam,
        dvValuesEndDate: newValuesEndDate || undefined,
        pageSize: pageSizeParam,
      },
    })
    // prevent double-fetching on initial component mount
    if (vehiclesStore.dvVehicle.name !== undefined) {
      debouncedGet()
    }
  }
)

function refetch() {
  const query = route.query

  const parsedDVPageSize = numUtils.queryParamToInt(
    query.pageSize,
    defaultPageSize
  )

  const parsedDVValueStartDate = numUtils.queryParamToNullableInt(
    query.dvValuesStartDate
  )
  vehiclesStore.dvValuesStartDate = parsedDVValueStartDate

  const parsedDVValueEndDate = numUtils.queryParamToNullableInt(
    query.dvValuesEndDate
  )
  vehiclesStore.dvValuesEndDate = parsedDVValueEndDate

  const defaultDVValuesStartDate = dateUtils.getEpochOneYearAgo()

  vehiclesStore.dvHydrate(
    route.params.id,
    parsedDVValueStartDate &&
      parsedDVValueStartDate !== defaultDVValuesStartDate
      ? parsedDVValueStartDate
      : defaultDVValuesStartDate,
    parsedDVValueEndDate,
    parsedDVPageSize
  )
}

onMounted(() => {
  refetch()
  vehiclesStore.getVehicleForDV()
})

onUnmounted(() => vehiclesStore.dvDehydrate())

const resetVehicleForm = () => {
  vehiclesStore.revertDVVehicleToCache()
}

const saveVehicle = () => {
  vehiclesStore.updateVehicle()
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: Form and Values Table -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Left: Vehicle Form -->
      <div class="card bg-base-100 shadow-md md:col-span-1">
        <div class="card-body">
          <h2 class="card-title">Vehicle Details</h2>
          <form class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="label">Name*</label>
              <input
                v-model="vehiclesStore.dvVehicle.name"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Make*</label>
              <input
                v-model="vehiclesStore.dvVehicle.make"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Model*</label>
              <input
                v-model="vehiclesStore.dvVehicle.model"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Year*</label>
              <input
                v-model="vehiclesStore.dvVehicle.year"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Type*</label>
              <select
                v-model="vehiclesStore.dvVehicle.type"
                class="select select-bordered w-full"
              >
                <option>bicycle</option>
                <option>car</option>
                <option>motorcycle</option>
                <option>truck</option>
              </select>
            </div>
            <div>
              <label class="label">Title Holder*</label>
              <input
                v-model="vehiclesStore.dvVehicle.titleHolder"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">License Plate Number*</label>
              <input
                v-model="vehiclesStore.dvVehicle.licensePlateNumber"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Purchase Date*</label>
              <DatePicker
                v-model:date="vehiclesStore.dvVehicle.purchaseDate"
                placeholder="pick a date"
                required
              />
            </div>
            <div>
              <label class="label">Initial Value*</label>
              <input
                v-model="vehiclesStore.dvVehicle.initialValue"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Initial Value Date*</label>
              <DatePicker
                v-model:date="vehiclesStore.dvVehicle.initialValueDate"
                placeholder="pick a date"
                required
              />
            </div>
            <div>
              <label class="label">Annual Depreciation (%)*</label>
              <input
                v-model="vehiclesStore.dvVehicle.annualDepreciationPercent"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status*</label>
              <select
                v-model="vehiclesStore.dvVehicle.status"
                class="select select-bordered w-full"
              >
                <option>in_use</option>
                <option>sold</option>
                <option>retired</option>
              </select>
            </div>
            <div class="flex justify-end gap-2 pt-4">
              <button
                type="button"
                @click="saveVehicle"
                class="btn btn-primary"
              >
                Save
              </button>
              <button
                type="button"
                @click="resetVehicleForm"
                class="btn btn-secondary"
              >
                Reset
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right: Value Table -->
      <div class="card bg-base-100 shadow-md md:col-span-1">
        <div class="card-body">
          <div class="flex items-center justify-between mb-4">
            <h2 class="card-title">Value History</h2>
            <button class="btn btn-neutral btn-xs" v-on:click="showEditor()">
              <font-awesome-icon :icon="['fas', 'plus']" />
              Add New Value
            </button>
          </div>
          <div class="overflow-x-auto h-82">
            <table class="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Date</th>
                  <th class="text-right">Value</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(value, index) in vehiclesStore.dvVehicle.values"
                  :key="index"
                >
                  <td>{{ dateUtils.epochToShortLocalDate(value.date) }}</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(value.value) }}
                  </td>
                  <td>
                    <div class="flex items-center gap-3">
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Edit"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Delete"
                      >
                        <font-awesome-icon :icon="['fas', 'trash']" />
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Half: Line Chart -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body">
        <h2 class="card-title">Value History Chart</h2>
        <line-chart />
      </div>
    </div>
  </div>
</template>