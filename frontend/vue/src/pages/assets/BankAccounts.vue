<script setup>
import { watch, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"

import { useBankAccountsStore } from "@/stores/bankAccountsStore"

import debounce from "lodash.debounce"

import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"

import LineChart from "@/components/assets/BankLineChart.vue"
import DatePicker from "@/components/common/DatePicker.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const bankAccountsStore = useBankAccountsStore()
const defaultPageSize = ev.getDefaultPageSize()

const debouncedFilterBankAccounts = debounce(() => {
  bankAccountsStore.filterBankAccounts()
}, 300)

watch(
  [
    () => bankAccountsStore.lvFilter,
    () => bankAccountsStore.lvBalancesStartDate,
    () => bankAccountsStore.lvBalancesEndDate,
    () => bankAccountsStore.lvPageSize,
  ],
  ([
    newLVFilter,
    newLVBalancesStartDate,
    newLVBalancesEndDate,
    newLVPageSize,
  ]) => {
    const lvPageSizeParam =
      Number.isInteger(newLVPageSize) && newLVPageSize !== defaultPageSize
        ? newLVPageSize
        : undefined

    const defaultLVBalancesStartDate = dateUtils.getEpochOneYearAgo()
    const lvBalancesStartDateParam =
      Number.isInteger(newLVBalancesStartDate) &&
      newLVBalancesStartDate !== defaultLVBalancesStartDate
        ? newLVBalancesStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        lvFilter: newLVFilter || undefined,
        lvBalancesStartDate: lvBalancesStartDateParam,
        lvBalancesEndDate: newLVBalancesEndDate || undefined,
        lvPageSize: lvPageSizeParam,
      },
    })
    debouncedFilterBankAccounts()
  }
)

const showAddBankAccountDialog = () => {
  bankAccountsStore.resetLVAddBankAccountDialog()
  lvAddBankAccountDialog.showModal()
}

function refetch() {
  const query = route.query

  const parsedLVPageSize = numUtils.queryParamToInt(
    query.lvPageSize,
    defaultPageSize
  )

  const parsedLVBalancesStartDate = numUtils.queryParamToNullableInt(
    query.lvBalancesStartDate
  )
  bankAccountsStore.lvBalancesStartDate = parsedLVBalancesStartDate

  const parsedLVBalancesEndDate = numUtils.queryParamToNullableInt(
    query.lvBalancesEndDate
  )
  bankAccountsStore.lvBalancesEndDate = parsedLVBalancesEndDate

  const defaultLVBalancesStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.lvHydrate(
    query.lvFilter?.toString() || "",
    parsedLVBalancesStartDate &&
      parsedLVBalancesStartDate !== defaultLVBalancesStartDate
      ? parsedLVBalancesStartDate
      : defaultLVBalancesStartDate,
    parsedLVBalancesEndDate,
    parsedLVPageSize
  )

  debouncedFilterBankAccounts()
}

onMounted(() => refetch())
onUnmounted(() => bankAccountsStore.lvDehydrate())

const createBankAccount = async () => {
  const res = await bankAccountsStore.createBankAccount()
  if (!res.errorMessage) {
    lvAddBankAccountDialog.close()
    bankAccountsStore.resetLVAddBankAccountDialog()
  }
}

const showDeleteBankAccountConfirmation = (accountId) => {
  bankAccountsStore.getAccountToDeleteById(accountId)
  lvDeleteBankAccountDialog.showModal()
}

const cancelDeleteBankAccount = () => {
  lvDeleteBankAccountDialog.close()
  bankAccountsStore.resetLVDeleteBankAccountDialog()
}

