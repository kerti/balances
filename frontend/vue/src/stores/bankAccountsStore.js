import { useBankAccountsService } from '@/services/bankAccountsService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useBankAccountsStore = defineStore('bankAccounts', () => {
    const svc = useBankAccountsService()

    // reactive state
    // TODO: fill in the search parameters based on URL parameter
    const filter = ref('')
    const pageSize = ref(10)
    const accounts = ref([])

    // actions
    async function hydrate() {
        console.log('hydrating bank accounts store...')
        accounts.value = await search(filter.value, pageSize.value)
    }

    async function search(filter, pageSize) {
        const res = await svc.searchBankAccounts(filter, pageSize)
        return res
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
