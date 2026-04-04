import {defineConfig} from 'vite'
import vue from "@vitejs/plugin-vue"
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {ElementPlusResolver} from "unplugin-vue-components/resolvers";
// @ts-ignore
import path from 'path'
import { fileURLToPath, URL } from 'node:url' // 1. 导入 node:url 模块


const __dirname = path.dirname(fileURLToPath(import.meta.url))


export default defineConfig({
    // 本地测试环境
    server: {
        host: '0.0.0.0',
        port: 3600,
        proxy: {
            "/api": {
                target: `http://127.0.0.1:3601`,
                // target: `http://www.mubai.us.ci:3601`,
                changeOrigin: true, // 允许跨域
                rewrite: path => path.replace(/^\/api/, '')
            }
        },
    },

    // nginx发布构建时使用此配置
    // server: {
    //     host: 'localhost',
    //     port: 3600,
    //     proxy: {
    //         "/api": {
    //             target: `http://localhost`,
    //             changeOrigin: true, // 允许跨域
    //             rewrite: path => path.replace(/^\/api/,'')
    //         }
    //     },
    // },

    plugins: [
        vue(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
        }),
        Components({
            resolvers: [ElementPlusResolver()],
        }),
    ],
    css: {devSourcemap: true},
    esbuild: {
        // 生产环境移除 console 和 debugger
        drop: ['console', 'debugger'],
    },
    resolve: {
        alias: {
            // 2. 配置 @ 指向 src 目录的绝对路径
            '@': path.resolve(__dirname, './src'),
        },
    },

})
