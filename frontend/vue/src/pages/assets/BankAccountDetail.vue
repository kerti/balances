<script setup>
import { onMounted, onUnmounted, watch } from "vue"
import { useRoute, useRouter } from "vue-router"

import { useBankAccountsStore } from "@/stores/bankAccountsStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/BankDetailLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const bankAccountsStore = useBankAccountsStore()
const defaultPageSize = ev.getDefaultPageSize() * 31 // assume maximum of 31 balances per month

const debouncedGet = debounce(() => {
  bankAccountsStore.getBankAccountForDV()
}, 300)

// TODO: add controls to leverage this
watch(
  [
    () => bankAccountsStore.dvBalancesStartDate,
    () => bankAccountsStore.dvBalancesEndDate,
    () => bankAccountsStore.dvPageSize,
  ],
  ([newBalanceStartDate, newBalanceEndDate, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultDVBalanceStartDate = dateUtils.getEpochOneYearAgo()
    const dvBalancesStartDateParam =
      Number.isInteger(newBalanceStartDate) &&
      newBalanceStartDate !== defaultDVBalanceStartDate
        ? newBalanceStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        dvBalancesStartDate: dvBalancesStartDateParam,
        dvBalancesEndDate: newBalanceEndDate || undefined,
        pageSize: pageSizeParam,
      },
    })
    // prevent double-fetching on initial component mount
    if (bankAccountsStore.dvAccount.bankName !== undefined) {
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

  const parsedDVBalanceStartDate = numUtils.queryParamToNullableInt(
    query.dvBalancesStartDate
  )
  bankAccountsStore.dvBalancesStartDate = parsedDVBalanceStartDate

  const parsedDVBalanceEndDate = numUtils.queryParamToNullableInt(
    query.dvBalancesEndDate
  )
  bankAccountsStore.dvBalancesEndDate = parsedDVBalanceEndDate

  const defaultDVBalanceStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.dvHydrate(
    route.params.id,
    parsedDVBalanceStartDate &&
      parsedDVBalanceStartDate !== defaultDVBalanceStartDate
      ? parsedDVBalanceStartDate
      : defaultDVBalanceStartDate,
    parsedDVBalanceEndDate,
    parsedDVPageSize
  )
}

onMounted(() => {
  refetch()
  bankAccountsStore.getBankAccountForDV()
})

onUnmounted(() => bankAccountsStore.dvDehydrate())

const resetAccountForm = () => {
  bankAccountsStore.revertDVBankAccountToCache()
}

const saveAccount = () => {
  bankAccountsStore.updateBankAccount()
}

const showEditor = (balanceId) => {
  if (balanceId) {
    bankAccountsStore.dvBalanceEditorMode = "Edit"
    bankAccountsStore.getBankAccountBalanceById(balanceId)
  } else {
    bankAccountsStore.dvBalanceEditorMode = "Add"
    bankAccountsStore.prepDVBlankBankAccountBalance()
  }
  balanceEditor.showModal()
}

const resetBalanceForm = () => {
  bankAccountsStore.revertDVBankAccountBalanceToCache()
}

const saveBalance = async () => {
  if (bankAccountsStore.dvBalanceEditorMode == "Edit") {
    const res = await bankAccountsStore.updateBankAccountBalance()
    if (!res.error) {
      balanceEditor.close()
    }
  } else if (bankAccountsStore.dvBalanceEditorMode == "Add") {
    const res = await bankAccountsStore.createBankAccountBalance()
    if (!res.error) {
      balanceEditor.close()
    }
  }
}

const showBalanceDeleteConfirmaton = (balanceId) => {
  bankAccountsStore.getBankAccountBalanceById(balanceId)
  bdConfirm.showModal()
}

const cancelBalanceDelete = () => {
  bdConfirm.close()
}

const deleteBalance = async () => {
  const res = await bankAccountsStore.deleteBankAccountBalance()
  if (!res.error) {
    bdConfirm.close()
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
              <label class="label">Account Name*</label>
              <input
                v-model="bankAccountsStore.dvAccount.accountName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Bank Name*</label>
              <input
                v-model="bankAccountsStore.dvAccount.bankName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Holder Name*</label>
              <input
                v-model="bankAccountsStore.dvAccount.accountHolderName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Number*</label>
              <input
                v-model="bankAccountsStore.dvAccount.accountNumber"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status*</label>
              <select
                v-model="bankAccountsStore.dvAccount.status"
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
            <button class="btn btn-neutral btn-xs" v-on:click="showEditor()">
              <font-awesome-icon :icon="['fas', 'plus']" />
              Add New Balance
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
                  v-for="(balance, index) in bankAccountsStore.dvAccount
                    .balances"
                  :key="index"
                >
                  <td>{{ dateUtils.epochToShortLocalDate(balance.date) }}</td>
                  <td class="text-right">
                    {{ numUtils.numericToMoney(balance.balance) }}
                  </td>
                  <td>
                    <div class="flex items-center gap-3">
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Edit"
                        v-on:click="showEditor(balance.id)"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                      <button
                        class="btn btn-neutral btn-sm tooltip"
                        data-tip="Delete"
                        v-on:click="showBalanceDeleteConfirmaton(balance.id)"
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
        <h2 class="card-title">Balance History Chart</h2>
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
      <h3 class="text-lg font-bold pb-5">
        {{ bankAccountsStore.dvBalanceEditorMode }} Bank Account Balance
      </h3>
      <form class="grid grid-cols-1 gap-4">
        <div>
          <label class="label">Balance*</label>
          <input
            v-model="bankAccountsStore.dvEditBankAccountBalance.balance"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Date*</label>
          <DatePicker
            v-model:date="bankAccountsStore.dvEditBankAccountBalance.date"
            placeholder="pick a date"
            required
          />
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

  <!-- Dialog Box: Confirm Balance Delete -->
  <dialog id="bdConfirm" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Confirm Balance Delete</h3>
      <form class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-2 grid-rows-2 gap-4">
          <div>Balance</div>
          <div>
            {{
              numUtils.numericToMoney(
                bankAccountsStore.dvEditBankAccountBalance.balance
              )
            }}
          </div>
          <div>Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                bankAccountsStore.dvEditBankAccountBalance.date
              )
            }}
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="deleteBalance()"
            class="btn btn-primary"
          >
            Confirm
          </button>
          <button
            type="button"
            @click="cancelBalanceDelete()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>