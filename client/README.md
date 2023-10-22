# Film Client

## 简介 

- Client 是项目的前端组成部分, 由 Vite 初始化构建, 使用 vue3 + ts 模板
- 项目中使用 ElementPlus 作为UI组件, 使用Axios作为HTTP请求工具 
- 项目打包构建后使用nginx发布运行



## 项目结构

- components 目录用于存放vue的公共组件
- views 目录用于存放vue-router组件
- request.ts 对Axios进行一些简单的封装处理
- router.ts 配置路由路径和组件之间的映射关系

```text
client                              
├─ public                           
│  └─ favicon.ico                   
├─ src                              
│  ├─ assets                        
│  │  ├─ css                        
│  │  │  ├─ classify.css            
│  │  │  ├─ film.css                
│  │  │  └─ pagination.css          
│  │  └─ image                      
│  │     ├─ 404.png                 
│  │     └─ play.png                
│  ├─ components                    
│  │  ├─ Loading                    
│  │  │  ├─ index.ts                
│  │  │  └─ Loading.vue             
│  │  ├─ FilmList.vue               
│  │  ├─ Footer.vue                 
│  │  ├─ Header.vue                 
│  │  ├─ RelateList.vue             
│  │  └─ Util.vue                   
│  ├─ router                        
│  │  └─ router.ts                  
│  ├─ utils                         
│  │  ├─ cookie.ts                  
│  │  └─ request.ts                 
│  ├─ views                         
│  │  ├─ error                      
│  │  │  └─ Error404.vue            
│  │  ├─ index                      
│  │  │  ├─ FilmClassify.vue        
│  │  │  ├─ FilmClassifySearch.vue  
│  │  │  ├─ FilmDetails.vue         
│  │  │  ├─ Home.vue                
│  │  │  ├─ Play.vue                
│  │  │  └─ SearchFilm.vue          
│  │  └─ IndexHome.vue              
│  ├─ App.vue                       
│  ├─ main.ts                       
│  ├─ style.css                     
│  └─ vite-env.d.ts                 
├─ auto-imports.d.ts                
├─ components.d.ts                  
├─ index.html                       
├─ package.json                     
├─ README.md                        
├─ tsconfig.json                    
├─ tsconfig.node.json               
└─ vite.config.ts                   
```

## 启动方式

### 本地运行方式

1. 进入client文件夹下
2. 执行 `npm install` 安装项目相关依赖 (请确保已经安装NodeJS)
3. 修改 vite.config.ts 文件中的 后端接口地址

```typescript
# 将 target 属性值设置为后端接口的请求地址
server: {
    host: '0.0.0.0',
    port: 3600,
    proxy: {
        "/api": {
            target: `http://127.0.0.1:3601`,
            changeOrigin: true, // 允许跨域
            rewrite: path => path.replace(/^\/api/,'')
        }
    },
},
```

4. 使用 `npm run dev` 启动项目
5. 打开浏览器访问 [http://127.0.0.1:3600](http://127.0.0.1:3600) 



### nginx部署

1. 进入client文件夹下
2. 执行 `npm install` 安装项目相关依赖 (请确保已经安装NodeJS)
3. 修改 vite.config.ts 文件中的 后端接口地址

```typescript
# # 将 target 属性值设置为nginx服务器url地址
server: {
    host: 'localhost',
    port: 3600,
    proxy: {
        "/api": {
            target: `http://localhost`,
            changeOrigin: true, // 允许跨域
            rewrite: path => path.replace(/^\/api/,'')
        }
    },
},
```

4. 执行 `npm run build` 将项目进行打包构建
5. 将打包后生成的 dist目录中的所有文件复制到 nginx的 html文件夹下
6. 使用浏览器访问 nginx 的服务地址, 例:  IP:Port, http://localhost:3600



  









































