import { useToast } from '@/composables/useToast'
import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()
    const toast = useToast()

    ////// templates
    const blankBankAccount = {
        id: '',
        accountName: '',
        bankName: '',
        accountHolderName: '',
        accountNumber: '',
        lastBalance: 0,
        lastBalanceDate: (new Date()).getTime(),
        status: 'active',
    }
    const blankBankAccountBalance = {
        id: '',
        bankAccountId: '',
        date: (new Date()).getTime(),
        balance: 0,
    }

    ////// reactive state

    //// list view
    // main page
    const lvFilter = ref('')
    const lvBalancesStartDate = ref(0)
    const lvBalancesEndDate = ref(0)
    const lvPageSize = ref(10)
    const lvBankAccounts = ref([])
    const lvChartData = ref([])
    // add bank account dialog box
    const lvAddBankAccount = ref({})
    // delete bank account dialog box
    const lvDeleteBankAccount = ref({})

    //// detail view
    // main page
    const dvBankAccountId = ref('')
    const dvBalancesStartDate = ref(0)
    const dvBalancesEndDate = ref(0)
    const dvPageSize = ref(10)
    const dvAccount = ref({})
    const dvAccountCache = ref({})
    const dvChartData = ref([])
    // balance editor dialog box
    const dvBalanceEditorMode = ref('Add')
    const dvEditBankAccountBalance = ref({})
    const dvEditBankAccountBalanceCache = ref({})

    ////// actions

    //// list view
    // hydration
    async function lvHydrate(initLVFilter, initLVBalancesStartDate, initLVBalancesEndDate, initLVPageSize) {
        lvFilter.value = initLVFilter
        lvBalancesStartDate.value = initLVBalancesStartDate
        lvBalancesEndDate.value = initLVBalancesEndDate
        lvPageSize.value = initLVPageSize
    }

    function lvDehydrate() {
        lvFilter.value = ''
        lvBalancesStartDate.value = 0
        lvBalancesEndDate.value = 0
        lvPageSize.value = 10
        lvBankAccounts.value = []
        lvChartData.value = []
        lvAddBankAccount.value = {}
    }

    // CRUD

    async function createBankAccount() {
        const res = await svc.createBankAccount({
            accountName: lvAddBankAccount.value.accountName,
            bankName: lvAddBankAccount.value.bankName,
            accountHolderName: lvAddBankAccount.value.accountHolderName,
            accountNumber: lvAddBankAccount.value.accountNumber,
            lastBalance: lvAddBankAccount.value.lastBalance,
            lastBalanceDate: lvAddBankAccount.value.lastBalanceDate,
            status: lvAddBankAccount.value.status,
        })
        if (!res.error) {
            filterBankAccounts()
            toast.showToast('Account created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create account: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function filterBankAccounts() {
        lvBankAccounts.value = await svc.searchBankAccounts(
            lvFilter.value,
            lvBalancesStartDate.value,
            lvBalancesEndDate.value,
            lvPageSize.value)
        extractLVChartData()
    }

    async function getAccountToDeleteById(id) {
        lvDeleteBankAccount.value = await svc.getBankAccount(id, null, null, 0)
    }

    async function deleteBankAccount() {
        const res = await svc.deleteBankAccount(lvDeleteBankAccount.value.id)
        if (!res.error) {
            filterBankAccounts()
            toast.showToast('Account deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete account: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    // cache prep and reset

    function resetLVAddBankAccountDialog() {
        lvAddBankAccount.value = JSON.parse(JSON.stringify(blankBankAccount))
    }

    function resetLVDeleteBankAccountDialog() {
        lvDeleteBankAccount.value = JSON.parse(JSON.stringify(blankBankAccount))
    }

    // chart utils

    function extractLVChartData() {
        lvChartData.value = lvBankAccounts.value.map(acc => {
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

    async function dvHydrate(initDVBankAccountId, initDVBalanceStartDate, initDVBalanceEndDate, initDVPageSize) {
        dvBankAccountId.value = initDVBankAccountId
        dvBalancesStartDate.value = initDVBalanceStartDate
        dvBalancesEndDate.value = initDVBalanceEndDate
        dvPageSize.value = initDVPageSize
    }

    function dvDehydrate() {
        dvBankAccountId.value = ''
        dvBalancesStartDate.value = 0
        dvBalancesEndDate.value = 0
        dvPageSize.value = 10
        dvAccount.value = {}
        dvAccountCache.value = {}
        dvChartData.value = []
        dvBalanceEditorMode.value = 'Add'
        dvEditBankAccountBalance.value = {}
        dvEditBankAccountBalanceCache.value = {}
    }

    // CRUD

    async function createBankAccountBalance() {
        const res = await svc.createBankAccountBalance({
            bankAccountId: dvEditBankAccountBalance.value.bankAccountId,
            date: dvEditBankAccountBalance.value.date,
            balance: dvEditBankAccountBalance.value.balance
        })
        if (!res.error) {
            getBankAccountForDV()
            toast.showToast('Balance created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create balance: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function getBankAccountForDV() {
        const fetchedAccount = await svc.getBankAccount(
            dvBankAccountId.value,
            dvBalancesStartDate.value,
            dvBalancesEndDate.value,
            dvPageSize.value)

        dvAccount.value = JSON.parse(JSON.stringify(fetchedAccount))
        dvAccountCache.value = JSON.parse(JSON.stringify(fetchedAccount))

        extractDVChartData()
    }

    async function getBankAccountBalanceById(id) {
        const fetchedBalance = await svc.getBankAccountBalance(id)
        dvEditBankAccountBalance.value = JSON.parse(JSON.stringify(fetchedBalance))
        dvEditBankAccountBalanceCache.value = JSON.parse(JSON.stringify(fetchedBalance))
    }

    async function updateBankAccount() {
        const res = await svc.updateBankAccount(dvAccount.value)
        if (!res.error) {
            // preserve balances records not fetched during bank account update
            res.balances = JSON.parse(JSON.stringify(dvAccountCache.value.balances))
            // then sync the store to the latest data from update
            dvAccount.value = JSON.parse(JSON.stringify(res))
            dvAccountCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Account updated!', 'success')
        } else {
            toast.showToast('Failed to save account: ' + res.error.message, 'error')
        }
    }

    async function updateBankAccountBalance() {
        const res = await svc.updateBankAccountBalance(dvEditBankAccountBalance.value)
        if (!res.error) {
            getBankAccountForDV()
            getBankAccountBalanceById(res.id)
            toast.showToast('Balance updated!', 'success')
            return res
        } else {
            toast.showToast('Failed to save balance: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function deleteBankAccountBalance() {
        const res = await svc.deleteBankAccountBalance(dvEditBankAccountBalance.value.id)
        if (!res.error) {
            getBankAccountForDV()
            toast.showToast('Balance deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete balance: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    // cache prep and reset

    function revertDVBankAccountToCache() {
        if (dvAccountCache.value) {
            dvAccount.value = JSON.parse(JSON.stringify(dvAccountCache.value))
        }
    }

    function revertDVBankAccountBalanceToCache() {
        if (dvEditBankAccountBalanceCache.value) {
            dvEditBankAccountBalance.value = JSON.parse(JSON.stringify(dvEditBankAccountBalanceCache.value))
        }
    }

    function prepDVBlankBankAccountBalance() {
        const template = JSON.parse(JSON.stringify(blankBankAccountBalance))
        template.bankAccountId = dvBankAccountId.value
        dvEditBankAccountBalance.value = JSON.parse(JSON.stringify(template))
        dvEditBankAccountBalanceCache.value = JSON.parse(JSON.stringify(template))
    }

    // chart utils

    function extractDVChartData() {
        dvChartData.value = [{
            label: dvAccount.value.accountName,
            data: dvAccount.value.balances.map(balance => {
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
        // main page
        lvFilter,
        lvBalancesStartDate,
        lvBalancesEndDate,
        lvPageSize,
        lvBankAccounts,
        lvChartData,
        // add bank account dialog box
        lvAddBankAccount,
        // delete bank account dialog box
        lvDeleteBankAccount,

        //// detail view
        // main page
        dvBankAccountId,
        dvBalancesStartDate,
        dvBalancesEndDate,
        dvPageSize,
        dvAccount,
        dvAccountCache,
        dvChartData,
        // balance editor dialog box
        dvBalanceEditorMode,
        dvEditBankAccountBalance,
        dvEditBankAccountBalanceCache,

        ////// actions

        //// list view
        // hydration
        lvHydrate,
        lvDehydrate,
        // CRUD
        createBankAccount,
        filterBankAccounts,
        getAccountToDeleteById,
        deleteBankAccount,
        // cache and prep
        resetLVAddBankAccountDialog,
        resetLVDeleteBankAccountDialog,

        //// detail view
        // hydration
        dvHydrate,
        dvDehydrate,
        // CRUD
        createBankAccountBalance,
        getBankAccountForDV,
        getBankAccountBalanceById,
        updateBankAccount,
        updateBankAccountBalance,
        deleteBankAccountBalance,
        // cache and prep
        revertDVBankAccountToCache,
        revertDVBankAccountBalanceToCache,
        prepDVBlankBankAccountBalance,
    }
})
