import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {ElementPlusResolver} from "unplugin-vue-components/resolvers";

// https://vitejs.dev/config/
export default defineConfig({
    // 本地测试环境
    server: {
        host: '0.0.0.0',
        port: 3600,
        proxy: {
            "/api": {
                target: `http://127.0.0.1:3601`,
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
    build: {
        terserOptions: {
            compress: {
                drop_console:true,
                drop_debugger: true,
            }
        }
    }

})
