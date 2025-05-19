<script setup>
import TheHeader from "@/components/TheHeader.vue"
import TheSidebar from "@/components/TheSidebar.vue"
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
