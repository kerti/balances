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
} from '@/api/bankAccountsApi';

export function useBankAccountsService() {

    //// bank accounts CRUD

    // create
    async function createBankAccount(account) {
        const payload = {
            accountName: account.accountName,
            bankName: account.bankName,
            accountHolderName: account.accountHolderName,
            accountNumber: account.accountNumber,
            lastBalance: parseInt(account.lastBalance),
            lastBalanceDate: account.lastBalanceDate,
            status: account.status,
        }

        const result = await createBankAccountWithAPI(payload)

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    // read
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

    // read
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

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    // delete
    async function deleteBankAccount(id) {
        const result = await deleteBankAccountWithAPI(id)

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
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

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
            }
        }
    }

    // read
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

    // update
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

    // delete
    async function deleteBankAccountBalance(id) {
        const result = await deleteBankAccountBalanceWithAPI(id)

        if (!result.errorMessage) {
            return result.data
        } else {
            return {
                errorMessage: result.errorMessage
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