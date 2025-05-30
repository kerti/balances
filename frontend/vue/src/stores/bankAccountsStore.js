import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()

    // reactive state
    const filter = ref('')
    const pageSize = ref(10)
    const accounts = ref([])

    // actions
    async function hydrate(initFilter, initPageSize) {
        filter.value = initFilter
        pageSize.value = initPageSize
        await search(filter.value, pageSize.value)
    }

    async function search(filter, pageSize) {
        accounts.value = await svc.searchBankAccounts(filter, pageSize)
    }

    return {
        // reactive state
        filter,
        pageSize,
        accounts,
        // actions
        hydrate,
        search,
    }
})
