import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

function getPageSize(pageSize) {
    const ev = useEnvUtils()
    return pageSize ? pageSize : ev.getDefaultPageSize()
}

export async function searchBankAccountsFromAPI(filter, pageSize) {
    try {
        const ev = useEnvUtils()
        const { data } = await axiosInstance.post('bankAccounts/search', {
            keyword: filter ? filter : '',
            pageSize: getPageSize(pageSize)
        })
        return data
    } catch (error) {
        return {
            errorMessage: 'API - ' + error.message
        }
    }
}

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
            errorMessage: 'API - ' + error.message
        }
    }
}