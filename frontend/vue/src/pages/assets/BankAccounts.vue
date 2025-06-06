<script setup>
import debounce from "lodash.debounce"
import { useDateUtils } from "@/composables/useDateUtils"
import { useEnvUtils } from "@/composables/useEnvUtils"
import { useNumUtils } from "@/composables/useNumUtils"
import { useBankAccountsStore } from "@/stores/bankAccountsStore"
import { watch, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import LineChart from "@/components/assets/BankLineChart.vue"

const dateUtils = useDateUtils()
const numUtils = useNumUtils()
const ev = useEnvUtils()
const route = useRoute()
const router = useRouter()
const bankAccountsStore = useBankAccountsStore()
const defaultPageSize = ev.getDefaultPageSize()

const debouncedSearch = debounce(() => {
  bankAccountsStore.search()
}, 300)

watch(
  [
    () => bankAccountsStore.filter,
    () => bankAccountsStore.balancesStartDate,
    () => bankAccountsStore.balancesEndDate,
    () => bankAccountsStore.pageSize,
  ],
  ([newFilter, newBalancesStartDate, newBalancesEndDate, newPageSize]) => {
    const pageSizeParam =
      Number.isInteger(newPageSize) && newPageSize !== defaultPageSize
        ? newPageSize
        : undefined

    const defaultBalancesStartDate = dateUtils.getEpochOneYearAgo()
    const balancesStartDateParam =
      Number.isInteger(newBalancesStartDate) &&
      newBalancesStartDate !== defaultBalancesStartDate
        ? newBalancesStartDate
        : undefined

    router.replace({
      query: {
        ...route.query,
        filter: newFilter || undefined,
        balancesStartDate: balancesStartDateParam,
        balancesEndDate: newBalancesEndDate || undefined,
        pageSize: pageSizeParam,
      },
    })
    debouncedSearch()
  }
)

function refetch() {
  const query = route.query

  const parsedPageSize = numUtils.queryParamToInt(
    query.pageSize,
    defaultPageSize
  )

  const parsedBalancesStartDate = numUtils.queryParamToNullableInt(
    query.balancesStartDate
  )
  bankAccountsStore.balancesStartDate = parsedBalancesStartDate

  const parsedBalancesEndDate = numUtils.queryParamToNullableInt(
    query.balancesEndDate
  )
  bankAccountsStore.balancesEndDate = parsedBalancesEndDate

  const defaultBalancesStartDate = dateUtils.getEpochOneYearAgo()
  bankAccountsStore.hydrate(
    query.filter?.toString() || "",
    parsedBalancesStartDate &&
      parsedBalancesStartDate !== defaultBalancesStartDate
      ? parsedBalancesStartDate
      : defaultBalancesStartDate,
    parsedBalancesEndDate,
    parsedPageSize
  )
}

onMounted(() => refetch())
onUnmounted(() => bankAccountsStore.dehydrate())
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
              class="input input-bordered w-64 mr-1"
            />
            <button
              class="btn btn-neutral btn-circle tooltip ml-1"
              data-tip="Add New Bank Account"
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
                    <button class="btn btn-neutral tooltip" data-tip="Activate">
                      <font-awesome-icon :icon="['fas', 'eye']" />
                    </button>
                    <button
                      class="btn btn-neutral tooltip"
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
        <line-chart />
      </div>
    </div>
  </div>
</template>