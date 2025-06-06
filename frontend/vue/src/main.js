import { createApp } from 'vue'
import router from '@/configs/routes'
import { createPinia } from 'pinia'
import App from '@/App.vue'
import { useAuthStore } from '@/stores/authStore'

import '@/style.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
    faEdit,
    faEye,
    faEyeSlash,
    faPlus,
    faTrash
} from '@fortawesome/free-solid-svg-icons'

library.add(
    faEdit,
    faEye,
    faEyeSlash,
    faPlus,
    faTrash)

const pinia = createPinia()

createApp(App)
    .use(pinia)
    .use(router)
    .component('font-awesome-icon', FontAwesomeIcon)
    .mount('#app')

const authStore = useAuthStore()
authStore.hydrate()
