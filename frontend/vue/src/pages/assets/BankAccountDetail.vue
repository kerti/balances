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
  bankAccountsStore.getBankAccountForDetailView()
}, 300)

// TODO: add controls to leverage this
watch(
  [
    () => bankAccountsStore.detailViewBalancesStartDate,
    () => bankAccountsStore.detailViewBalancesEndDate,
    () => bankAccountsStore.detailViewPageSize,
  ],
  ([newBalanceStartDate, newBalanceEndDate, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultDetailViewBalanceStartDate = dateUtils.getEpochOneYearAgo()
    const detailViewBalancesStartDateParam =
      Number.isInteger(newBalanceStartDate) &&
      newBalanceStartDate !== defaultDetailViewBalanceStartDate
        ? newBalanceStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        detailViewBalancesStartDate: detailViewBalancesStartDateParam,
        detailViewBalancesEndDate: newBalanceEndDate || undefined,
        pageSize: pageSizeParam,
      },
    })
    // prevent double-fetching on initial component mount
    if (bankAccountsStore.detailViewAccount.bankName !== undefined) {
      debouncedGet()
    }
  }
)

function refetch() {
  const query = route.query

  const parsedDetailViewPageSize = numUtils.queryParamToInt(
    query.pageSize,
    defaultPageSize
  )

  const parsedDetailViewBalanceStartDate = numUtils.queryParamToNullableInt(
    query.detailViewBalancesStartDate
  )
  bankAccountsStore.detailViewBalancesStartDate =
    parsedDetailViewBalanceStartDate

  const parsedDetailViewBalanceEndDate = numUtils.queryParamToNullableInt(
    query.detailViewBalancesEndDate
  )
  bankAccountsStore.detailViewBalancesEndDate = parsedDetailViewBalanceEndDate

  const defaultDetailViewBalanceStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.dvHydrate(
    route.params.id,
    parsedDetailViewBalanceStartDate &&
      parsedDetailViewBalanceStartDate !== defaultDetailViewBalanceStartDate
      ? parsedDetailViewBalanceStartDate
      : defaultDetailViewBalanceStartDate,
    parsedDetailViewBalanceEndDate,
    parsedDetailViewPageSize
  )
}

onMounted(() => {
  refetch()
  bankAccountsStore.getBankAccountForDetailView()
})

onUnmounted(() => bankAccountsStore.dvDehydrate())

const resetAccountForm = () => {
  bankAccountsStore.revertAccountToCache()
}

const saveAccount = () => {
  bankAccountsStore.updateBankAccount()
}

const showEditor = (balanceId) => {
  if (balanceId) {
    bankAccountsStore.detailViewBalanceEditorMode = "Edit"
    bankAccountsStore.getBankAccountBalanceById(balanceId)
  } else {
    bankAccountsStore.detailViewBalanceEditorMode = "Add"
    bankAccountsStore.prepBlankBalance()
  }
  balanceEditor.showModal()
}

const resetBalanceForm = () => {
  bankAccountsStore.revertBalanceToCache()
}

const saveBalance = async () => {
  if (bankAccountsStore.detailViewBalanceEditorMode == "Edit") {
    const res = await bankAccountsStore.updateBankAccountBalance()
    if (!res.errorMessage) {
      balanceEditor.close()
    }
  } else if (bankAccountsStore.detailViewBalanceEditorMode == "Add") {
    const res = await bankAccountsStore.createBankAccountBalance()
    if (!res.errorMessage) {
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
  if (!res.errorMessage) {
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
              <label class="label">Account Name</label>
              <input
                v-model="bankAccountsStore.detailViewAccount.accountName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Bank Name</label>
              <input
                v-model="bankAccountsStore.detailViewAccount.bankName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Holder Name</label>
              <input
                v-model="bankAccountsStore.detailViewAccount.accountHolderName"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Account Number</label>
              <input
                v-model="bankAccountsStore.detailViewAccount.accountNumber"
                type="text"
                class="input input-bordered w-full"
              />
            </div>
            <div>
              <label class="label">Status</label>
              <select
                v-model="bankAccountsStore.detailViewAccount.status"
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
                  v-for="(balance, index) in bankAccountsStore.detailViewAccount
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
                        class="btn btn-neutral tooltip"
                        data-tip="Edit"
                        v-on:click="showEditor(balance.id)"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </button>
                      <button
                        class="btn btn-neutral tooltip"
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
        {{ bankAccountsStore.detailViewBalanceEditorMode }} Bank Account Balance
      </h3>
      <form class="grid grid-cols-1 gap-4">
        <div>
          <label class="label">Balance</label>
          <input
            v-model="bankAccountsStore.detailViewEditBankAccountBalance.balance"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Date</label>
          <DatePicker
            v-model:date="
              bankAccountsStore.detailViewEditBankAccountBalance.date
            "
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
                bankAccountsStore.detailViewEditBankAccountBalance.balance
              )
            }}
          </div>
          <div>Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                bankAccountsStore.detailViewEditBankAccountBalance.date
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