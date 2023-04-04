# GoFilm

一个 vite vue 和 gin 实现的在线观影网站

演示站点: [点击访问](http://119.23.231.91:3600)

## 简介

**GoFilm** 

项目采用vite + vue作为前端技术栈, 使用 ElementPlus 作为UI 框架进行开发

后端程序使用 Gin + gorm + go-redis 等相关框架提供接口服务, 使用 gocolly 和 robfig/cron 进行公共影视资源采集和定时更新功能

## 目录结构

- client 客户端项目目录
- server 服务端接口项目目录
- film 项目部署相关配置目录
- 详细说明请查看具体目录中的README文件

```text
GoFilm                          
├─ client                             
│  ├─ public                     
│  │  └─ favicon.ico             
│  ├─ src                        
│  │  ├─ assets                  
│  │  │  ├─ image                
│  │  │  │  ├─ 404.png           
│  │  │  │  ├─ cartoon.png       
│  │  │  │  ├─ film.png          
│  │  │  │  ├─ play.png          
│  │  │  │  └─ tv.png            
│  │  │  └─ svg                  
│  │  │     ├─ cartoon.svg       
│  │  │     ├─ film.svg          
│  │  │     └─ tv.svg            
│  │  ├─ components              
│  │  │  ├─ Footer.vue           
│  │  │  ├─ Header.vue           
│  │  │  ├─ RelateList.vue       
│  │  │  └─ Util.vue             
│  │  ├─ router                  
│  │  │  └─ router.ts            
│  │  ├─ utils                   
│  │  │  └─ request.ts           
│  │  ├─ views                   
│  │  │  ├─ error                
│  │  │  │  └─ Error404.vue      
│  │  │  ├─ index                
│  │  │  │  ├─ CategoryFilm.vue  
│  │  │  │  ├─ FilmDetails.vue   
│  │  │  │  ├─ Home.vue          
│  │  │  │  ├─ Play.vue          
│  │  │  │  └─ SearchFilm.vue    
│  │  │  └─ IndexHome.vue        
│  │  ├─ App.vue                 
│  │  ├─ main.ts                 
│  │  ├─ style.css               
│  │  └─ vite-env.d.ts           
│  ├─ auto-imports.d.ts          
│  ├─ components.d.ts            
│  ├─ index.html                 
│  ├─ package-lock.json          
│  ├─ package.json               
│  ├─ README.md                  
│  ├─ tsconfig.json              
│  ├─ tsconfig.node.json         
│  └─ vite.config.ts             
├─ film                          
│  ├─ data                       
│  │  ├─ nginx                   
│  │  │  ├─ html                 
│  │  │  └─ nginx.conf           
│  │  └─ redis                   
│  │     └─ redis.conf           
│  ├─ docker-compose.yml         
│  ├─ Dockerfile                 
│  └─ README.md                  
├─ server                        
│  ├─ config                     
│  │  └─ DataConfig.go           
│  ├─ controller                 
│  │  ├─ IndexController.go      
│  │  └─ SpiderController.go     
│  ├─ logic                      
│  │  └─ IndexLogic.go           
│  ├─ model                      
│  │  ├─ Categories.go           
│  │  ├─ LZJson.go               
│  │  ├─ Movies.go               
│  │  └─ Search.go               
│  ├─ plugin                     
│  │  ├─ common                  
│  │  │  ├─ ProcessCategory.go   
│  │  │  └─ ProcessMovies.go     
│  │  ├─ db                      
│  │  │  ├─ mysql.go             
│  │  │  └─ redis.go             
│  │  └─ spider                  
│  │     ├─ Spider.go            
│  │     ├─ SpiderCron.go        
│  │     └─ SpiderRequest.go     
│  ├─ router                     
│  │  └─ router.go               
│  ├─ go.mod                     
│  ├─ go.sum                     
│  ├─ main.go                    
│  └─ README.md                  
├─ LICENSE                       
└─ README.md                     
```



## 项目说明

本项目出于个人对于学习技术栈的一次实践, 也算是完成最初学习时想做却做不了的一个念想

后续可能会陆续进行一些修改和完善, 项目规范性可能没有那么强, 可能也存在一些问题

如有发现Bug, 或者有好的建议, 可以进行反馈, 欢迎各位大佬来指定一二 
