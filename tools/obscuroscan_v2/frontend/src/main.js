import './assets/main.css'

import { createApp } from 'vue'
import {createPinia, setActivePinia} from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import Config from "@/lib/config";

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
setActivePinia(pinia)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(ElementPlus)
app.config.globalProperties.RunConfig =new Config()
app.use(router)


app.mount('#app')
