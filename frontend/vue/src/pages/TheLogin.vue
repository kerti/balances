<script setup>
import { useAuthStore } from "@/stores/authStore"
import { ref } from "vue"
import { useRouter } from "vue-router"

const authStore = useAuthStore()
const router = useRouter()

const username = ref("")
const password = ref("")
const errorMessage = ref("")

const updateUsername = (event) => {
  username.value = event.target.value
}

const updatePassword = (event) => {
  password.value = event.target.value
}

async function authenticate(event) {
  event.preventDefault()

  const result = await authStore.authenticate(username.value, password.value)

  const success = !result.error

  if (success) {
    errorMessage.value = ""
    router.push("/")
  } else {
    if (result.error.message.includes("Network Error")) {
      errorMessage.value =
        "Our systems are experiencing disruptions, please try again later."
    } else {
      errorMessage.value = "Login failed. Please check your credentials."
    }
  }
}

const version = __APP_VERSION__
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-base-100">
    <div class="grid grid-cols-2 p-2 gap-3">
      <div
        class="card bg-base-100 shadow-sm w-96 flex items-center justify-center"
      >
        <div>
          <h1 class="text-3xl">Balances</h1>
          <p class="pt-2">Your personal finance, simplified.</p>
          <p class="text-xs opacity-50 pt-15">version {{ version }}</p>
        </div>
      </div>
      <form v-on:submit="authenticate" data-testid="frmLogin_form">
        <fieldset
          class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4"
        >
          <legend class="fieldset-legend">Login</legend>

          <label class="label" for="username">Username*</label>
          <input
            type="text"
            class="input"
            placeholder="username"
            autocomplete="username"
            id="username"
            data-testid="frmLogin_txtUsername"
            :value="username"
            @input="updateUsername"
          />

          <label class="label" for="password">Password*</label>
          <input
            type="password"
            class="input"
            placeholder="password"
            autocomplete="current-password"
            id="password"
            data-testid="frmLogin_txtPassword"
            :value="password"
            @input="updatePassword"
          />
          <button
            type="submit"
            class="btn btn-primary mt-4"
            data-testid="frmLogin_btnLogin"
          >
            Login
          </button>

          <div
            v-if="errorMessage"
            class="text-secondary-content mt-2 font-bold"
            data-testid="frmLogin_divErrorMessage"
          >
            {{ errorMessage }}
          </div>
        </fieldset>
      </form>
    </div>
  </div>
</template>