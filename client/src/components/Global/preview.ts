import ImageViewer from "./ImageViewer.vue";
import {createApp} from "vue";

const Preview = (options:any) =>{
    // 默认创建 ImageViewer 组件时为显示状态
    options.show = true
    // 创建节点用户挂载
    const el = document.createElement("div")
    document.body.appendChild(el)
    const app = createApp(ImageViewer, {
        options,
        remove(){
            app.unmount()
            document.body.removeChild(el)
        }
    })
    return app.mount(el)
}

export {Preview}