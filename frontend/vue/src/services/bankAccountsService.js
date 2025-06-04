import {
    searchBankAccountBalancesFromAPI,
    searchBankAccountsFromAPI,
    getBankAccountFromAPI,
} from '@/api/bankAccountsApi';

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
                    })
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

    async function getBankAccount(id, balanceStartDate, balanceEndDate, pageSize) {
        const account = await getBankAccountFromAPI(id, balanceStartDate, balanceEndDate, pageSize)

        if (!account.errorMessage) {
            return account.data
        } else {
            return {
                errorMessage: account.errorMessage
            }
        }
    }

    return {
        searchBankAccounts,
        getBankAccount,
    }
}