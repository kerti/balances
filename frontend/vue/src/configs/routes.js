import { createRouter, createWebHistory } from 'vue-router'
import { useAuthCookie } from '@/composables/useAuthCookie'

import Dashboard from '@/pages/Dashboard.vue'
import TheLogin from '@/pages/TheLogin.vue'

const routes = [
    // dashboard
    {
        path: '/',
        name: 'base',
        component: Dashboard,
        meta: { requiresAuth: true, pageTitle: 'Dashboard' }
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        component: Dashboard,
        meta: { requiresAuth: true, pageTitle: 'Dashboard' }
    },
    // assets
    {
        path: '/assets/bank-accounts',
        name: 'assets.bankaccounts',
        component: () => import('@/pages/assets/BankAccounts.vue'),
        meta: { requiresAuth: true, pageTitle: 'Bank Accounts' }
    },
    {
        path: '/assets/bank-accounts/:id',
        name: 'assets.bankaccount.detail',
        component: () => import('@/pages/assets/BankAccountDetail.vue'),
        meta: { requiresAuth: true, pageTitle: 'Bank Account Details' }
    },
    {
        path: '/assets/vehicles',
        name: 'assets.vehicles',
        component: () => import('@/pages/assets/Vehicles.vue'),
        meta: { requiresAuth: true, pageTitle: 'Vehicles' }
    },
    {
        path: '/assets/properties',
        name: 'assets.properties',
        component: () => import('@/pages/assets/Properties.vue'),
        meta: { requiresAuth: true, pageTitle: 'Properties' }
    },
    // liabilities
    {
        path: '/liabilities/personal-debts',
        name: 'liabilities.personaldebts',
        component: () => import('@/pages/liabilities/PersonalDebts.vue'),
        meta: { requiresAuth: true, pageTitle: 'Personal Debts' }
    },
    {
        path: '/liabilities/institutional-debts',
        name: 'liabilities.institutionaldebts',
        component: () => import('@/pages/liabilities/InstitutionalDebts.vue'),
        meta: { requiresAuth: true, pageTitle: 'Institutional Debts' }
    },
    // investments
    {
        path: '/investments/deposits',
        name: 'investments.deposits',
        component: () => import('@/pages/investments/Deposits.vue'),
        meta: { requiresAuth: true, pageTitle: 'Deposits' }
    },
    {
        path: '/investments/obligations',
        name: 'investments.obligations',
        component: () => import('@/pages/investments/Obligations.vue'),
        meta: { requiresAuth: true, pageTitle: 'Obligations' }
    },
    {
        path: '/investments/gold',
        name: 'investments.gold',
        component: () => import('@/pages/investments/Gold.vue'),
        meta: { requiresAuth: true, pageTitle: 'Gold' }
    },
    {
        path: '/investments/mutual-funds',
        name: 'investments.mutualfunds',
        component: () => import('@/pages/investments/MutualFunds.vue'),
        meta: { requiresAuth: true, pageTitle: 'Mutual Funds' }
    },
    {
        path: '/investments/stocks',
        name: 'investments.stocks',
        component: () => import('@/pages/investments/Stocks.vue'),
        meta: { requiresAuth: true, pageTitle: 'Stocks' }
    },
    {
        path: '/investments/p2p-lendings',
        name: 'investments.p2plendings',
        component: () => import('@/pages/investments/P2PLendings.vue'),
        meta: { requiresAuth: true, pageTitle: 'P2P Lendings' }
    },
    // authentication and authorization
    {
        path: '/login',
        name: 'login',
        component: TheLogin,
        meta: { requiresAuth: false, pageTitle: 'Login' }
    },
    // miscellaneous
    {
        path: '/about',
        name: 'misc.about',
        component: () => import('@/pages/About.vue'),
        meta: { requiresAuth: false, pageTitle: 'About' }
    },
    {
        path: '/docs',
        name: 'misc.docs',
        component: () => import('@/pages/Docs.vue'),
        meta: { requiresAuth: false, pageTitle: 'Documentation' }
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    linkActiveClass: 'menu-active'
})

router.beforeEach((to, from, next) => {
    const { getAuthTokenFromCookie } = useAuthCookie()
    const isAuthenticated = !!getAuthTokenFromCookie()

    if (to.meta.requiresAuth && !isAuthenticated) {
        next({ name: 'login' })
    } else {
        next()
    }
})

export default router