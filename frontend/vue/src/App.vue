<script setup>
import TheHeader from "@/components/TheHeader.vue"
import TheSidebar from "@/components/TheSidebar.vue"
import TheFooter from "./components/TheFooter.vue"
import { useAuthStore } from "@/stores/authStore"

const authStore = useAuthStore()
</script>

<template>
  <div id="toast-container" class="toast toast-top toast-end fixed z-50"></div>
  <template v-if="authStore.isLoggedIn">
    <div class="flex h-screen">
      <TheSidebar />
      <div class="flex flex-col flex-1">
        <TheHeader />
        <main class="flex-1 overflow-y-auto overflow-x-hidden bg-base-200 p-4">
          <router-view />
        </main>
        <TheFooter />
      </div>
    </div>
  </template>
  <template v-else>
    <router-view />
  </template>
</template>

<style scoped></style>
