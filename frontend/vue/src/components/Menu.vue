<script setup></script>

<template>
  <ul class="menu w-full px-4 py-0">
    <!-- dashboard section -->
    <li>
      <h2 class="menu-title">
        <router-link to="/dashboard">Dashboard</router-link>
      </h2>
    </li>
    <!-- assets section -->
    <li>
      <details open>
        <summary>
          <h2 class="menu-title">Assets</h2>
        </summary>
        <ul>
          <li>
            <router-link to="/assets/bank-accounts">Bank Accounts</router-link>
          </li>
          <li>
            <router-link to="/assets/vehicles">Vehicles</router-link>
          </li>
          <li>
            <router-link to="/assets/properties">Properties</router-link>
          </li>
        </ul>
      </details>
    </li>
    <!-- liabilities section -->
    <li>
      <details open>
        <summary>
          <h2 class="menu-title">Liabilities</h2>
        </summary>
        <ul>
          <li>
            <router-link to="/liabilities/personal-debts"
              >Personal Debts</router-link
            >
          </li>
          <li>
            <router-link to="/liabilities/institutional-debts"
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
            <router-link to="/investments/deposits">Deposits</router-link>
          </li>
          <li>
            <router-link to="/investments/obligations">Obligations</router-link>
          </li>
          <li>
            <router-link to="/investments/gold">Gold</router-link>
          </li>
          <li>
            <router-link to="/investments/mutual-funds"
              >Mutual Funds</router-link
            >
          </li>
          <li>
            <router-link to="/investments/stocks">Stocks</router-link>
          </li>
          <li>
            <router-link to="/investments/p2p-lendings"
              >P2P Lendings</router-link
            >
          </li>
        </ul>
      </details>
    </li>
    <li></li>
    <li>
      <router-link to="/about">About</router-link>
    </li>
    <li>
      <router-link to="/docs">Documentation</router-link>
    </li>
  </ul>
</template>