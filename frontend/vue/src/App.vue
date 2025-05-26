<script setup>
import TheHeader from "@/components/TheHeader.vue"
import TheSidebar from "@/components/TheSidebar.vue"
import { useUiStore } from "./stores/ui"

const uiStore = useUiStore()
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
      </main>
    </div>
  </div>
</template>

<style scoped></style>
