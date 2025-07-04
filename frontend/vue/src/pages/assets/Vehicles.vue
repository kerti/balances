<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"

import { useVehiclesStore } from "@/stores/vehiclesStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/VehicleLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const vehiclesStore = useVehiclesStore()
const defaultPageSize = ev.getDefaultPageSize()

const debouncedFilterVehicles = debounce(() => {
  vehiclesStore.filterVehicles()
}, 300)

watch(
  [
    () => vehiclesStore.lvFilter,
    () => vehiclesStore.lvValuesStartDate,
    () => vehiclesStore.lvValuesEndDate,
    () => vehiclesStore.lvPageSize,
  ],
  ([newLVFilter, newLVValuesStartDate, newLVValuesEndDate, newLVPageSize]) => {
    const lvPageSizeParam =
      Number.isInteger(newLVPageSize) && newLVPageSize !== defaultPageSize
        ? newLVPageSize
        : undefined

    const defaultLVValuesStartDate = dateUtils.getEpochOneYearAgo()
    const lvValuesStartDateParam =
      Number.isInteger(newLVValuesStartDate) &&
      newLVValuesStartDate !== defaultLVValuesStartDate
        ? newLVValuesStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        lvFilter: newLVFilter || undefined,
        lvValuesStartDate: lvValuesStartDateParam,
        lvValuesEndDate: newLVValuesEndDate || undefined,
        lvPageSize: lvPageSizeParam,
      },
    })
    debouncedFilterVehicles()
  }
)

const showAddVehicleDialog = () => {
  vehiclesStore.resetLVAddVehicleDialog()
  lvAddVehicleDialog.showModal()
}

function refetch() {
  const query = route.query

  const parsedLVPageSize = numUtils.queryParamToInt(
    query.lvPageSize,
    defaultPageSize
  )

  const parsedLVValuesStartDate = numUtils.queryParamToNullableInt(
    query.lvValuesStartDate
  )
  vehiclesStore.lvValuesStartDate = parsedLVValuesStartDate

  const parsedLVValuesEndDate = numUtils.queryParamToNullableInt(
    query.lvValuesEndDate
  )
  vehiclesStore.lvValuesEndDate = parsedLVValuesEndDate

  const defaultLVValuesStartDate = dateUtils.getEpochOneYearAgo()
  vehiclesStore.lvHydrate(
    query.lvFilter?.toString() || "",
    parsedLVValuesStartDate &&
      parsedLVValuesStartDate !== defaultLVValuesStartDate
      ? parsedLVValuesStartDate
      : defaultLVValuesStartDate,
    parsedLVValuesEndDate,
    parsedLVPageSize
  )

  debouncedFilterVehicles()
}

onMounted(() => refetch())
onUnmounted(() => vehiclesStore.lvDehydrate())

const createVehicle = async () => {
  const res = await vehiclesStore.createVehicle()
  if (!res.error) {
    lvAddVehicleDialog.close()
    vehiclesStore.resetLVAddVehicleDialog()
  }
}

const showDeleteVehicleConfirmation = async (vehicleId) => {
  const res = await vehiclesStore.getVehicleToDeleteById(vehicleId)
  if (!res.error) {
    lvDeleteVehicleDialog.showModal()
  }
}

const cancelDeleteVehicle = () => {
  lvDeleteVehicleDialog.close()
  vehiclesStore.resetLVDeleteVehicleDialog()
}

