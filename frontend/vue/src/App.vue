<script setup>
import TheHeader from "@/components/TheHeader.vue"
import TheSidebar from "@/components/TheSidebar.vue"
import { useAuthStore } from "@/stores/authStore"
import { ref } from "vue"

const authStore = useAuthStore()
const usernameValue = ref("")
const passwordValue = ref("")

const updateUsernameValue = (event) => {
  usernameValue.value = event.target.value
}

const updatePasswordValue = (event) => {
  passwordValue.value = event.target.value
}

function authenticate() {
  authStore.authenticate(usernameValue.value, passwordValue.value)
}
</script>

<template>
  <div class="flex h-screen">
    <TheSidebar />
    <div class="flex flex-col flex-1">
      <TheHeader />
      <main class="flex-1 overflow-y-auto bg-base-200 p-4">
        <router-view />
        <fieldset
          class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4"
        >
          <legend class="fieldset-legend">Login</legend>

          <label class="label">Username</label>
          <input
            type="text"
            class="input"
            placeholder="Username"
            :value="usernameValue"
            @input="updateUsernameValue"
          />

          <label class="label">Password</label>
          <input
            type="password"
            class="input"
            placeholder="Password"
            :value="passwordValue"
            @input="updatePasswordValue"
          />
          <div class="grid grid-cols-3 gap-2">
            <button class="btn btn-primary mt-4" @click="authenticate()">
              Login
            </button>
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
        </fieldset>
        <div>Authentication Response: {{ authStore.isLoggedIn }}</div>
      </main>
    </div>
  </div>
</template>

<style scoped></style>
