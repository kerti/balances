import { searchBankAccountsFromAPI } from '@/api/bankAccountsApi';

export function useBankAccountsService() {
    async function searchBankAccounts(filter, pageSize) {
        const res = await searchBankAccountsFromAPI(filter, pageSize)
        if (!res.errorMessage) {
            return res.data.items
        } else {
            return res
        }
    }

    return {
        searchBankAccounts,
    }
}