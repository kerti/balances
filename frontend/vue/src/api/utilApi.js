import axiosInstance from '@/api/index'

export async function getServerHealth() {
    const { data } = await axiosInstance.get('/health')
    return data
}