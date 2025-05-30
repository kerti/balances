<script setup>
import LineChart from "@/components/assets/BankLineChart.vue"
import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"
import debounce from "lodash.debounce"
import { ref, watch, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"

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

// use actual backend
const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const bankAccountsStore = useBankAccountsStore()

const defaultPageSize = ev.getDefaultPageSize()

const debouncedSearch = debounce(() => {
  bankAccountsStore.search(bankAccountsStore.filter, bankAccountsStore.pageSize)
}, 300)

watch(
  () => bankAccountsStore.filter,
  () => {
    debouncedSearch()
  }
)

watch(
  [() => bankAccountsStore.filter, () => bankAccountsStore.pageSize],
  ([newFilter, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    router.replace({
      query: {
        ...route.query,
        filter: newFilter || undefined,
        pageSize: pageSizeParam,
      },
    })
  }
)

onMounted(() => {
  const query = route.query

  bankAccountsStore.filter = query.filter?.toString() || ""

  const parsedPageSize = numUtils.queryParamToInt(
    query.pageSize,
    defaultPageSize
  )
  bankAccountsStore.pageSize =
    parsedPageSize !== defaultPageSize ? parsedPageSize : defaultPageSize

  bankAccountsStore.hydrate(query.filter?.toString() || "", parsedPageSize)
})
</script>

<template>
  <div class="space-y-6">
    <!-- Top Half: List of Accounts -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <h2 class="card-title">List of Accounts</h2>
          <div class="form-control">
            <input
              type="text"
              v-model="bankAccountsStore.filter"
              placeholder="Search accounts..."
              class="input input-bordered w-64"
            />
          </div>
        </div>
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
                v-for="(account, index) in bankAccountsStore.accounts"
                :key="index"
                class="hover:bg-base-300"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ account.accountName }}</div>
                      <div class="text-sm opacity-50">
                        {{ account.accountHolderName }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div class="font-bold">{{ account.bankName }}</div>
                      <div class="text-sm opacity-50">
                        # {{ account.accountNumber }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div class="font-bold">
                        {{ numUtils.numericToMoney(account.lastBalance) }}
                      </div>
                      <div class="text-sm opacity-50">
                        at
                        {{
                          dateUtils.epochToLocalDate(account.lastBalanceDate)
                        }}
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