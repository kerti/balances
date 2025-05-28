import { useAuthService } from '@/services/authService'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const authService = useAuthService()

    // reactive state
    const isLoggedIn = ref(false)
    const user = ref({})

    // actions
    async function authenticate(uname, password) {
        const authenticationResult = await authService.authenticate(uname, password)
        if (!authenticationResult.errorMessage) {
            isLoggedIn.value = true
            user.value = authenticationResult.data.user
        }
        return authenticationResult
    }

    function deauthenticate() {
        authService.deauthenticate()
        isLoggedIn.value = false
        user.value = {}
    }

    async function refreshToken() {
        const refreshResult = await authService.refreshToken()
        if (refreshResult.errorMessage) {
            deauthenticate()
        }
    }

    return {
        isLoggedIn,
        user,
        authenticate,
        deauthenticate,
        refreshToken,
    }
})