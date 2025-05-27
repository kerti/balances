import { useAuthService } from "@/services/authService";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useAuthStore = defineStore('auth', () => {
    const authService = useAuthService()

    // reactive state
    const isLoggedIn = ref(false)
    const user = ref({})
    const token = ref('')

    // actions
    async function authenticate(uname, password) {
        const authenticationResult = await authService.authenticate(uname, password)
        isLoggedIn.value = true
        user.value = authenticationResult.data.user
        token.value = authenticationResult.data.token
    }

    function deauthenticate() {
        isLoggedIn.value = false
        user.value = {}
        token.value = ''
    }

    async function refreshToken() {
        if (token.value === '') {
            return
        }

        const refreshResult = await authService.refreshToken(token.value)
        if (refreshResult.errorMessage) {
            deauthenticate()
        } else {
            token.value = refreshResult.data.token
        }
    }

    return {
        isLoggedIn,
        user,
        token,
        authenticate,
        deauthenticate,
        refreshToken,
    }
})