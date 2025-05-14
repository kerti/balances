<script setup>
import TheSidebar from "./components/TheSidebar.vue"
import TheHeader from "./components/TheHeader.vue"
</script>

<template>
  <div class="flex h-screen">
    <TheSidebar />
    <div class="flex flex-col flex-1">
      <TheHeader />
      <main class="flex-1 overflow-y-auto bg-base-200 p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped></style>
