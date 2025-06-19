import axiosInstance from '@/api/index'

export async function getServerHealthFromAPI() {
    try {
        const { data } = await axiosInstance.get('/health')
        return data
    } catch (error) {
        return {
            error: error
        }
    }

}