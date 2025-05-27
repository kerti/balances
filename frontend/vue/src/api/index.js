import axios from "axios"

const axiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
})

axiosInstance.interceptors.request.use((config) => {
    // add bearer token or other headers
    // TODO: use cookies instead?
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = 'Bearer ${token}'
    }
    return config
})

axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // handle logout or redirect to login
            console.error('unauthorized, logging out...')
        }
        // show error message
        return Promise.reject(error)
    }
)

export default axiosInstance