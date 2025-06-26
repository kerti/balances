<script setup>
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/VehicleLineChart.vue"

const numUtils = useNumUtils()

const showAddVehicleDialog = () => {
  // bankAccountsStore.resetLVAddBankAccountDialog()
  lvAddVehicleDialog.showModal()
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
                <th class="text-right">Purchase Value</th>
                <th class="text-right">Current Value</th>
                <th>Type</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">John's Car</div>
                      <div class="text-sm opacity-50">2008 Toyota Yaris E</div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">John Fitzgerald Doe</div>
                      <div class="text-sm opacity-50">AB 1234 RFT</div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(120000) }}
                      </div>
                      <div class="text-sm opacity-50">at 12 April 2020</div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(100000) }}
                      </div>
                      <div class="text-sm opacity-50">at 12 April 2025</div>
                    </div>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">CAR</span>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">IN USE</span>
                  </div>
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
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Make*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Model*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Year*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Type*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Title Holder*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">License Plate Number*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Purchase Date*</label>
          <DatePicker placeholder="pick a date" required />
        </div>
        <div>
          <label class="label">Initial Value*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Initial Value Date*</label>
          <DatePicker placeholder="pick a date" required />
        </div>
        <div>
          <label class="label">Current Value*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div>
          <label class="label">Current Value Date*</label>
          <DatePicker placeholder="pick a date" required />
        </div>
        <div>
          <label class="label">Annual Depreciation Rate (%)*</label>
          <input type="text" class="input input-bordered w-full" />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" class="btn btn-primary">Save</button>
          <button type="button" class="btn btn-secondary">Reset</button>
        </div>
      </form>
    </div>
  </dialog>
</template>