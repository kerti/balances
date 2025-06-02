import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()

    // reactive state
    const filter = ref('')
    const balancesStartDate = ref(0)
    const balancesEndDate = ref(0)
    const pageSize = ref(10)
    const accounts = ref([])

    // actions
    async function hydrate(initFilter, initBalancesStartDate, initBalancesEndDate, initPageSize) {
        filter.value = initFilter
        balancesStartDate.value = initBalancesStartDate
        balancesEndDate.value = initBalancesEndDate
        pageSize.value = initPageSize
    }

    async function search(filter, balancesStartDate, balancesEndDate, pageSize) {
        accounts.value = await svc.searchBankAccounts(filter, balancesStartDate, balancesEndDate, pageSize)
    }

    return {
        // reactive state
        filter,
        balancesStartDate,
        balancesEndDate,
        pageSize,
        accounts,
        // actions
        hydrate,
        search,
    }
})
