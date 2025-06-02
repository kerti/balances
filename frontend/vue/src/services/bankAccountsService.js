import { searchBankAccountBalancesFromAPI, searchBankAccountsFromAPI } from '@/api/bankAccountsApi';

export function useBankAccountsService() {
    async function searchBankAccounts(filter, balancesStartDate, balancesEndDate, pageSize) {
        // get the account data
        const accounts = await searchBankAccountsFromAPI(filter, pageSize)

        if (!accounts.errorMessage) {

            // then get the balances data
            const bankAccountIds = accounts.data.items.map(account => account.id)
            const balances = await searchBankAccountBalancesFromAPI(bankAccountIds, balancesStartDate, balancesEndDate, 99999, 1)
            if (!balances.errorMessage) {
                return accounts.data.items.map(account => {
                    account.balances = balances.data.items.filter(function (bal) {
                        return bal.bankAccountId == account.id
                    }).sort((a, b) => a.date - b.date)
                    return account
                })
            } else {
                return {
                    errorMessage: balances.errorMessage,
                }
            }

        } else {
            return {
                errorMessage: accounts.errorMessage,
            }
        }
    }

    return {
        searchBankAccounts,
    }
}