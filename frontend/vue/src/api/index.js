import axios from 'axios'
import { useAuthCookie } from '@/composables/useAuthCookie'
import { useEnvUtils } from '@/composables/useEnvUtils'

const { getAuthTokenFromCookie, removeAuthTokenFromCookie, removeUserDataFromCookie } = useAuthCookie()
const ev = useEnvUtils()

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
        // logout immediately when unauthorized access detected
        if (error.response?.status === 401) {
            removeAuthTokenFromCookie()
            removeUserDataFromCookie()
        }
        // refine the error object if possible and return it
        if (error.response?.data.error) {
            const respData = error.response?.data.error
            let errObj = {}
            errObj.code = respData.code
            if (respData.operation) {
                errObj.operation = respData.operation
            }
            if (respData.entity) {
                errObj.entity = respData.entity
            }
            errObj.message = respData.message
            return Promise.reject(errObj)
        } else {
            return Promise.reject(error)
        }
    }
)

export default axiosInstance