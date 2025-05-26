<script setup>
import TheHeader from "@/components/TheHeader.vue"
import TheSidebar from "@/components/TheSidebar.vue"
import { useUiStore } from "@/stores/uiStore"
import { useUtilStore } from "@/stores/utilStore"

const uiStore = useUiStore()
const api = import.meta.env.VITE_API_BASE_URL

const utilStore = useUtilStore()
utilStore.getServerHealth()
</script>

<template>
  <div class="flex h-screen">
    <TheSidebar />
    <div class="flex flex-col flex-1">
      <TheHeader />
      <main class="flex-1 overflow-y-auto bg-base-200 p-4">
        <router-view />
        <div>Theme is {{ uiStore.theme }}.</div>
        <div>Light theme is {{ uiStore.lightTheme }}.</div>
        <div>Dark theme is {{ uiStore.darkTheme }}.</div>
        <div>API Endpoint is {{ api }}</div>
        <div>Server Health Response: {{ utilStore.serverIsHealthy }}</div>
      </main>
    </div>
  </div>
</template>

<style scoped></style>
