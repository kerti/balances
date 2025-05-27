import { authenticateFromAPI, refreshTokenFromAPI } from "@/api/authApi";

export function useAuthService() {
    async function authenticate(username, password) {
        const result = await authenticateFromAPI(username, password)
        return result
    }

    async function refreshToken(currentToken) {
        const result = await refreshTokenFromAPI(currentToken)
        return result
    }

    return {
        authenticate,
        refreshToken,
    }
}