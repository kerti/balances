<script setup>
import Sidebar from "./components/Sidebar.vue"
import Header from "./components/Header.vue"
</script>

<template>
  <div class="flex h-screen">
    <Sidebar />
    <div class="flex flex-col flex-1">
      <Header />
      <main class="flex-1 overflow-y-auto bg-base-200 p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped></style>
