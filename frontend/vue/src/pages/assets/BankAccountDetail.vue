<script setup>
import { onMounted, watch } from "vue"
import LineChart from "@/components/assets/BankDetailLineChart.vue"
import { useRoute, useRouter } from "vue-router"
import { useNumUtils } from "@/composables/useNumUtils"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"
import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import debounce from "lodash.debounce"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const bankAccountsStore = useBankAccountsStore()
const defaultPageSize = ev.getDefaultPageSize() * 12 // assume about 10 balances per month

const debouncedGet = debounce(() => {
  bankAccountsStore.get()
}, 300)

// TODO: add controls to leverage this
watch(
  [
    () => bankAccountsStore.detailBalanceStartDate,
    () => bankAccountsStore.detailBalanceEndDate,
    () => bankAccountsStore.detailPageSize,
  ],
  ([newBalanceStartDate, newBalanceEndDate, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultBalanceStartDate = dateUtils.getEpochOneYearAgo()
    const balanceStartDateParam =
      Number.isInteger(newBalanceStartDate) &&
      newBalanceStartDate !== defaultBalanceStartDate
        ? newBalanceStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        balanceStartDate: balanceStartDateParam,
        balanceEndDate: newBalanceEndDate || undefined,
        pageSize: pageSizeParam,
      },
    })
    // prevent double-fetching on initial component mount
    if (bankAccountsStore.account.bankName !== undefined) {
      debouncedGet()
    }
  }
)

function refetch() {
  const query = route.query

  const parsedPageSize = numUtils.queryParamToInt(
    query.pageSize,
    defaultPageSize
  )

  const parsedBalanceStartDate = numUtils.queryParamToNullableInt(
    query.balanceStartDate
  )
  bankAccountsStore.balancesStartDate = parsedBalanceStartDate

  const parsedBalanceEndDate = numUtils.queryParamToNullableInt(
    query.balanceEndDate
  )
  bankAccountsStore.balancesEndDate = parsedBalanceEndDate

  const defaultBalanceStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.hydrateDetail(
    route.params.id,
    parsedBalanceStartDate && parsedBalanceStartDate !== defaultBalanceStartDate
      ? parsedBalanceStartDate
      : defaultBalanceStartDate,
    parsedBalanceEndDate,
    parsedPageSize
  )
}

onMounted(() => {
  refetch()
  bankAccountsStore.get()
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
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Left: Account Form -->
      <div class="card bg-base-100 shadow-md md:col-span-2">
        <div class="card-body">
          <!-- TODO: Use the account name instead -->
          <h2 class="card-title">Account Details</h2>
          <form class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="label">Account Name</label>
              <input
                v-model="bankAccountsStore.account.accountName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Bank Name</label>
              <input
                v-model="bankAccountsStore.account.bankName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Holder Name</label>
              <input
                v-model="bankAccountsStore.account.accountHolderName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Number</label>
              <input
                v-model="bankAccountsStore.account.accountNumber"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status</label>
              <select
                v-model="bankAccountsStore.account.status"
                class="select select-bordered w-full"
              >
                <option>active</option>
                <option>inactive</option>
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
                  <th>Date</th>
                  <th class="text-right">Balance</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(entry, index) in bankAccountsStore.account.balances"
                  :key="index"
                >
                  <td>{{ dateUtils.epochToLocalDate(entry.date) }}</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(entry.balance) }}
                  </td>
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
        <line-chart />
      </div>
    </div>
  </div>
</template>