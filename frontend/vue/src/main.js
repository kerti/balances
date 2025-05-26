import { createApp } from 'vue'
import '@/style.css'
import router from '@/configs/routes'
import App from '@/App.vue'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faEdit, faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons'

library.add(faEdit, faEye, faEyeSlash)

createApp(App)
    .use(router)
    .component('font-awesome-icon', FontAwesomeIcon)
    .mount('#app')
