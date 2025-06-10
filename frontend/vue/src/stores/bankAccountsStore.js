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
    const listViewFilter = ref('')
    const listViewBalancesStartDate = ref(0)
    const listViewBalancesEndDate = ref(0)
    const listViewPageSize = ref(10)
    const listViewBankAccounts = ref([])
    const listViewChartData = ref([])
    // add bank account dialog box
    const listViewAddBankAccount = ref({})
    // delete bank account dialog box
    const listViewDeleteBankAccount = ref({})

    //// detail view
    // main page
    const detailViewBankAccountId = ref('')
    const detailViewBalancesStartDate = ref(0)
    const detailViewBalancesEndDate = ref(0)
    const detailViewPageSize = ref(10)
    const detailViewAccount = ref({})
    const detailViewAccountCache = ref({})
    const detailViewChartData = ref([])
    // balance editor dialog box
    const detailViewBalanceEditorMode = ref('Add')
    const detailViewEditBankAccountBalance = ref({})
    const detailViewEditBankAccountBalanceCache = ref({})
    // balance delete confirmation dialog box
    const bdBalance = ref({})

    ////// actions

    //// list view
    // hydration
    async function lvHydrate(initListViewFilter, initListViewBalancesStartDate, initListViewBalancesEndDate, initListViewPageSize) {
        listViewFilter.value = initListViewFilter
        listViewBalancesStartDate.value = initListViewBalancesStartDate
        listViewBalancesEndDate.value = initListViewBalancesEndDate
        listViewPageSize.value = initListViewPageSize
    }

    function lvDehydrate() {
        listViewFilter.value = ''
        listViewBalancesStartDate.value = 0
        listViewBalancesEndDate.value = 0
        listViewPageSize.value = 10
        listViewBankAccounts.value = []
        listViewChartData.value = []
        listViewAddBankAccount.value = {}
    }

    // CRUD

    async function createBankAccount() {
        const res = await svc.createBankAccount({
            accountName: listViewAddBankAccount.value.accountName,
            bankName: listViewAddBankAccount.value.bankName,
            accountHolderName: listViewAddBankAccount.value.accountHolderName,
            accountNumber: listViewAddBankAccount.value.accountNumber,
            lastBalance: listViewAddBankAccount.value.lastBalance,
            lastBalanceDate: listViewAddBankAccount.value.lastBalanceDate,
            status: listViewAddBankAccount.value.status,
        })
        if (!res.errorMessage) {
            filterBankAccounts()
            toast.showToast('Account created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create account: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    async function filterBankAccounts() {
        listViewBankAccounts.value = await svc.searchBankAccounts(
            listViewFilter.value,
            listViewBalancesStartDate.value,
            listViewBalancesEndDate.value,
            listViewPageSize.value)
        extractListViewChartData()
    }

    async function getAccountToDeleteById(id) {
        listViewDeleteBankAccount.value = await svc.getBankAccount(id, null, null, 0)
    }

    async function deleteBankAccount() {
        const res = await svc.deleteBankAccount(listViewDeleteBankAccount.value.id)
        if (!res.errorMessage) {
            filterBankAccounts()
            toast.showToast('Account deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete account: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    // cache prep, and reset

    function resetListViewAddBankAccountDialog() {
        listViewAddBankAccount.value = JSON.parse(JSON.stringify(blankBankAccount))
    }

    function resetListViewDeleteBankAccountDialog() {
        listViewDeleteBankAccount.value = JSON.parse(JSON.stringify(blankBankAccount))
    }

    function revertAccountToCache() {
        if (detailViewAccountCache.value) {
            detailViewAccount.value = JSON.parse(JSON.stringify(detailViewAccountCache.value))
        }
    }

    function prepBlankAccount() {
        detailViewAccount.value = JSON.parse(JSON.stringify(blankBankAccount))
        detailViewAccountCache.value = JSON.parse(JSON.stringify(blankBankAccount))
    }

    // chart utils

    function extractListViewChartData() {
        listViewChartData.value = listViewBankAccounts.value.map(acc => {
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

    async function dvHydrate(initDetailViewBankAccountId, initDetailViewBalanceStartDate, initDetailViewBalanceEndDate, initDetailViewPageSize) {
        detailViewBankAccountId.value = initDetailViewBankAccountId
        detailViewBalancesStartDate.value = initDetailViewBalanceStartDate
        detailViewBalancesEndDate.value = initDetailViewBalanceEndDate
        detailViewPageSize.value = initDetailViewPageSize
    }

    function dvDehydrate() {
        detailViewBankAccountId.value = ''
        detailViewBalancesStartDate.value = 0
        detailViewBalancesEndDate.value = 0
        detailViewPageSize.value = 10
        detailViewAccount.value = {}
        detailViewAccountCache.value = {}
        detailViewChartData.value = []
        detailViewBalanceEditorMode.value = 'Add'
        detailViewEditBankAccountBalance.value = {}
        detailViewEditBankAccountBalanceCache.value = {}
        bdBalance.value = {}
    }

    // CRUD

    async function createBankAccountBalance() {
        const res = await svc.createBankAccountBalance({
            bankAccountId: detailViewEditBankAccountBalance.value.bankAccountId,
            date: detailViewEditBankAccountBalance.value.date,
            balance: detailViewEditBankAccountBalance.value.balance
        })
        if (!res.errorMessage) {
            getBankAccountForDetailView()
            toast.showToast('Balance created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create balance: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    async function getBankAccountForDetailView() {
        const fetchedAccount = await svc.getBankAccount(
            detailViewBankAccountId.value,
            detailViewBalancesStartDate.value,
            detailViewBalancesEndDate.value,
            detailViewPageSize.value)

        detailViewAccount.value = JSON.parse(JSON.stringify(fetchedAccount))
        detailViewAccountCache.value = JSON.parse(JSON.stringify(fetchedAccount))

        extractDetailViewChartData()
    }

    async function getBankAccountBalanceById(id) {
        const fetchedBalance = await svc.getBankAccountBalance(id)
        detailViewEditBankAccountBalance.value = JSON.parse(JSON.stringify(fetchedBalance))
        detailViewEditBankAccountBalanceCache.value = JSON.parse(JSON.stringify(fetchedBalance))
    }

    async function updateBankAccount() {
        const res = await svc.updateBankAccount(detailViewAccount.value)
        if (!res.errorMessage) {
            detailViewAccount.value = JSON.parse(JSON.stringify(res))
            detailViewAccountCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Account updated!', 'success')
        } else {
            toast.showToast('Failed to save account: ' + res.errorMessage, 'error')
        }
    }

    async function updateBankAccountBalance() {
        const res = await svc.updateBankAccountBalance(detailViewEditBankAccountBalance.value)
        if (!res.errorMessage) {
            getBankAccountForDetailView()
            getBankAccountBalanceById(res.id)
            toast.showToast('Balance updated!', 'success')
            return res
        } else {
            toast.showToast('Failed to save balance: ' + res.errorMessage)
            return {
                errorMessage: res.errorMessage
            }
        }
    }

    async function deleteBankAccountBalance() {
        const res = await svc.deleteBankAccountBalance(detailViewEditBankAccountBalance.value.id)
        if (!res.errorMessage) {
            getBankAccountForDetailView()
            toast.showToast('Balance deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete balance: ' + res.errorMessage)
            return {
                errorMessages: res.errorMessage
            }
        }
    }

    // cache prep and reset

    function revertBalanceToCache() {
        if (detailViewEditBankAccountBalanceCache.value) {
            detailViewEditBankAccountBalance.value = JSON.parse(JSON.stringify(detailViewEditBankAccountBalanceCache.value))
        }
    }

    function prepBlankBalance() {
        const template = JSON.parse(JSON.stringify(blankBankAccountBalance))
        template.bankAccountId = detailViewBankAccountId.value
        detailViewEditBankAccountBalance.value = JSON.parse(JSON.stringify(template))
        detailViewEditBankAccountBalanceCache.value = JSON.parse(JSON.stringify(template))
    }

    // chart utils

    function extractDetailViewChartData() {
        detailViewChartData.value = [{
            label: detailViewAccount.value.accountName,
            data: detailViewAccount.value.balances.map(balance => {
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
        listViewFilter,
        listViewBalancesStartDate,
        listViewBalancesEndDate,
        listViewPageSize,
        listViewBankAccounts,
        listViewChartData,
        // add bank account dialog box
        listViewAddBankAccount,
        // delete bank account dialog box
        listViewDeleteBankAccount,

        //// detail view
        // main page
        detailViewBankAccountId,
        detailViewBalancesStartDate,
        detailViewBalancesEndDate,
        detailViewPageSize,
        detailViewAccount,
        detailViewAccountCache,
        detailViewChartData,
        // balance editor dialog box
        detailViewBalanceEditorMode,
        detailViewEditBankAccountBalance,
        detailViewEditBankAccountBalanceCache,
        //// balance delete confirmation dialog box
        bdBalance,

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
        resetListViewAddBankAccountDialog,
        resetListViewDeleteBankAccountDialog,
        revertAccountToCache,
        prepBlankAccount,

        //// detail view
        // hydration
        dvHydrate,
        dvDehydrate,
        // CRUD
        createBankAccountBalance,
        getBankAccountForDetailView,
        getBankAccountBalanceById,
        updateBankAccount,
        updateBankAccountBalance,
        deleteBankAccountBalance,
        // cache and prep
        revertBalanceToCache,
        prepBlankBalance,
    }
})
