import { createRouter, createWebHistory } from 'vue-router'

import Dashboard from '@/pages/Dashboard.vue'

import Login from '@/pages/Login.vue'

const routes = [
    // dashboard
    { path: '/', name: 'base', component: Dashboard },
    { path: '/dashboard', name: 'dashboard', component: Dashboard },
    // assets
    { path: '/assets/bank-accounts', name: 'assets.bankaccounts', component: () => import('@/pages/assets/BankAccounts.vue') },
    { path: '/assets/bank-accounts/:id', name: 'assets.bankaccount.detail', component: () => import('@/pages/assets/BankAccountDetail.vue') },
    { path: '/assets/vehicles', name: 'assets.vehicles', component: () => import('@/pages/assets/Vehicles.vue') },
    { path: '/assets/properties', name: 'assets.properties', component: () => import('@/pages/assets/Properties.vue') },
    // liabilities
    { path: '/liabilities/personal-debts', name: 'liabilities.personaldebts', component: () => import('@/pages/liabilities/PersonalDebts.vue') },
    { path: '/liabilities/institutional-debts', name: 'liabilities.institutionaldebts', component: () => import('@/pages/liabilities/InstitutionalDebts.vue') },
    // investments
    { path: '/investments/deposits', name: 'investments.deposits', component: () => import('@/pages/investments/Deposits.vue') },
    { path: '/investments/obligations', name: 'investments.obligations', component: () => import('@/pages/investments/Obligations.vue') },
    { path: '/investments/gold', name: 'investments.gold', component: () => import('@/pages/investments/Gold.vue') },
    { path: '/investments/mutual-funds', name: 'investments.mutualfunds', component: () => import('@/pages/investments/MutualFunds.vue') },
    { path: '/investments/stocks', name: 'investments.stocks', component: () => import('@/pages/investments/Stocks.vue') },
    { path: '/investments/p2p-lendings', name: 'investments.p2plendings', component: () => import('@/pages/investments/P2PLendings.vue') },
    // authentication and authorization
    { path: '/login', name: 'login', component: Login },
    // miscellaneous
    { path: '/about', name: 'misc.about', component: () => import('@/pages/About.vue') },
    { path: '/docs', name: 'misc.docs', component: () => import('@/pages/Docs.vue') },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    linkActiveClass: 'menu-active'
})

export default router