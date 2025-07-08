<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"

import { usePropertiesStore } from "@/stores/propertiesStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/PropertyLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const propertiesStore = usePropertiesStore()
const defaultPageSize = ev.getDefaultPageSize()

const debouncedFilterProperties = debounce(() => {
  propertiesStore.filterProperties()
}, 300)

watch(
  [
    () => propertiesStore.lvFilter,
    () => propertiesStore.lvValuesStartDate,
    () => propertiesStore.lvValuesEndDate,
    () => propertiesStore.lvPageSize,
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
    debouncedFilterProperties()
  }
)

const showAddPropertyDialog = () => {
  propertiesStore.resetLVAddPropertyDialog()
  lvAddPropertyDialog.showModal()
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
  propertiesStore.lvValuesStartDate = parsedLVValuesStartDate

  const parsedLVValuesEndDate = numUtils.queryParamToNullableInt(
    query.lvValuesEndDate
  )
  propertiesStore.lvValuesEndDate = parsedLVValuesEndDate

  const defaultLVValuesStartDate = dateUtils.getEpochOneYearAgo()
  propertiesStore.lvHydrate(
    query.lvFilter?.toString() || "",
    parsedLVValuesStartDate &&
      parsedLVValuesStartDate !== defaultLVValuesStartDate
      ? parsedLVValuesStartDate
      : defaultLVValuesStartDate,
    parsedLVValuesEndDate,
    parsedLVPageSize
  )

  debouncedFilterProperties()
}

onMounted(() => refetch())
onUnmounted(() => propertiesStore.lvDehydrate())

const createProperty = async () => {
  const res = await propertiesStore.createProperty()
  if (!res.error) {
    lvAddPropertyDialog.close()
    propertiesStore.resetLVAddPropertyDialog()
  }
}

const showDeletePropertyConfirmation = async (propertyId) => {
  const res = await propertiesStore.getPropertyToDeleteById(propertyId)
  if (!res.error) {
    lvDeletePropertyDialog.showModal()
  }
}

const cancelDeleteProperty = () => {
  lvDeletePropertyDialog.close()
  propertiesStore.resetLVDeletePropertyDialog()
}