const deleteVehicle = async () => {
  const res = await vehiclesStore.deleteVehicle()
  if (!res.error) {
    lvDeleteVehicleDialog.close()
    vehiclesStore.resetLVAddVehicleDialog()
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: List of Vehicles -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div class="flex gap-3">
            <h2 class="card-title">List of Vehicles</h2>
            <button
              class="btn btn-neutral btn-xs"
              v-on:click="showAddVehicleDialog()"
            >
              <font-awesome-icon :icon="['fas', 'plus']" />
              New Vehicle
            </button>
          </div>
          <input
            type="text"
            placeholder="Search vehicles..."
            class="input input-bordered w-64"
          />
        </div>
        <div class="overflow-x-auto h-88">
          <table class="table table-zebra w-full table-pin-rows">
            <thead>
              <tr>
                <th>Name</th>
                <th>Identification</th>
                <th class="text-right">Initial Value</th>
                <th class="text-right">Current Value</th>
                <th>Type</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(vehicle, index) in vehiclesStore.lvVehicles"
                :key="index"
                class="hover:bg-base-300"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ vehicle.name }}</div>
                      <div class="text-sm opacity-50">
                        {{ vehicle.year }} {{ vehicle.make }}
                        {{ vehicle.model }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ vehicle.titleHolder }}</div>
                      <div class="text-sm opacity-50">
                        {{ vehicle.licensePlateNumber }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(vehicle.initialValue) }}
                      </div>
                      <div class="text-sm opacity-50">
                        at
                        {{
                          dateUtils.epochToLocalDate(vehicle.initialValueDate)
                        }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(vehicle.currentValue) }}
                      </div>
                      <div class="text-sm opacity-50">
                        at
                        {{
                          dateUtils.epochToLocalDate(vehicle.currentValueDate)
                        }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">{{
                      vehicle.type
                    }}</span>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">{{
                      vehicle.status
                    }}</span>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <router-link
                      :to="{
                        name: 'assets.vehicle.detail',
                        params: { id: vehicle.id },
                      }"
                    >
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Edit"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                    </router-link>
                    <button
                      class="btn btn-neutral btn-sm tooltip"
                      data-tip="Delete"
                      v-on:click="showDeleteVehicleConfirmation(vehicle.id)"
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

    <!-- Bottom Half: Line Chart of the Vehicle's Values -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body flex-1 min-h-0">
        <h2 class="card-title">Value History Chart</h2>
        <line-chart class="flex-1 min-h-0" />
      </div>
    </div>
  </div>

  <!-- Dialog Box: Vehicle Adder -->
  <dialog id="lvAddVehicleDialog" class="modal">
    <div class="modal-box w-11/12 max-w-5xl overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Add New Vehicle</h3>
      <form class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Vehicle Name*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.name"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Make*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.make"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Model*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.model"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Year*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.year"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Type*</label>
          <select
            v-model="vehiclesStore.lvAddVehicle.type"
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
            v-model="vehiclesStore.lvAddVehicle.titleHolder"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">License Plate Number*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.licensePlateNumber"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Purchase Date*</label>
          <DatePicker
            v-model:date="vehiclesStore.lvAddVehicle.purchaseDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Initial Value*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.initialValue"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Initial Value Date*</label>
          <DatePicker
            v-model:date="vehiclesStore.lvAddVehicle.initialValueDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Current Value*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.currentValue"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Current Value Date*</label>
          <DatePicker
            v-model:date="vehiclesStore.lvAddVehicle.currentValueDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Annual Depreciation Rate (%)*</label>
          <input
            v-model="vehiclesStore.lvAddVehicle.annualDepreciationPercent"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="createVehicle" class="btn btn-primary">
            Save
          </button>
          <button
            type="button"
            @click="vehiclesStore.resetLVAddVehicleDialog()"
            class="btn btn-secondary"
          >
            Reset
          </button>
        </div>
      </form>
    </div>
  </dialog>

  <!-- Dialog Box: Confirm Vehicle Delete -->
  <dialog id="lvDeleteVehicleDialog" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Confirm Vehicle Delete</h3>
      <form class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-2 grid-rows-7 gap-4">
          <div>Name</div>
          <div>
            {{ vehiclesStore.lvDeleteVehicle.name }}
          </div>
          <div>Make</div>
          <div>{{ vehiclesStore.lvDeleteVehicle.make }}</div>
          <div>Model</div>
          <div>
            {{ vehiclesStore.lvDeleteVehicle.model }}
          </div>
          <div>Year</div>
          <div>
            {{ vehiclesStore.lvDeleteVehicle.year }}
          </div>
          <div>Type</div>
          <div>{{ vehiclesStore.lvDeleteVehicle.type }}</div>
          <div>Title Holder</div>
          <div>{{ vehiclesStore.lvDeleteVehicle.titleHolder }}</div>
          <div>License Plate Number</div>
          <div>{{ vehiclesStore.lvDeleteVehicle.licensePlateNumber }}</div>
          <div>Purchase Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                vehiclesStore.lvDeleteVehicle.purchaseDate
              )
            }}
          </div>
          <div>Initial Value</div>
          <div>
            {{
              numUtils.numericToMoney(
                vehiclesStore.lvDeleteVehicle.initialValue
              )
            }}
          </div>
          <div>Initial Value Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                vehiclesStore.lvDeleteVehicle.initialValueDate
              )
            }}
          </div>
          <div>Current Value</div>
          <div>
            {{
              numUtils.numericToMoney(
                vehiclesStore.lvDeleteVehicle.currentValue
              )
            }}
          </div>
          <div>Current Value Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                vehiclesStore.lvDeleteVehicle.currentValueDate
              )
            }}
          </div>
          <div>Annual Depreciation Rate (%)</div>
          <div>
            {{ vehiclesStore.lvDeleteVehicle.annualDepreciationPercent }}
          </div>
          <div>Status</div>
          <div>{{ vehiclesStore.lvDeleteVehicle.status }}</div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="deleteVehicle()"
            class="btn btn-primary"
          >
            Confirm
          </button>
          <button
            type="button"
            @click="cancelDeleteVehicle()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>