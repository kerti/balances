import { createRouter, createWebHistory } from 'vue-router'

import Dashboard from '../pages/Dashboard.vue'

import Login from '../pages/Login.vue'

const routes = [
    // dashboard
    { path: '/dashboard', name: 'Dashboard', component: Dashboard },
    // assets
    { path: '/assets/bank-accounts', name: 'Bank Accounts', component: () => import('../pages/assets/BankAccounts.vue') },
    { path: '/assets/vehicles', name: 'Vehicles', component: () => import('../pages/assets/Vehicles.vue') },
    { path: '/assets/properties', name: 'Properties', component: () => import('../pages/assets/Properties.vue') },
    // liabilities
    { path: '/liabilities/personal-debts', name: 'Personal Debts', component: () => import('../pages/liabilities/PersonalDebts.vue') },
    { path: '/liabilities/institutional-debts', name: 'Institutional Debts', component: () => import('../pages/liabilities/InstitutionalDebts.vue') },
    // investments
    { path: '/investments/deposits', name: 'Deposits', component: () => import('../pages/investments/Deposits.vue') },
    { path: '/investments/obligations', name: 'Obligations', component: () => import('../pages/investments/Obligations.vue') },
    { path: '/investments/gold', name: 'Gold', component: () => import('../pages/investments/Gold.vue') },
    { path: '/investments/mutual-funds', name: 'Mutual Funds', component: () => import('../pages/investments/MutualFunds.vue') },
    { path: '/investments/stocks', name: 'Stocks', component: () => import('../pages/investments/Stocks.vue') },
    { path: '/investments/p2p-lendings', name: 'P2P Lendings', component: () => import('../pages/investments/P2PLendings.vue') },
    // authentication and authorization
    { path: '/login', name: 'Login', component: Login },
    // miscellaneous
    { path: '/about', name: 'About', component: () => import('../pages/About.vue') },
    { path: '/docs', name: 'Documentation', component: () => import('../pages/Docs.vue') },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router