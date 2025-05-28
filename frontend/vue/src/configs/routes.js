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
        meta: { requiresAuth: true }
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        component: Dashboard,
        meta: { requiresAuth: true }
    },
    // assets
    {
        path: '/assets/bank-accounts',
        name: 'assets.bankaccounts',
        component: () => import('@/pages/assets/BankAccounts.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/assets/bank-accounts/:id',
        name: 'assets.bankaccount.detail',
        component: () => import('@/pages/assets/BankAccountDetail.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/assets/vehicles',
        name: 'assets.vehicles',
        component: () => import('@/pages/assets/Vehicles.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/assets/properties',
        name: 'assets.properties',
        component: () => import('@/pages/assets/Properties.vue'),
        meta: { requiresAuth: true }
    },
    // liabilities
    {
        path: '/liabilities/personal-debts',
        name: 'liabilities.personaldebts',
        component: () => import('@/pages/liabilities/PersonalDebts.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/liabilities/institutional-debts',
        name: 'liabilities.institutionaldebts',
        component: () => import('@/pages/liabilities/InstitutionalDebts.vue'),
        meta: { requiresAuth: true }
    },
    // investments
    {
        path: '/investments/deposits',
        name: 'investments.deposits',
        component: () => import('@/pages/investments/Deposits.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/investments/obligations',
        name: 'investments.obligations',
        component: () => import('@/pages/investments/Obligations.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/investments/gold',
        name: 'investments.gold',
        component: () => import('@/pages/investments/Gold.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/investments/mutual-funds',
        name: 'investments.mutualfunds',
        component: () => import('@/pages/investments/MutualFunds.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/investments/stocks',
        name: 'investments.stocks',
        component: () => import('@/pages/investments/Stocks.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/investments/p2p-lendings',
        name: 'investments.p2plendings',
        component: () => import('@/pages/investments/P2PLendings.vue'),
        meta: { requiresAuth: true }
    },
    // authentication and authorization
    {
        path: '/login',
        name: 'login',
        component: TheLogin,
        meta: { requiresAuth: false }
    },
    // miscellaneous
    {
        path: '/about',
        name: 'misc.about',
        component: () => import('@/pages/About.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/docs',
        name: 'misc.docs',
        component: () => import('@/pages/Docs.vue'),
        meta: { requiresAuth: false }
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