import axiosInstance from "@/api/index";

export async function authenticateFromAPI(username, password) {
    try {
        const basicAuth = btoa(`${username}:${password}`)
        const { data } = await axiosInstance.post('auth/login', null, {
            headers: {
                Authorization: `Basic ${basicAuth}`,
            },
        })
        return data
    } catch (error) {
        return {
            errorMessage: 'API - ' + error.message
        }
    }
}

export async function refreshTokenFromAPI(currentToken) {
    try {
        const { data } = await axiosInstance.get('auth/token', {
            headers: {
                Authorization: `Bearer ${currentToken}`,
            },
        })
        return data
    } catch (error) {
        return {
            errorMessage: 'API - ' + error.message
        }
    }
}