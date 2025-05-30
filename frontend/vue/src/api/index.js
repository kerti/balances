import axios from 'axios'
import { useAuthCookie } from '@/composables/useAuthCookie'
import { useEnvUtils } from '@/composables/useEnvUtils'
import { useRouter } from 'vue-router'

const { getAuthTokenFromCookie, removeAuthTokenFromCookie, removeUserDataFromCookie } = useAuthCookie()
const ev = useEnvUtils()
const router = useRouter()

const axiosInstance = axios.create({
    baseURL: ev.getAPIBaseURL(),
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
            removeAuthTokenFromCookie()
            removeUserDataFromCookie()
            router.push({ name: 'login' })
        }
        // show error message
        return Promise.reject(error)
    }
)

export default axiosInstance