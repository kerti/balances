import { useToast } from '@/composables/useToast'
import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()
    const toast = useToast()

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
    // balance editor
    const balanceEditorMode = ref('Add')
    const beBalance = ref({})
    const beBalanceCache = ref({})

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
            detailBalanceStartDate.value,
            detailBalanceEndDate.value,
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

    function revertBalanceToCache() {
        if (beBalanceCache.value) {
            beBalance.value = JSON.parse(JSON.stringify(beBalanceCache.value))
        }
    }

    async function update() {
        const res = await svc.updateBankAccount(account.value)
        if (!res.errorMessage) {
            account.value = JSON.parse(JSON.stringify(res))
            accountCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Account updated!', 'success')
        } else {
            toast.showToast('Failed to save account: ' + res.errorMessage, 'error')
        }
    }

    async function updateBalance() {
        const res = await svc.updateBankAccountBalance(beBalance.value)
        if (!res.errorMessage) {
            get()
            getBalanceById(res.id)
            toast.showToast('Balance updated!', 'success')
            return res
        } else {
            toast.showToast('Failed to save balance: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    async function getBalanceById(id) {
        const fetchedBalance = await svc.getBankAccountBalance(id)
        beBalance.value = JSON.parse(JSON.stringify(fetchedBalance))
        beBalanceCache.value = JSON.parse(JSON.stringify(fetchedBalance))
    }

    function prepBlankBalance() {
        const blankBalance = {
            bankAccountId: account.value.id,
            date: (new Date()).getTime(),
            balance: 0,
        }
        beBalance.value = JSON.parse(JSON.stringify(blankBalance))
        beBalanceCache.value = JSON.parse(JSON.stringify(blankBalance))
    }

    async function createBalance() {
        const res = await svc.createBankAccountBalance({
            bankAccountId: beBalance.value.bankAccountId,
            date: beBalance.value.date,
            balance: beBalance.value.balance
        })
        if (!res.errorMessage) {
            get()
            toast.showToast('Balance created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create balance: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    function prepBlankAccount() {
        const blankAccount = {
            accountName: '',
            bankName: '',
            accountHolderName: '',
            accountNumber: '',
            lastBalance: 0,
            lastBalanceDate: (new Date()).getTime(),
            status: 'active',
        }
        account.value = JSON.parse(JSON.stringify(blankAccount))
        accountCache.value = JSON.parse(JSON.stringify(blankAccount))
    }

    async function createAccount() {
        const res = await svc.createBankAccount({
            accountName: account.value.accountName,
            bankName: account.value.bankName,
            accountHolderName: account.value.accountHolderName,
            accountNumber: account.value.accountNumber,
            lastBalance: account.value.lastBalance,
            lastBalanceDate: account.value.lastBalanceDate,
            status: account.value.status,
        })
        if (!res.errorMessage) {
            search()
            toast.showToast('Account created!', 'success')
            return res
        } else {
            toast.showToasst('Failed to create account: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
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
        // balance editor
        balanceEditorMode,
        beBalance,
        beBalanceCache,
        //// actions
        hydrate,
        dehydrate,
        hydrateDetail,
        dehydrateDetail,
        search,
        get,
        revertAccountToCache,
        revertBalanceToCache,
        update,
        updateBalance,
        getBalanceById,
        prepBlankBalance,
        createBalance,
        prepBlankAccount,
        createAccount,
    }
})
