<script setup>
const version = __APP_VERSION__
</script>

<template>
  <footer
    class="footer justify-between items-center sm:footer-horizontal bg-base-100 text-base-content px-10 py-1"
  >
    <aside class="flex-1 text-left">
      <p>Disclaimer</p>
    </aside>
    <aside class="flex-1 text-center">
      <div>&copy; 2025 Raditya Kertiyasa</div>
    </aside>
    <aside class="flex-1 text-right">
      <div>v{{ version }}</div>
    </aside>
  </footer>
</template>

<style scoped>
</style>