import { createRouter, createWebHistory } from 'vue-router'

import Dashboard from '../pages/Dashboard.vue'

import BankAccounts from '../pages/assets/BankAccounts.vue'
import Vehicles from '../pages/assets/Vehicles.vue'
import Properties from '../pages/assets/Properties.vue'

import PersonalDebts from '../pages/liabilities/PersonalDebts.vue'
import InstitutionalDebts from '../pages/liabilities/InstitutionalDebts.vue'

import Deposits from '../pages/investments/Deposits.vue'
import Obligations from '../pages/investments/Obligations.vue'
import Gold from '../pages/investments/Gold.vue'
import MutualFunds from '../pages/investments/MutualFunds.vue'
import Stocks from '../pages/investments/Stocks.vue'
import P2PLendings from '../pages/investments/P2PLendings.vue'

import Login from '../pages/Login.vue'

import About from '../pages/About.vue'
import Docs from '../pages/Docs.vue'

const routes = [
    // dashboard
    { path: '/dashboard', name: 'Dashboard', component: Dashboard },
    // assets
    { path: '/assets/bank-accounts', name: 'Bank Accounts', component: BankAccounts },
    { path: '/assets/vehicles', name: 'Vehicles', component: Vehicles },
    { path: '/assets/properties', name: 'Properties', component: Properties },
    // liabilities
    { path: '/liabilities/personal-debts', name: 'Personal Debts', component: PersonalDebts },
    { path: '/liabilities/institutional-debts', name: 'Institutional Debts', component: InstitutionalDebts },
    // investments
    { path: '/investments/deposits', name: 'Deposits', component: Deposits },
    { path: '/investments/obligations', name: 'Obligations', component: Obligations },
    { path: '/investments/gold', name: 'Gold', component: Gold },
    { path: '/investments/mutual-funds', name: 'Mutual Funds', component: MutualFunds },
    { path: '/investments/stocks', name: 'Stocks', component: Stocks },
    { path: '/investments/p2p-lendings', name: 'P2P Lendings', component: P2PLendings },
    // authentication and authorization
    { path: '/login', name: 'Login', component: Login },
    // miscellaneous
    { path: '/about', name: 'About', component: About },
    { path: '/docs', name: 'Documentation', component: Docs },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router