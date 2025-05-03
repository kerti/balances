import { createApp } from 'vue'
import './style.css'
import router from './configs/routes'
import App from './App.vue'

createApp(App)
    .use(router)
    .mount('#app')
