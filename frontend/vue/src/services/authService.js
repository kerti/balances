import { authenticateFromAPI, refreshTokenFromAPI } from '@/api/authApi'
import { useAuthCookie } from '@/composables/useAuthCookie'

const {
    setAuthTokenToCookie,
    getAuthTokenFromCookie,
    removeAuthTokenFromCookie,
    setUserDataToCookie,
    getUserDataFromCookie,
    removeUserDataFromCookie,
} = useAuthCookie()

export function useAuthService() {
    async function authenticate(username, password) {
        const result = await authenticateFromAPI(username, password)
        if (!result.error) {
            setAuthTokenToCookie(result.data.token)
            setUserDataToCookie(result.data.user)
        } else {
            deauthenticate()
        }
        return result
    }

    async function refreshToken() {
        const result = await refreshTokenFromAPI()
        if (!result.error) {
            setAuthTokenToCookie(result.data.token)
        } else {
            deauthenticate()
        }
        return result
    }

    function deauthenticate() {
        removeAuthTokenFromCookie()
        removeUserDataFromCookie()
    }

    function isLoggedIn() {
        const token = getAuthTokenFromCookie()
        return token !== undefined
    }

    function getUserData() {
        return getUserDataFromCookie()
    }

    return {
        authenticate,
        refreshToken,
        deauthenticate,
        isLoggedIn,
        getUserData,
    }
}