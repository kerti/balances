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
    () => bankAccountsStore.listViewFilter,
    () => bankAccountsStore.listViewBalancesStartDate,
    () => bankAccountsStore.listViewBalancesEndDate,
    () => bankAccountsStore.listViewPageSize,
  ],
  ([
    newListViewFilter,
    newListViewBalancesStartDate,
    newListViewBalancesEndDate,
    newPageSize,
  ]) => {
    const listViewPageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultListViewBalancesStartDate = dateUtils.getEpochOneYearAgo()
    const listViewBalancesStartDateParam =
      Number.isInteger(newListViewBalancesStartDate) &&
      newListViewBalancesStartDate !== defaultListViewBalancesStartDate
        ? newListViewBalancesStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        filter: newListViewFilter || undefined,
        listViewBalancesStartDate: listViewBalancesStartDateParam,
        listViewBalancesEndDate: newListViewBalancesEndDate || undefined,
        listViewPageSize: listViewPageSizeParam,
      },
    })
    debouncedFilterBankAccounts()
  }
)

const showAdder = () => {
  bankAccountsStore.prepBlankAccount()
  accountAdder.showModal()
}

function refetch() {
  const query = route.query

  const parsedPageSize = numUtils.queryParamToInt(
    query.listViewPageSize,
    defaultPageSize
  )

  const parsedListViewBalancesStartDate = numUtils.queryParamToNullableInt(
    query.listViewBalancesStartDate
  )
  bankAccountsStore.listViewBalancesStartDate = parsedListViewBalancesStartDate

  const parsedBalancesEndDate = numUtils.queryParamToNullableInt(
    query.listViewBalancesEndDate
  )
  bankAccountsStore.listViewBalancesEndDate = parsedBalancesEndDate

  const defaultListViewBalancesStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.hydrate(
    query.filter?.toString() || "",
    parsedListViewBalancesStartDate &&
      parsedListViewBalancesStartDate !== defaultListViewBalancesStartDate
      ? parsedListViewBalancesStartDate
      : defaultListViewBalancesStartDate,
    parsedBalancesEndDate,
    parsedPageSize
  )
}

onMounted(() => refetch())
onUnmounted(() => bankAccountsStore.dehydrate())

const saveAccount = async () => {
  const res = await bankAccountsStore.createBankAccount()
  if (!res.errorMessage) {
    accountAdder.close()
  }
}

const showAccountDeleteConfirmation = (accountId) => {
  bankAccountsStore.getById(accountId)
  adConfirm.showModal()
}

const cancelAccountDelete = () => {
  adConfirm.close()
}

const deleteAccount = async () => {
  const res = await bankAccountsStore.deleteAccount()
  if (!res.errorMessage) {
    adConfirm.close()
  }
}
</script>

<template>
  <div class="flex flex-col h-full space-y-6">
    <!-- Top Half: List of Accounts -->
    <div class="card bg-base-100 shadow-md">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <h2 class="card-title">List of Accounts</h2>
          <div class="form-control flex gap-3">
            <input
              type="text"
              v-model="bankAccountsStore.listViewFilter"
              placeholder="Search accounts..."
              class="input input-bordered w-64"
            />
            <button
              class="btn btn-neutral btn-circle tooltip"
              data-tip="Add New Bank Account"
              v-on:click="showAdder()"
            >
              <font-awesome-icon :icon="['fas', 'plus']" />
            </button>
          </div>
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
                v-for="(
                  account, index
                ) in bankAccountsStore.listViewBankAccounts"
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
                    <button class="btn btn-neutral tooltip" data-tip="Edit">
                      <router-link
                        :to="{
                          name: 'assets.bankaccount.detail',
                          params: { id: account.id },
                        }"
                      >
                        <font-awesome-icon :icon="['fas', 'edit']" />
                      </router-link>
                    </button>
                    <button
                      class="btn btn-neutral tooltip"
                      data-tip="Delete"
                      v-on:click="showAccountDeleteConfirmation(account.id)"
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
  <dialog id="accountAdder" class="modal">
    <div class="modal-box overflow-visible relative z-50">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold pb-5">Add New Bank Account</h3>
      <form class="grid grid-cols-1 gap-4">
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
          <label class="label">Initial Balance</label>
          <input
            v-model="bankAccountsStore.account.lastBalance"
            type="text"
            class="input input-bordered w-full"
          />
        </div>
        <div>
          <label class="label">Initial Balance Date</label>
          <DatePicker
            v-model:date="bankAccountsStore.account.lastBalanceDate"
            placeholder="pick a date"
            required
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button type="button" @click="saveAccount" class="btn btn-primary">
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
  </dialog>

  <!-- Dialog Box: Confirm Account Delete -->
  <dialog id="adConfirm" class="modal">
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
          <div>{{ bankAccountsStore.account.accountName }}</div>
          <div>Bank Name</div>
          <div>{{ bankAccountsStore.account.bankName }}</div>
          <div>Account Holder Name</div>
          <div>{{ bankAccountsStore.account.accountHolderName }}</div>
          <div>Account Number</div>
          <div>{{ bankAccountsStore.account.accountNumber }}</div>
          <div>Status</div>
          <div>{{ bankAccountsStore.account.status }}</div>
          <div>Last Balance</div>
          <div>
            {{ numUtils.numericToMoney(bankAccountsStore.account.lastBalance) }}
          </div>
          <div>Last Balance Date</div>
          <div>
            {{
              dateUtils.epochToLocalDate(
                bankAccountsStore.account.lastBalanceDate
              )
            }}
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <button
            type="button"
            @click="deleteAccount()"
            class="btn btn-primary"
          >
            Confirm
          </button>
          <button
            type="button"
            @click="cancelAccountDelete()"
            class="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  </dialog>
</template>