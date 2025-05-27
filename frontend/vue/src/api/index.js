import axios from 'axios'
import { useAuthCookie } from '@/composables/useAuthCookie'

const { getAuthTokenFromCookie, removeAuthTokenFromCookie } = useAuthCookie()


const axiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
})

axiosInstance.interceptors.request.use((config) => {
    if (!config.headers.Authorization) {
        const token = getAuthTokenFromCookie()
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
    }
    return config
})

axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // handle logout or redirect to login
            removeAuthTokenFromCookie()
            console.error('unauthorized, logging out...')
        }
        // show error message
        return Promise.reject(error)
    }
)

export default axiosInstance