const deleteBankAccount = async () => {
  const res = await bankAccountsStore.deleteBankAccount()
  if (!res.errorMessage) {
    lvDeleteBankAccountDialog.close()
    bankAccountsStore.resetLVDeleteBankAccountDialog()
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: List of Accounts -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div class="flex gap-3">
            <h2 class="card-title">List of Accounts</h2>
            <button
              class="btn btn-neutral btn-xs"
              v-on:click="showAddBankAccountDialog()"
            >
              <font-awesome-icon :icon="['fas', 'plus']" />
              New Account
            </button>
          </div>
          <input
            type="text"
            v-model="bankAccountsStore.lvFilter"
            placeholder="Search accounts..."
            class="input input-bordered w-64"
          />
        </div>
        <div class="overflow-x-auto h-88">
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
                v-for="(account, index) in bankAccountsStore.lvBankAccounts"
                :key="index"
                class="hover:bg-base-300"
              >
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div
                        class="font-bold"
                        :class="{ 'opacity-50': account.status == 'inactive' }"
                      >
                        {{ account.accountName }}
                      </div>
                      <div class="text-sm opacity-50">
                        {{ account.accountHolderName }}
                      </div>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <div>
                      <div
                        class="font-bold"
                        :class="{ 'opacity-50': account.status == 'inactive' }"
                      >
                        {{ account.bankName }}
                      </div>
                      <div class="text-sm opacity-50">
                        # {{ account.accountNumber }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="text-right">
                  <div class="items-end">
                    <div>
                      <div
                        class="font-bold"
                        :class="{ 'opacity-50': account.status == 'inactive' }"
                      >
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
                    <span
                      class="badge badge-sm badge-neutral"
                      :class="{ 'opacity-50': account.status == 'inactive' }"
                      >{{ account.status }}</span
                    >
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-3">
                    <router-link
                      :to="{
                        name: 'assets.bankaccount.detail',
                        params: { id: account.id },
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
                      v-on:click="showDeleteBankAccountConfirmation(account.id)"
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

    <!-- Bottom Half: Line Chart of the Accounts' Balances -->
    <div class="card bg-base-100 shadow-md flex flex-1 min-h-0">
      <div class="card-body flex-1 min-h-0">
        <h2 class="card-title">Balance History Chart</h2>
        <line-chart class="flex-1 min-h-0" />
      </div>
    </div>
  </div>

  <!-- Dialog Box: Account Adder -->
  <dialog id="lvAddBankAccountDialog" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Add New Bank Account</h3>
      <form class="grid grid-cols-1 gap-4">
        <div>
          <label class="label">Account Name*</label>
          <input
            v-model="bankAccountsStore.lvAddBankAccount.accountName"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Bank Name*</label>
          <input
            v-model="bankAccountsStore.lvAddBankAccount.bankName"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Account Holder Name*</label>
          <input
            v-model="bankAccountsStore.lvAddBankAccount.accountHolderName"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Account Number*</label>
          <input
            v-model="bankAccountsStore.lvAddBankAccount.accountNumber"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Initial Balance*</label>
          <input
            v-model="bankAccountsStore.lvAddBankAccount.lastBalance"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Initial Balance Date*</label>
          <DatePicker
            v-model:date="bankAccountsStore.lvAddBankAccount.lastBalanceDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="createBankAccount"
            class="btn btn-primary"
          >
            Save
          </button>
          <button
            type="button"
            @click="bankAccountsStore.resetLVAddBankAccountDialog()"
            class="btn btn-secondary"
          >
            Reset
          </button>
        </div>
      </form>
    </div>
  </dialog>

  <!-- Dialog Box: Confirm Account Delete -->
  <dialog id="lvDeleteBankAccountDialog" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Confirm Account Delete</h3>
      <form class="grid grid-cols-1 gap-4">
        <div class="grid grid-cols-2 grid-rows-7 gap-4">
          <div>Account Name</div>
          <div>
            {{ bankAccountsStore.lvDeleteBankAccount.accountName }}
          </div>
          <div>Bank Name</div>
          <div>{{ bankAccountsStore.lvDeleteBankAccount.bankName }}</div>
          <div>Account Holder Name</div>
          <div>
            {{ bankAccountsStore.lvDeleteBankAccount.accountHolderName }}
          </div>
          <div>Account Number</div>
          <div>
            {{ bankAccountsStore.lvDeleteBankAccount.accountNumber }}
          </div>
          <div>Status</div>
          <div>{{ bankAccountsStore.lvDeleteBankAccount.status }}</div>
          <div>Last Balance</div>
          <div>
            {{
              numUtils.numericToMoney(
                bankAccountsStore.lvDeleteBankAccount.lastBalance
              )
            }}
          </div>
          <div>Last Balance Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                bankAccountsStore.lvDeleteBankAccount.lastBalanceDate
              )
            }}
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="deleteBankAccount()"
            class="btn btn-primary"
          >
            Confirm
          </button>
          <button
            type="button"
            @click="cancelDeleteBankAccount()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>