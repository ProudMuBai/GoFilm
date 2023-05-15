import { createApp, reactive } from 'vue'

import customLoading from './Loading.vue'

const msg = reactive({
    show: false,
    title: '拼命加载中...'
})

const $loading = createApp(customLoading, {msg}).mount(document.createElement('div'))
const load = {
    start(title:string) { // 控制显示loading的方法
        msg.show = true
        msg.title = title
        document.body.appendChild($loading.$el)
        document.body.style.overflow = 'hidden'
    },
    close() { // 控制loading隐藏的方法
        msg.show = false
        document.body.style.overflow = 'auto'
    }
}
export  { load }


