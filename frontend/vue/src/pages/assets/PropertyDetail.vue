<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"

import { usePropertiesStore } from "@/stores/propertiesStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/PropertyDetailLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const propertiesStore = usePropertiesStore()
const defaultPageSize = ev.getDefaultPageSize() * 31 // assume maxiumum of 31 values per month

const debouncedGet = debounce(() => {
  propertiesStore.getPropertyForDV()
}, 300)

// TODO: add controls to leverage this
watch(
  [
    () => propertiesStore.dvValuesStartDate,
    () => propertiesStore.dvValuesEndDate,
    () => propertiesStore.dvPageSize,
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
    if (propertiesStore.dvProperty.name !== undefined) {
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
  propertiesStore.dvValuesStartDate = parsedDVValueStartDate

  const parsedDVValueEndDate = numUtils.queryParamToNullableInt(
    query.dvValuesEndDate
  )
  propertiesStore.dvValuesEndDate = parsedDVValueEndDate

  const defaultDVValuesStartDate = dateUtils.getEpochOneYearAgo()

  propertiesStore.dvHydrate(
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
  propertiesStore.getPropertyForDV()
})

onUnmounted(() => propertiesStore.dvDehydrate())

const resetPropertyForm = () => {
  propertiesStore.revertDVPropertyToCache()
}

const saveProperty = () => {
  propertiesStore.updateProperty()
}

const showEditor = (valueId) => {
  if (valueId) {
    propertiesStore.dvValueEditorMode = "Edit"
    propertiesStore.getPropertyValueById(valueId)
  } else {
    propertiesStore.dvValueEditorMode = "Add"
    propertiesStore.prepDVBlankPropertyValue()
  }
  valueEditor.showModal()
}

const resetValueForm = () => {
  propertiesStore.revertDVPropertyValueToCache()
}

const saveValue = async () => {
  if (propertiesStore.dvValueEditorMode == "Edit") {
    const res = await propertiesStore.updatePropertyValue()
    if (!res.error) {
      valueEditor.close()
    }
  } else if (propertiesStore.dvValueEditorMode == "Add") {
    const res = await propertiesStore.createPropertyValue()
    if (!res.error) {
      valueEditor.close()
    }
  }
}

const showValueDeleteConfirmation = (valueId) => {
  propertiesStore.getPropertyValueById(valueId)
  vdConfirm.showModal()
}

const cancelValueDelete = () => {
  vdConfirm.close()
}

const deleteValue = async () => {
  const res = await propertiesStore.deletePropertyValue()
  if (!res.error) {
    vdConfirm.close()
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: Form and Values Table -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Left: Property Form -->
      <div class="card bg-base-100 shadow-md md:col-span-1">
        <div class="card-body">
          <h2 class="card-title">Property Details</h2>
          <form class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label class="label">Name*</label>
              <input
                v-model="propertiesStore.dvProperty.name"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Address*</label>
              <input
                v-model="propertiesStore.dvProperty.address"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Total Area*</label>
              <input
                v-model="propertiesStore.dvProperty.totalArea"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Building Area*</label>
              <input
                v-model="propertiesStore.dvProperty.buildingArea"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Area Unit*</label>
              <select
                v-model="propertiesStore.dvProperty.areaUnit"
                class="select select-bordered w-full"
              >
                <option>sqft</option>
                <option>sqm</option>
              </select>
            </div>
            <div>
              <label class="label">Type*</label>
              <select
                v-model="propertiesStore.dvProperty.type"
                class="select select-bordered w-full"
              >
                <option>land</option>
                <option>house</option>
                <option>apartment</option>
              </select>
            </div>
            <div>
              <label class="label">Title Holder*</label>
              <input
                v-model="propertiesStore.dvProperty.titleHolder"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Tax Identifier*</label>
              <input
                v-model="propertiesStore.dvProperty.taxIdentifier"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Purchase Date*</label>
              <DatePicker
                v-model:date="propertiesStore.dvProperty.purchaseDate"
                placeholder="pick a date"
                required
              />
            </div>
            <div>
              <label class="label">Initial Value*</label>
              <input
                v-model="propertiesStore.dvProperty.initialValue"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Initial Value Date*</label>
              <DatePicker
                v-model:date="propertiesStore.dvProperty.initialValueDate"
                placeholder="pick a date"
                required
              />
            </div>
            <div>
              <label class="label">Annual Appreciation (%)*</label>
              <input
                v-model="propertiesStore.dvProperty.annualAppreciationPercent"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status*</label>
              <select
                v-model="propertiesStore.dvProperty.status"
                class="select select-bordered w-full"
              >
                <option>in_use</option>
                <option>rented</option>
                <option>sold</option>
              </select>
            </div>
            <div class="flex justify-end gap-2 pt-4">
              <button
                type="button"
                @click="saveProperty"
                class="btn btn-primary"
              >
                Save
              </button>
              <button
                type="button"
                @click="resetPropertyForm"
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
                  v-for="(value, index) in propertiesStore.dvProperty.values"
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
                        v-on:click="showEditor(value.id)"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Delete"
                        v-on:click="showValueDeleteConfirmation(value.id)"
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
        <h2 class="card-title">Value History</h2>
        <line-chart />
      </div>
    </div>
  </div>

  <!-- Dialog Box: Value Editor -->
  <dialog id="valueEditor" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">
        {{ propertiesStore.dvValueEditorMode }} Property Value
      </h3>
      <form class="gri grid-cols-1 gap-4">
        <div>
          <label class="label">Value*</label>
          <input
            v-model="propertiesStore.dvEditPropertyValue.value"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Date*</label>
          <DatePicker
            v-model:date="propertiesStore.dvEditPropertyValue.date"
            placeholder="pick a date"
            required
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="saveValue" class="btn btn-primary">
            Save
          </button>
          <button
            type="button"
            @click="resetValueForm"
            class="btn btn-secondary"
          >
            Reset
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>

  <!-- Dialog Box: Confirm Value Delete -->
  <dialog id="vdConfirm" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Confirm Property Value Delete</h3>
      <form class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-2 grid-rows-2 gap-4">
          <div>Value</div>
          <div>
            {{
              numUtils.numericToMoney(propertiesStore.dvEditPropertyValue.value)
            }}
          </div>
          <div>Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                propertiesStore.dvEditPropertyValue.date
              )
            }}
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="deleteValue()" class="btn btn-primary">
            Confirm
          </button>
          <button
            type="button"
            @click="cancelValueDelete()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>