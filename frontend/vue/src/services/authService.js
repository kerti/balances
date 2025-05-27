import { authenticateFromAPI } from "@/api/authApi";

export function useAuthService() {
    async function authenticate(username, password) {
        const result = await authenticateFromAPI(username, password)
        // TODO: parse errors and return user object instead?
        return result
    }

    return {
        authenticate,
    }
}