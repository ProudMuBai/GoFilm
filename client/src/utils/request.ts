import axios from "axios";
import {ElMessage, ElLoading } from "element-plus";

// 定义加载动画对象
let loading:any
const startLoading = ()=>{
    loading = ElLoading.service({
        lock: true,
        text: `请求发送中...`,
        background: `rgba(255,255,255,0.5)`,
        // target: document.querySelector(`.content`)
    })
}
const closeLoading = ()=>{
    loading.close()
}

const http = (options: any) => {
    return new Promise((resolve, reject) => {
        // create an axios instance
        const service = axios.create({
            // baseURL: import.meta.env.VITE_URL_BASE, // api 的 base_url 注意 vue3
            baseURL: `/api`, // api 的 base_url 注意 vue3
            //   baseURL: 'https://www.baidu.com/api',  // 固定api
            timeout: 80000, // request timeout
        });

        // request interceptor
        service.interceptors.request.use(
            (config: any) => {
                // 开启loading动画
                startLoading()
                // let token: string = ""; //此处换成自己获取回来的token，通常存在在cookie或者store里面
                // if (token) {
                //     // 让每个请求携带token-- ['X-Token']为自定义key 请根据实际情况自行修改
                //     config.headers["X-Token"] = token;
                //
                //     config.headers.Authorization = +token;
                // }
                return config;
            },
            (error) => {
                // Do something with request error
                Promise.reject(error);
            }
        );

        // response interceptor
        service.interceptors.response.use(
            (response) => {
                closeLoading()

                return response.data;
            },
            (error) => {
                closeLoading()
                if (error.response.status == 403) {
                    ElMessage.error("请求异常: ", error)
                } else {
                    ElMessage.error("服务器繁忙，请稍后再试");
                }
                return Promise.reject(error);
            }
        );
        // 请求处理
        service(options)
            .then((res) => {
                resolve(res);
            })
            .catch((error) => {
                reject(error);
            });
    });
};

const ApiGet = (url:string, params?:any)=>{
   return http({url, method:`get`, params,})
}
const ApiPost = (url:string, data:any) =>{
    return http({url, method:`post`, data,})
}

export {http, ApiPost, ApiGet};