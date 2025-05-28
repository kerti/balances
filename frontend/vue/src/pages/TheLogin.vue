<script setup>
import { useAuthStore } from "@/stores/authStore"
import { ref } from "vue"
import { useRouter } from "vue-router"

const authStore = useAuthStore()
const router = useRouter()

const usernameValue = ref("")
const passwordValue = ref("")
const errorMessage = ref("")

const updateUsernameValue = (event) => {
  usernameValue.value = event.target.value
}

const updatePasswordValue = (event) => {
  passwordValue.value = event.target.value
}

async function authenticate(event) {
  event.preventDefault()

  const result = await authStore.authenticate(
    usernameValue.value,
    passwordValue.value
  )

  const success = !result.errorMessage

  if (success) {
    errorMessage.value = ""
    router.push("/")
  } else {
    errorMessage.value = "Login failed. Please check your credentials." // ✅ Show error on failed login
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-base-100">
    <form v-on:submit="authenticate">
      <fieldset
        class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4"
      >
        <legend class="fieldset-legend">Login</legend>

        <label class="label">Username</label>
        <input
          type="text"
          class="input"
          placeholder="Username"
          autocomplete="username"
          :value="usernameValue"
          @input="updateUsernameValue"
        />

        <label class="label">Password</label>
        <input
          type="password"
          class="input"
          placeholder="Password"
          autocomplete="current-password"
          :value="passwordValue"
          @input="updatePasswordValue"
        />
        <div class="grid grid-cols-3 gap-2">
          <button type="submit" class="btn btn-primary mt-4">Login</button>
          <button
            class="btn btn-secondary mt-4"
            @click="authStore.deauthenticate()"
          >
            Logout
          </button>
          <button
            class="btn btn-secondary mt-4"
            @click="authStore.refreshToken()"
          >
            Refresh Token
          </button>
        </div>

        <div v-if="errorMessage" class="text-error mt-2">
          {{ errorMessage }}
        </div>

        <div>Authentication Response: {{ authStore.isLoggedIn }}</div>
      </fieldset>
    </form>
  </div>
</template>