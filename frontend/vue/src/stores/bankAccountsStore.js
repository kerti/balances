import { useToast } from '@/composables/useToast'
import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { errorMessages } from 'vue/compiler-sfc'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()
    const toast = useToast()

    ////// templates
    const blankAccount = {}
    const blankAccountBalance = {}

    ////// reactive state
    //// list view
    const filter = ref('')
    const balancesStartDate = ref(0)
    const balancesEndDate = ref(0)
    const pageSize = ref(10)
    const accounts = ref([])
    const chartData = ref([])
    //// detail view
    const detailId = ref('')
    const detailBalanceStartDate = ref(0)
    const detailBalanceEndDate = ref(0)
    const detailPageSize = ref(10)
    const account = ref({})
    const accountCache = ref({})
    const detailChartData = ref([])
    //// balance editor
    const balanceEditorMode = ref('Add')
    const beBalance = ref({})
    const beBalanceCache = ref({})
    //// account adder

    ////// actions

    //// list view
    // hydration
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

    // CRUD

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
            toast.showToast('Failed to create account: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    async function search() {
        accounts.value = await svc.searchBankAccounts(
            filter.value,
            balancesStartDate.value,
            balancesEndDate.value,
            pageSize.value)
        extractChartData()
    }

    async function getById(id) {
        account.value = await svc.getBankAccount(id, null, null, 0)
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

    async function deleteAccount() {
        const res = await svc.deleteBankAccount(account.value.id)
        if (!res.errorMessage) {
            search()
            toast.showToast('Account deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete account: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    // cache and prep

    function revertAccountToCache() {
        if (accountCache.value) {
            account.value = JSON.parse(JSON.stringify(accountCache.value))
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

    // chart utils

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

    //// detail view

    // hydration

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

    // CRUD

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

    async function getBalanceById(id) {
        const fetchedBalance = await svc.getBankAccountBalance(id)
        beBalance.value = JSON.parse(JSON.stringify(fetchedBalance))
        beBalanceCache.value = JSON.parse(JSON.stringify(fetchedBalance))
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

    async function deleteAccountBalance() {
        const res = await svc.deleteBankAccountBalance(beBalance.value.id)
        if (!res.errorMessage) {
            get()
            toast.showToast('Balance deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete balance: ' + res.errorMessage)
            return {
                errorMessages: res.errorMessage
            }
        }
    }

    function revertBalanceToCache() {
        if (beBalanceCache.value) {
            beBalance.value = JSON.parse(JSON.stringify(beBalanceCache.value))
        }
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

    // chart utils

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

    return {
        ////// reactive state

        //// list view
        filter,
        balancesStartDate,
        balancesEndDate,
        pageSize,
        accounts,
        chartData,

        //// detail view
        detailId,
        detailBalanceStartDate,
        detailBalanceEndDate,
        detailPageSize,
        account,
        accountCache,
        detailChartData,

        //// balance editor
        balanceEditorMode,
        beBalance,
        beBalanceCache,

        //// account adder

        ////// actions

        //// list view
        // hydration
        hydrate,
        dehydrate,
        // CRUD
        createAccount,
        search,
        getById,
        update,
        deleteAccount,
        // cache and prep
        revertAccountToCache,
        prepBlankAccount,

        //// detail view
        // hydration
        hydrateDetail,
        dehydrateDetail,
        // CRUD
        createBalance,
        get,
        getBalanceById,
        updateBalance,
        deleteAccountBalance,
        // cache and prep
        revertBalanceToCache,
        prepBlankBalance,
    }
})
