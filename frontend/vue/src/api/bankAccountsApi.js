import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

//// bank accounts CRUD

// create
export async function createBankAccountWithAPI(account) {
    try {
        const { data } = await axiosInstance.post('bankAccounts', {
            accountName: account.accountName,
            bankName: account.bankName,
            accountHolderName: account.accountHolderName,
            accountNumber: account.accountNumber,
            lastBalance: account.lastBalance,
            lastBalanceDate: account.lastBalanceDate,
            status: account.status,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function searchBankAccountsFromAPI(filter, pageSize) {
    try {
        const { data } = await axiosInstance.post('bankAccounts/search', {
            keyword: filter ? filter : '',
            pageSize: getPageSize(pageSize)
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function getBankAccountFromAPI(bankAccountId, startDate, endDate, pageSize, page) {
    try {
        const params = new URLSearchParams()
        params.append('withBalances', true)
        if (startDate) {
            params.append('balanceStartDate', startDate)
        }
        if (endDate) {
            params.append('balanceEndDate', endDate)
        }
        if (pageSize) {
            params.append('pageSize', pageSize)
        }
        const { data } = await axiosInstance.get('bankAccounts/' + bankAccountId + '?' + params.toString())
        return data
    } catch (errror) {
        return {
            error: error
        }
    }
}

// update
export async function updateBankAccountWithAPI(account) {
    try {
        const { data } = await axiosInstance.patch('bankAccounts/' + account.id, account)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deleteBankAccountWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('bankAccounts/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

//// bank account balances CRUD

// create
export async function createBankAccountBalanceWithAPI(accountBalance) {
    try {
        const { data } = await axiosInstance.post('bankAccounts/balances', {
            bankAccountId: accountBalance.bankAccountId,
            date: accountBalance.date,
            balance: accountBalance.balance,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function searchBankAccountBalancesFromAPI(bankAccountIds, startDate, endDate, pageSize, page) {
    try {
        const { data } = await axiosInstance.post('bankAccounts/balances/search', {
            bankAccountIds: bankAccountIds,
            startDate: startDate ? startDate : null,
            endDate: endDate ? endDate : null,
            pageSize: getPageSize(pageSize),
            page: page ? page : 1,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function getBankAccountBalanceFromAPI(id) {
    try {
        const { data } = await axiosInstance.get('bankAccounts/balances/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// update
export async function updateAccountBalanceWithAPI(accountBalance) {
    try {
        const { data } = await axiosInstance.patch('bankAccounts/balances/' + accountBalance.id, accountBalance)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deleteBankAccountBalanceWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('bankAccounts/balances/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

//// utilities

function getPageSize(pageSize) {
    const ev = useEnvUtils()
    return pageSize ? pageSize : ev.getDefaultPageSize()
}