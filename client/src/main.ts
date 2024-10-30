import {createApp} from 'vue'
import './style.css'
import App from './App.vue'
import { router} from "./router/router";
// 引入elementPlus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import {createHead} from '@vueuse/head'

// 使用自定义loading

const app = createApp(App)


app.use(ElementPlus)
// 引入路由
app.use(router)

app.use(createHead())

app.mount('#app')

export default app


