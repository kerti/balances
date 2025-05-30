import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

export async function searchBankAccountsFromAPI(filter, pageSize) {
    try {
        const ev = useEnvUtils()
        const { data } = await axiosInstance.post('bankAccounts/search', {
            keyword: filter ? filter : '',
            pageSize: pageSize ? pageSize : ev.getDefaultPageSize()
        })
        return data
    } catch (error) {
        return {
            errorMessage: 'API - ' + error.message
        }
    }
}