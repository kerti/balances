import axiosInstance from '@/api/index'

export async function searchBankAccountsFromAPI(filter, pageSize) {
    try {
        const { data } = await axiosInstance.post('bankAccounts/search', {
            keyword: filter ? filter : '',
            pageSize: pageSize ? pageSize : 10
        })
        return data
    } catch (error) {
        return {
            errorMessage: 'API - ' + error.message
        }
    }
}