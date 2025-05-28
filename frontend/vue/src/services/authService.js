import { authenticateFromAPI, refreshTokenFromAPI } from '@/api/authApi'
import { useAuthCookie } from '@/composables/useAuthCookie'

const {
    setAuthTokenToCookie,
    removeAuthTokenFromCookie,
    setUserDataToCookie,
    removeUserDataFromCookie,
} = useAuthCookie()

export function useAuthService() {
    async function authenticate(username, password) {
        const result = await authenticateFromAPI(username, password)
        if (!result.errorMessage) {
            setAuthTokenToCookie(result.data.token)
            setUserDataToCookie(result.data.user)
        } else {
            deauthenticate()
        }
        return result
    }

    async function refreshToken() {
        const result = await refreshTokenFromAPI()
        if (!result.errorMessage) {
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

    return {
        authenticate,
        refreshToken,
        deauthenticate,
    }
}