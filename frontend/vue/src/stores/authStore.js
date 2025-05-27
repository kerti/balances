import { useAuthService } from "@/services/authService";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useAuthStore = defineStore('auth', () => {
    const authService = useAuthService()

    // reactive state
    const isLoggedIn = ref(false)
    const username = ref('')
    const name = ref('')
    const email = ref('')

    // actions
    async function authenticate(uname, password) {
        const authenticationResult = await authService.authenticate(uname, password)
        const user = authenticationResult.data.user
        isLoggedIn.value = true
        username.value = user.username
        name.value = user.name
        email.value = user.email
    }

    return {
        isLoggedIn,
        username,
        name,
        email,
        authenticate,
    }
})