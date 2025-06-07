<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useNumUtils } from "@/composables/useNumUtils"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"
import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import debounce from "lodash.debounce"

import LineChart from "@/components/assets/BankDetailLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

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

onUnmounted(() => bankAccountsStore.dehydrateDetail())

const resetAccountForm = () => {
  bankAccountsStore.revertAccountToCache()
}

const saveAccount = () => {
  bankAccountsStore.update()
}

const showEditor = (balanceId) => {
  if (balanceId) {
    bankAccountsStore.balanceEditorMode = "Edit"
    bankAccountsStore.getBalanceById(balanceId)
  } else {
    bankAccountsStore.balanceEditorMode = "Add"
    bankAccountsStore.prepBlankBalance()
  }
  balanceEditor.showModal()
}

const resetBalanceForm = () => {
  bankAccountsStore.revertBalanceToCache()
}

const saveBalance = async () => {
  if (bankAccountsStore.balanceEditorMode == "Edit") {
    const res = await bankAccountsStore.updateBalance()
    if (!res.errorMessage) {
      balanceEditor.close()
    }
  } else if (bankAccountsStore.balanceEditorMode == "Add") {
    const res = await bankAccountsStore.createBalance()
    if (!res.errorMessage) {
      balanceEditor.close()
    }
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: Form and Balances Table -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Left: Account Form -->
      <div class="card bg-base-100 shadow-md md:col-span-1">
        <div class="card-body">
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
                @click="resetAccountForm"
                class="btn btn-secondary"
              >
                Reset
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right: Balance Table -->
      <div class="card bg-base-100 shadow-md md:col-span-1">
        <div class="card-body">
          <div class="flex items-center justify-between mb-4">
            <h2 class="card-title">Balance History</h2>
            <button
              class="btn btn-neutral btn-circle tooltip"
              data-tip="Add New Balance"
              v-on:click="showEditor()"
            >
              <font-awesome-icon :icon="['fas', 'plus']" />
            </button>
          </div>
          <div class="overflow-x-auto h-64">
            <table class="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Date</th>
                  <th class="text-right">Balance</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(entry, index) in bankAccountsStore.account.balances"
                  :key="index"
                >
                  <td>{{ dateUtils.epochToShortLocalDate(entry.date) }}</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(entry.balance) }}
                  </td>
                  <td>
                    <div class="flex items-center gap-3">
                      <button
                        class="btn btn-neutral tooltip"
                        data-tip="Edit"
                        v-on:click="showEditor(entry.id)"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                      <button class="btn btn-neutral tooltip" data-tip="Delete">
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
        <h2 class="card-title">Account Balance Over Time (last 12 months)</h2>
        <line-chart />
      </div>
    </div>
  </div>

  <!-- Dialog Box: Balance Editor -->
  <dialog id="balanceEditor" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold">
        {{ bankAccountsStore.balanceEditorMode }} Bank Account Balance
      </h3>
      <form class="grid grid-cols-1 gap-4">
        <div>
          <label class="label">Balance</label>
          <!-- <label class="input"> -->
          <input
            v-model="bankAccountsStore.beBalance.balance"
            type="text"
            class="input input-bordered w-full"
          />
          <!-- </label> -->
        </div>
        <div>
          <label class="label">Date</label>
          <!-- <div class="input input-bordered p-0"> -->
          <DatePicker
            v-model:date="bankAccountsStore.beBalance.date"
            placeholder="pick a date"
            required
          />
          <!-- </div> -->
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="saveBalance" class="btn btn-primary">
            Save
          </button>
          <button
            type="button"
            @click="resetBalanceForm"
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
</template>