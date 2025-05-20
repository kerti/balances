<script setup>
import { ref } from "vue"
import LineChart from "@/components/assets/BankDetailLineChart.vue"

// TODO: fetch from real API
const account = ref({
  bank: "John's Main Account",
  holder: "John Fitzgerald Doe",
  number: "1234567890",
  status: "Active",
})

const balances = ref([
  { amount: "Rp21,000,000", date: "2024-05-01" },
  { amount: "Rp17,000,000", date: "2024-06-01" },
  { amount: "Rp14,000,000", date: "2024-07-01" },
  { amount: "Rp18,000,000", date: "2024-08-01" },
  { amount: "Rp12,000,000", date: "2024-09-01" },
  { amount: "Rp9,000,000", date: "2024-10-01" },
  { amount: "Rp14,000,000", date: "2024-11-01" },
  { amount: "Rp7,000,000", date: "2024-12-01" },
  { amount: "Rp10,000,000", date: "2025-01-01" },
  { amount: "Rp11,500,000", date: "2025-02-01" },
  { amount: "Rp12,000,000", date: "2025-03-01" },
  { amount: "Rp12,000,000", date: "2025-04-01" },
])

const chartData = ref({
  labels: balances.value.map((entry) => entry.date),
  datasets: [
    {
      label: "Account Balance",
      data: balances.value.map((entry) =>
        parseInt(entry.amount.replace(/Rp|,|\./g, ""))
      ),
      fill: true,
    },
  ],
})

const resetForm = () => {
  // TODO: use actual data to reset the form
  account.value = {
    bank: "John's Main Account",
    holder: "John Fitzgerald Doe",
    number: "1234567890",
    status: "Active",
  }
}

const saveAccount = () => {
  console.log("Saved account:", account.value)
  // TODO: hook this to an API call or store logic
}
</script>

<template>
  <div class="space-y-6">
    <!-- Top Half: Form and Balances Table -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Left: Account Form -->
      <div class="card bg-base-100 shadow-md">
        <div class="card-body">
          <!-- TODO: Use the account name instead -->
          <h2 class="card-title">Account Details: {{ $route.params.id }}</h2>
          <form class="space-y-4">
            <div>
              <label class="label">Bank Name</label>
              <input
                v-model="account.bank"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Holder Name</label>
              <input
                v-model="account.holder"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Number</label>
              <input
                v-model="account.number"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status</label>
              <select
                v-model="account.status"
                class="select select-bordered w-full"
              >
                <option>Active</option>
                <option>Inactive</option>
              </select>
            </div>
            <div class="flex justify-end gap-2 pt-4">
              <button
                type="button"
                @click="saveAccount"
                class="btn btn-primary"
              >
                Save
              </button>
              <button
                type="button"
                @click="resetForm"
                class="btn btn-secondary"
              >
                Reset
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right: Balance Table -->
      <div class="card bg-base-100 shadow-md">
        <div class="card-body">
          <h2 class="card-title">Balances</h2>
          <div class="overflow-x-auto h-96">
            <table class="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Balance</th>
                  <th>Date</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(entry, index) in balances" :key="index">
                  <td>{{ entry.amount }}</td>
                  <td>{{ entry.date }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Half: Line Chart -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <h2 class="card-title">Account Balance Over Time (last 12 months)</h2>
        <line-chart :chart-data="chartData" />
      </div>
    </div>
  </div>
</template>