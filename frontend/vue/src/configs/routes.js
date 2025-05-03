import { createRouter, createWebHistory } from 'vue-router'

import About from '../pages/About.vue'
import Docs from '../pages/Docs.vue'
import Login from '../pages/Login.vue'

const routes = [
    { path: '/about', name: 'About', component: About },
    { path: '/docs', name: 'Documentation', component: Docs },
    { path: '/login', name: 'Login', component: Login }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router