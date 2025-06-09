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
    const listViewFilter = ref('')
    const listViewBalancesStartDate = ref(0)
    const listViewBalancesEndDate = ref(0)
    const listViewPageSize = ref(10)
    const listViewBankAccounts = ref([])
    const listViewChartData = ref([])

    //// detail view
    const detailViewBankAccountId = ref('')
    const detailViewBalancesStartDate = ref(0)
    const detailViewBalancesEndDate = ref(0)
    const detailViewPageSize = ref(10)
    const detailViewAccount = ref({})
    const detailViewAccountCache = ref({})
    const detailViewChartData = ref([])

    //// account adder
    const aaBankAccount = ref({})
    const aaBankAccountCache = ref({})

    //// account deleter
    const adAccount = ref({})

    //// balance editor
    const balanceEditorMode = ref('Add')
    const beBalance = ref({})
    const beBalanceCache = ref({})

    //// balance deleter
    const bdBalance = ref({})

    ////// actions

    //// list view
    // hydration
    async function hydrate(initListViewFilter, initListViewBalancesStartDate, initListViewBalancesEndDate, initListViewPageSize) {
        listViewFilter.value = initListViewFilter
        listViewBalancesStartDate.value = initListViewBalancesStartDate
        listViewBalancesEndDate.value = initListViewBalancesEndDate
        listViewPageSize.value = initListViewPageSize
    }

    function dehydrate() {
        listViewFilter.value = ''
        listViewBalancesStartDate.value = 0
        listViewBalancesEndDate.value = 0
        listViewPageSize.value = 10
        listViewBankAccounts.value = []
        listViewChartData.value = []
        aaBankAccount.value = {}
        aaBankAccountCache.value = {}
        adAccount.value = {}
    }

    // CRUD

    async function createBankAccount() {
        const res = await svc.createBankAccount({
            accountName: detailViewAccount.value.accountName,
            bankName: detailViewAccount.value.bankName,
            accountHolderName: detailViewAccount.value.accountHolderName,
            accountNumber: detailViewAccount.value.accountNumber,
            lastBalance: detailViewAccount.value.lastBalance,
            lastBalanceDate: detailViewAccount.value.lastBalanceDate,
            status: detailViewAccount.value.status,
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

    async function getById(id) {
        detailViewAccount.value = await svc.getBankAccount(id, null, null, 0)
    }

    async function update() {
        const res = await svc.updateBankAccount(detailViewAccount.value)
        if (!res.errorMessage) {
            detailViewAccount.value = JSON.parse(JSON.stringify(res))
            detailViewAccountCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Account updated!', 'success')
        } else {
            toast.showToast('Failed to save account: ' + res.errorMessage, 'error')
        }
    }

    async function deleteAccount() {
        const res = await svc.deleteBankAccount(detailViewAccount.value.id)
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

    // cache and prep

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

    async function hydrateDetail(initDetailViewBankAccountId, initDetailViewBalanceStartDate, initDetailViewBalanceEndDate, initDetailViewPageSize) {
        detailViewBankAccountId.value = initDetailViewBankAccountId
        detailViewBalancesStartDate.value = initDetailViewBalanceStartDate
        detailViewBalancesEndDate.value = initDetailViewBalanceEndDate
        detailViewPageSize.value = initDetailViewPageSize
    }

    function dehydrateDetail() {
        detailViewBankAccountId.value = ''
        detailViewBalancesStartDate.value = 0
        detailViewBalancesEndDate.value = 0
        detailViewPageSize.value = 10
        detailViewAccount.value = {}
        detailViewAccountCache.value = {}
        detailViewChartData.value = []
        balanceEditorMode.value = 'Add'
        beBalance.value = {}
        beBalanceCache.value = {}
        bdBalance.value = {}
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
            detailViewBankAccountId.value,
            detailViewBalancesStartDate.value,
            detailViewBalancesEndDate.value,
            detailViewPageSize.value)

        detailViewAccount.value = JSON.parse(JSON.stringify(fetchedAccount))
        detailViewAccountCache.value = JSON.parse(JSON.stringify(fetchedAccount))

        extractDetailViewChartData()
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
        beBalance.value = JSON.parse(JSON.stringify(blankBankAccountBalance))
        beBalanceCache.value = JSON.parse(JSON.stringify(blankBankAccountBalance))
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
        listViewFilter,
        listViewBalancesStartDate,
        listViewBalancesEndDate,
        listViewPageSize,
        listViewBankAccounts,
        listViewChartData,

        //// detail view
        detailViewBankAccountId,
        detailViewBalancesStartDate,
        detailViewBalancesEndDate,
        detailViewPageSize,
        detailViewAccount,
        detailViewAccountCache,
        detailViewChartData,

        //// account adder
        aaBankAccount,
        aaBankAccountCache,

        //// account deleter
        adAccount,

        //// balance editor
        balanceEditorMode,
        beBalance,
        beBalanceCache,

        //// balance deleter
        bdBalance,

        ////// actions

        //// list view
        // hydration
        hydrate,
        dehydrate,
        // CRUD
        createBankAccount,
        filterBankAccounts,
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
