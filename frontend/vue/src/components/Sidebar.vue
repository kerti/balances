<script setup>
import Menu from "./Menu.vue"
</script>

<template>
  <div class="w-58 h-full bg-base-100 p-4 shadow-md">
    <div class="flex-1">
      <a class="btn btn-ghost text-xl"
        ><svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          class="inline-block h-5 w-5 stroke-current"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16"
          ></path></svg
        >Balances</a
      >
    </div>
    <div class="h-4"></div>
    <Menu />
  </div>
</template>
