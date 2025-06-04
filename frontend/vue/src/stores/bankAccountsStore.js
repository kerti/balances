import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()

    //// reactive state
    // list view
    const filter = ref('')
    const balancesStartDate = ref(0)
    const balancesEndDate = ref(0)
    const pageSize = ref(10)
    const accounts = ref([])
    const chartData = ref([])
    // detail view
    const detailId = ref('')
    const detailBalanceStartDate = ref(0)
    const detailBalanceEndDate = ref(0)
    const detailPageSize = ref(10)
    const account = ref({})
    const accountCache = ref({})
    const detailChartData = ref([])

    //// actions
    async function hydrate(initFilter, initBalancesStartDate, initBalancesEndDate, initPageSize) {
        filter.value = initFilter
        balancesStartDate.value = initBalancesStartDate
        balancesEndDate.value = initBalancesEndDate
        pageSize.value = initPageSize
    }

    function dehydrate() {
        filter.value = ''
        balancesStartDate.value = 0
        balancesEndDate.value = 0
        pageSize.value = 10
        accounts.value = []
        chartData.value = []
    }

    async function hydrateDetail(initId, initBalanceStartDate, initBalanceEndDate, initPageSize) {
        detailId.value = initId
        detailBalanceStartDate.value = initBalanceStartDate
        detailBalanceEndDate.value = initBalanceEndDate
        detailPageSize.value = initPageSize
    }

    function dehydrateDetail() {
        detailId.value = ''
        detailBalanceStartDate.value = 0
        detailBalanceEndDate.value = 0
        detailPageSize.value = 10
        account.value = {}
        accountCache.value = {}
        detailChartData.value = []
    }

    async function search() {
        accounts.value = await svc.searchBankAccounts(
            filter.value,
            balancesStartDate.value,
            balancesEndDate.value,
            pageSize.value)
        extractChartData()
    }

    async function get() {
        const fetchedAccount = await svc.getBankAccount(
            detailId.value,
            balancesStartDate.value,
            balancesEndDate.value,
            detailPageSize.value)

        account.value = JSON.parse(JSON.stringify(fetchedAccount))
        accountCache.value = JSON.parse(JSON.stringify(fetchedAccount))

        extractDetailChartData()
    }

    function extractChartData() {
        chartData.value = accounts.value.map(acc => {
            return {
                label: acc.accountName,
                data: acc.balances.map(balance => {
                    return {
                        x: balance.date,
                        y: balance.balance,
                    }
                })
            }
        })
    }

    function extractDetailChartData() {
        detailChartData.value = [{
            label: account.value.accountName,
            data: account.value.balances.map(balance => {
                return {
                    x: balance.date,
                    y: balance.balance,
                }
            })
        }]
    }

    function revertAccountToCache() {
        if (accountCache.value) {
            account.value = JSON.parse(JSON.stringify(accountCache.value))
        }
    }

    async function update() {
        svc.updateBankAccount(account.value)
    }

    return {
        //// reactive state
        // list view
        filter,
        balancesStartDate,
        balancesEndDate,
        pageSize,
        accounts,
        chartData,
        // detail view
        detailId,
        detailBalanceStartDate,
        detailBalanceEndDate,
        detailPageSize,
        account,
        accountCache,
        detailChartData,
        //// actions
        hydrate,
        dehydrate,
        hydrateDetail,
        dehydrateDetail,
        search,
        get,
        revertAccountToCache,
        update,
    }
})
