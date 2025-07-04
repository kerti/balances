import {
    // bank accounts
    createBankAccountWithAPI,
    searchBankAccountsFromAPI,
    getBankAccountFromAPI,
    updateBankAccountWithAPI,
    deleteBankAccountWithAPI,
    // bank account balances
    createBankAccountBalanceWithAPI,
    searchBankAccountBalancesFromAPI,
    getBankAccountBalanceFromAPI,
    updateAccountBalanceWithAPI,
    deleteBankAccountBalanceWithAPI,
} from '@/api/bankAccountsApi'

export function useBankAccountsService() {

    //// bank accounts CRUD

    // create
    async function createBankAccount(account) {
        const payload = {
            accountName: account.accountName,
            bankName: account.bankName,
            accountHolderName: account.accountHolderName,
            accountNumber: account.accountNumber,
            lastBalance: parseFloat(account.lastBalance),
            lastBalanceDate: account.lastBalanceDate,
            status: account.status,
        }

        const result = await createBankAccountWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function searchBankAccounts(filter, balancesStartDate, balancesEndDate, pageSize) {
        // get the account data
        const accounts = await searchBankAccountsFromAPI(filter, pageSize)

        if (!accounts.error) {
            // then get the balances data
            const bankAccountIds = accounts.data.items.map(account => account.id)
            const balances = await searchBankAccountBalancesFromAPI(bankAccountIds, balancesStartDate, balancesEndDate, 99999, 1)
            if (!balances.error) {
                return accounts.data.items.map(account => {
                    account.balances = balances.data.items.filter(function (bal) {
                        return bal.bankAccountId == account.id
                    }).sort((a, b) => a.date - b.date)
                    return account
                })
            } else {
                return {
                    error: balances.error
                }
            }

        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function getBankAccount(id, balanceStartDate, balanceEndDate, pageSize) {
        const account = await getBankAccountFromAPI(id, balanceStartDate, balanceEndDate, pageSize)

        if (!account.error) {
            account.data.balances.sort((a, b) => a.date - b.date)
            return account.data
        } else {
            return {
                error: account.error
            }
        }
    }

    // update
    async function updateBankAccount(account) {
        const payload = {
            id: account.id,
            accountName: account.accountName,
            bankName: account.bankName,
            accountHolderName: account.accountHolderName,
            accountNumber: account.accountNumber,
            status: account.status,
        }

        const result = await updateBankAccountWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // delete
    async function deleteBankAccount(id) {
        const result = await deleteBankAccountWithAPI(id)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    //// bank account balances CRUD

    // create
    async function createBankAccountBalance(accountBalance) {
        const payload = {
            bankAccountId: accountBalance.bankAccountId,
            date: accountBalance.date,
            balance: parseInt(accountBalance.balance)
        }

        const result = await createBankAccountBalanceWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function getBankAccountBalance(id) {
        const accountBalance = await getBankAccountBalanceFromAPI(id)

        if (!accountBalance.error) {
            return accountBalance.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // update
    async function updateBankAccountBalance(accountBalance) {
        const payload = {
            id: accountBalance.id,
            bankAccountId: accountBalance.bankAccountId,
            date: accountBalance.date,
            balance: parseInt(accountBalance.balance),
        }

        const result = await updateAccountBalanceWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // delete
    async function deleteBankAccountBalance(id) {
        const result = await deleteBankAccountBalanceWithAPI(id)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    return {
        // bank accounts
        createBankAccount,
        searchBankAccounts,
        getBankAccount,
        updateBankAccount,
        deleteBankAccount,
        // bank account balances
        createBankAccountBalance,
        getBankAccountBalance,
        updateBankAccountBalance,
        deleteBankAccountBalance,
    }
}