const deleteProperty = async () => {
  const res = await propertiesStore.deleteProperty()
  if (!res.error) {
    lvDeletePropertyDialog.close()
    propertiesStore.resetLVAddPropertyDialog()
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: List of Properties -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div class="flex gap-3">
            <h2 class="card-title">List of Properties</h2>
            <button
              class="btn btn-neutral btn-xs"
              v-on:click="showAddPropertyDialog()"
            >
              <font-awesome-icon :icon="['fas', 'plus']" />
              New Property
            </button>
          </div>
          <input
            type="text"
            placeholder="Search properties..."
            class="input input-bordered w-64"
          />
        </div>
        <div class="overflow-x-auto h-88">
          <table class="table table-zebra w-full table-pin-rows">
            <thead>
              <tr>
                <th>Identification</th>
                <th>Location</th>
                <th class="text-right">Initial Value</th>
                <th class="text-right">Current Value</th>
                <th>Type</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(property, index) in propertiesStore.lvProperties"
                :key="index"
                class="hover:bg-base-300"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ property.name }}</div>
                      <div class="text-sm opacity-50">
                        {{ property.titleHolder }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ property.address }}</div>
                      <div class="text-sm opacity-50">
                        Total: {{ property.totalArea }}
                        {{ property.areaUnit }} | Building:
                        {{ property.buildingArea }} {{ property.areaUnit }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(property.initialValue) }}
                      </div>
                      <div class="text-sm opacity-50">
                        at
                        {{
                          dateUtils.epochToLocalDate(property.initialValueDate)
                        }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(property.currentValue) }}
                      </div>
                      <div class="text-sm opacity-50">
                        at
                        {{
                          dateUtils.epochToLocalDate(property.currentValueDate)
                        }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">{{
                      property.type
                    }}</span>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">{{
                      property.status
                    }}</span>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <router-link
                      :to="{
                        name: 'assets.property.detail',
                        params: { id: property.id },
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
                      v-on:click="showDeletePropertyConfirmation(property.id)"
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

    <!-- Bottom Half: Line Chart of the Property's Values -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body flex-1 min-h-0">
        <h2 class="card-title">Value History Chart</h2>
        <line-chart class="flex-1 min-h-0" />
      </div>
    </div>
  </div>

  <!-- Dialog Box: Property Adder -->
  <dialog id="lvAddPropertyDialog" class="modal">
    <div class="modal-box w-11/12 max-w-5xl overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Add New Property</h3>
      <form class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Property Name*</label>
          <input
            v-model="propertiesStore.lvAddProperty.name"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Address*</label>
          <input
            v-model="propertiesStore.lvAddProperty.address"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Total Area*</label>
          <input
            v-model="propertiesStore.lvAddProperty.totalArea"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Building Area*</label>
          <input
            v-model="propertiesStore.lvAddProperty.buildingArea"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Area Unit*</label>
          <select
            v-model="propertiesStore.lvAddProperty.areaUnit"
            class="select select-bordered w-full"
          >
            <option>sqft</option>
            <option>sqm</option>
          </select>
        </div>
        <div>
          <label class="label">Type*</label>
          <select
            v-model="propertiesStore.lvAddProperty.type"
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
            v-model="propertiesStore.lvAddProperty.titleHolder"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Tax Identifier*</label>
          <input
            v-model="propertiesStore.lvAddProperty.taxIdentifier"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Purchase Date*</label>
          <DatePicker
            v-model:date="propertiesStore.lvAddProperty.purchaseDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Initial Value*</label>
          <input
            v-model="propertiesStore.lvAddProperty.initialValue"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Initial Value Date*</label>
          <DatePicker
            v-model:date="propertiesStore.lvAddProperty.initialValueDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Current Value*</label>
          <input
            v-model="propertiesStore.lvAddProperty.currentValue"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Current Value Date*</label>
          <DatePicker
            v-model:date="propertiesStore.lvAddProperty.currentValueDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div>
          <label class="label">Annual Appreciation Rate (%)*</label>
          <input
            v-model="propertiesStore.lvAddProperty.annualAppreciationPercent"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="createProperty" class="btn btn-primary">
            Save
          </button>
          <button
            type="button"
            @click="propertiesStore.resetLVAddPropertyDialog()"
            class="btn btn-secondary"
          >
            Reset
          </button>
        </div>
      </form>
    </div>
  </dialog>

  <!-- Dialog Box: Confirm Property Delete -->
  <dialog id="lvDeletePropertyDialog" class="modal">
    <div class="modal-box w-11/12 max-w-5xl overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Confirm Property Delete</h3>
      <form class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-2 grid-rows-7 gap-4">
          <div>Name</div>
          <div>
            {{ propertiesStore.lvDeleteProperty.name }}
          </div>
          <div>Address</div>
          <div>{{ propertiesStore.lvDeleteProperty.address }}</div>
          <div>Total Area</div>
          <div>
            {{ propertiesStore.lvDeleteProperty.totalArea }}
          </div>
          <div>Building Area</div>
          <div>
            {{ propertiesStore.lvDeleteProperty.buildingArea }}
          </div>
          <div>Area Unit</div>
          <div>
            {{ propertiesStore.lvDeleteProperty.areaUnit }}
          </div>
          <div>Type</div>
          <div>{{ propertiesStore.lvDeleteProperty.type }}</div>
          <div>Title Holder</div>
          <div>{{ propertiesStore.lvDeleteProperty.titleHolder }}</div>
          <div>Tax Identifier</div>
          <div>{{ propertiesStore.lvDeleteProperty.taxIdentifier }}</div>
          <div>Purchase Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                propertiesStore.lvDeleteProperty.purchaseDate
              )
            }}
          </div>
          <div>Initial Value</div>
          <div>
            {{
              numUtils.numericToMoney(
                propertiesStore.lvDeleteProperty.initialValue
              )
            }}
          </div>
          <div>Initial Value Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                propertiesStore.lvDeleteProperty.initialValueDate
              )
            }}
          </div>
          <div>Current Value</div>
          <div>
            {{
              numUtils.numericToMoney(
                propertiesStore.lvDeleteProperty.currentValue
              )
            }}
          </div>
          <div>Current Value Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                propertiesStore.lvDeleteProperty.currentValueDate
              )
            }}
          </div>
          <div>Annual Appreciation Rate (%)</div>
          <div>
            {{ propertiesStore.lvDeleteProperty.annualAppreciationPercent }}
          </div>
          <div>Status</div>
          <div>{{ propertiesStore.lvDeleteProperty.status }}</div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="deleteProperty()"
            class="btn btn-primary"
          >
            Confirm
          </button>
          <button
            type="button"
            @click="cancelDeleteProperty()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>