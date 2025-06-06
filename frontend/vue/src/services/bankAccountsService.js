import {
    searchBankAccountBalancesFromAPI,
    searchBankAccountsFromAPI,
    getBankAccountFromAPI,
    updateAccountWithAPI,
    getBankAccountBalanceFromAPI,
    updateAccountBalanceWithAPI,
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

    async function getBankAccount(id, balanceStartDate, balanceEndDate, pageSize) {
        const account = await getBankAccountFromAPI(id, balanceStartDate, balanceEndDate, pageSize)

        if (!account.errorMessage) {
            account.data.balances.sort((a, b) => a.date - b.date)
            return account.data
        } else {
            return {
                errorMessage: account.errorMessage
            }
        }
    }

    async function updateBankAccount(account) {
        const payload = {
            id: account.id,
            accountName: account.accountName,
            bankName: account.bankName,
            accountHolderName: account.accountHolderName,
            accountNumber: account.accountNumber,
            status: account.status,
        }

        const result = await updateAccountWithAPI(payload)

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    async function getBankAccountBalance(id) {
        const accountBalance = await getBankAccountBalanceFromAPI(id)

        if (!accountBalance.errorMessage) {
            return accountBalance.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    async function updateBankAccountBalance(accountBalance) {
        const payload = {
            id: accountBalance.id,
            bankAccountId: accountBalance.bankAccountId,
            date: accountBalance.date,
            balance: parseInt(accountBalance.balance),
        }

        const result = await updateAccountBalanceWithAPI(payload)

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    return {
        searchBankAccounts,
        getBankAccount,
        updateBankAccount,
        getBankAccountBalance,
        updateBankAccountBalance,
    }
}