import './assets/main.css'

import { createApp } from 'vue'
import { createPinia, type Pinia } from 'pinia'
import App from './App.vue'
import router from './router'

const pinia: Pinia = createPinia()
const app: ReturnType<typeof createApp> = createApp(App)

app.use(router).use(pinia).mount('#app')
