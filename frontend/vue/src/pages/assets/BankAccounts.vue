<script setup>
import LineChart from "@/components/assets/BankLineChart.vue"
import { ref } from "vue"

const bankAccounts = ref([
  {
    name: "John's Main Account",
    holder: "John Fitzgerald Doe",
    bank: "First Bank of John",
    number: "1234567890",
    lastDate: "2025-05-18",
    balance: "Rp1,200,000.00",
    status: "active",
  },
  {
    name: "John's Savings",
    holder: "John Fitzgerald Doe",
    bank: "Second Bank of John",
    number: "0987654321",
    lastDate: "2025-05-18",
    balance: "Rp30,000,000.00",
    status: "active",
  },
  {
    name: "Jane's Main Account",
    holder: "Jane Montgomery Doe",
    bank: "First Bank of Jane",
    number: "1357924680",
    lastDate: "2025-05-18",
    balance: "Rp1,250,000.00",
    status: "inactive",
  },
  {
    name: "Jane's Savings",
    holder: "Jane Montgomery Doe",
    bank: "Second Bank of Jane",
    number: "0864297531",
    lastDate: "2025-05-18",
    balance: "Rp27,500,000.00",
    status: "active",
  },
  {
    name: "Jack's Main Account",
    holder: "Jack Reacher Doe",
    bank: "First Bank of Jack",
    number: "1470258369",
    lastDate: "2025-05-18",
    balance: "Rp799,000.00",
    status: "active",
  },
  {
    name: "Jack's Savings",
    holder: "Jack Reacher Doe",
    bank: "Second Bank of Jack",
    number: "0741963852",
    lastDate: "2025-05-18",
    balance: "Rp14,764,000.00",
    status: "inactive",
  },
])

const chartData = ref({
  labels: ["Jan", "Feb", "Mar", "Apr", "May"],
  datasets: [
    {
      label: "John's Main Account",
      data: [10000000, 11000000, 11500000, 12000000, 1200000],
    },
    {
      label: "John's Savings",
      data: [25000000, 26000000, 27500000, 29500000, 30000000],
    },
    {
      label: "Jane's Main Account",
      data: [10530000, 15000000, 7500000, 12780000, 1250000],
    },
    {
      label: "Jane's Savings",
      data: [25753200, 22000000, 22750000, 31477000, 27500000],
    },
    {
      label: "Jack's Main Account",
      data: [11530000, 10000000, 9500000, 8780000, 799000],
    },
    {
      label: "Jack's Savings",
      data: [22753200, 17000000, 14750000, 13477000, 14764000],
    },
  ],
})
</script>

<template>
  <div class="space-y-6">
    <!-- Top Half: List of Accounts -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <h2 class="card-title">List of Accounts</h2>
        <div class="overflow-x-auto h-96">
          <table class="table table-zebra w-full table-pin-rows">
            <thead>
              <tr>
                <th>Account</th>
                <th>Bank</th>
                <th class="text-right">Balance</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(account, index) in bankAccounts"
                :key="index"
                class="hover:bg-base-300"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ account.name }}</div>
                      <div class="text-sm opacity-50">
                        {{ account.holder }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ account.bank }}</div>
                      <div class="text-sm opacity-50">
                        # {{ account.number }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">{{ account.balance }}</div>
                      <div class="text-sm opacity-50">
                        at {{ account.lastDate }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div>
                    <span class="badge badge-sm badge-neutral">{{
                      account.status
                    }}</span>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <button class="btn btn-primary tooltip" data-tip="Edit">
                      <font-awesome-icon :icon="['fas', 'edit']" />
                    </button>
                    <button class="btn btn-primary tooltip" data-tip="Activate">
                      <font-awesome-icon :icon="['fas', 'eye']" />
                    </button>
                    <button
                      class="btn btn-primary tooltip"
                      data-tip="Deactivate"
                    >
                      <font-awesome-icon :icon="['fas', 'eye-slash']" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Bottom Half: Line Chart of the Accounts' Balances -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <h2 class="card-title">Balance Over Time</h2>
        <line-chart :chart-data="chartData" />
      </div>
    </div>
  </div>
</template>