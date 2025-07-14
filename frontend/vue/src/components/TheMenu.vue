<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/authStore"

const authStore = useAuthStore()
const router = useRouter()

function handleLogout() {
  authStore.deauthenticate()
  router.push({ name: "login" })
}
</script>

<template>
  <ul class="menu w-full px-4 py-0">
    <!-- dashboard section -->
    <li>
      <router-link :to="{ name: 'dashboard' }">Dashboard</router-link>
    </li>
    <!-- assets section -->
    <li>
      <router-link :to="{ name: 'assets' }">Assets</router-link>
      <ul>
        <li>
          <router-link :to="{ name: 'assets.bankaccounts' }"
            >Bank Accounts</router-link
          >
        </li>
        <li>
          <router-link :to="{ name: 'assets.vehicles' }">Vehicles</router-link>
        </li>
        <li>
          <router-link :to="{ name: 'assets.properties' }"
            >Properties</router-link
          >
        </li>
      </ul>
    </li>
    <!-- liabilities section -->
    <li>
      <details open>
        <summary>
          <h2 class="menu-title">Liabilities</h2>
        </summary>
        <ul>
          <li>
            <router-link :to="{ name: 'liabilities.personaldebts' }"
              >Personal Debts</router-link
            >
          </li>
          <li>
            <router-link :to="{ name: 'liabilities.institutionaldebts' }"
              >Institutional Debts</router-link
            >
          </li>
        </ul>
      </details>
    </li>
    <!-- investments section -->
    <li>
      <details open>
        <summary>
          <h2 class="menu-title">Investments</h2>
        </summary>
        <ul>
          <li>
            <router-link :to="{ name: 'investments.deposits' }"
              >Deposits</router-link
            >
          </li>
          <li>
            <router-link :to="{ name: 'investments.obligations' }"
              >Obligations</router-link
            >
          </li>
          <li>
            <router-link :to="{ name: 'investments.gold' }">Gold</router-link>
          </li>
          <li>
            <router-link :to="{ name: 'investments.mutualfunds' }"
              >Mutual Funds</router-link
            >
          </li>
          <li>
            <router-link :to="{ name: 'investments.stocks' }"
              >Stocks</router-link
            >
          </li>
          <li>
            <router-link :to="{ name: 'investments.p2plendings' }"
              >P2P Lendings</router-link
            >
          </li>
        </ul>
      </details>
    </li>
    <li></li>
    <li>
      <router-link :to="{ name: 'misc.about' }">About</router-link>
    </li>
    <li>
      <router-link :to="{ name: 'misc.docs' }">Documentation</router-link>
    </li>
    <li></li>
    <li>
      <a href="#" @click="handleLogout" data-testid="divSidebar_btnLogout"
        >Logout</a
      >
    </li>
  </ul>
</template>