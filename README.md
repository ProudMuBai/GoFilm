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



## 起源

从正式接触编程语言到第一次动手敲代码, , 当时有动手做一些东西的想法,也正是在那时喜欢追番迷二次元, 曾想过做一个自己的动漫站,

但因为知识面匮乏, 总是在进行到某一步时就会遇到一些盲区, 从最开始的静态页面到后面的伪数据, 也实现过一些当时能做到的部分, 

后面慢慢学习的过程中也渐渐遗忘了这个想法, 但因为一些偶然的因素, 想要做一个自己的开源项目, 于是就从零开始慢慢实现并完善了这个

影视站的各个部分, 期间也一点点修改颠覆了一些最开始的思路, 但目前主体功能基本完善, 后续也会定期进行一些bug修复和新功能的更新

如有发现Bug, 或者有好的建议, 可以进行反馈, 欢迎各位大佬来指点一二





## JetBrains 开源证书

感谢Jetbrains提供的免费开源许可, 项目开发中使用GoLang和WebStam让编程变得更加的便捷高效.



<a href="https://www.jetbrains.com/?from=GoFilm" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo."></a>



