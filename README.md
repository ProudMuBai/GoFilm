# GoFilm

一个基于 vue 和 gin 实现的在线观影网站

效果展示: <a href="https://m.mubai.link/" target="_blank">点击访问演示站点</a> 

## 简介

**GoFilm** 

项目采用vite + vue作为前端技术栈, 使用 ElementPlus 作为UI 框架进行开发

后端程序使用 Gin + gorm + go-redis 等相关框架提供接口服务, 使用 gocolly 和 robfig/cron 进行公共影视资源采集和定时更新功能



## 新版本说明

**网站前台**

- 前台部分对网站名称以及播放源等部分信息与后台数据进行关联, 可通过后台进行修改
- 影片详情部分以及首页导航数据结构发生变化, 样式保持一致
- 默认访问地址: `服务器IP:默认端口 [http://127.0.0.1:3600]`

**管理后台**

- 新增管理后台功能组, 主要用于对 `采集站点`, `定时更新`, `网站基本信息`, `图片同步`, `影片分类`, `影片信息` 等进行管理 (部分功能正在完善中, 不影响已有功能使用)
- 管理后台访问需进行登录, 默认账号/密码: `admin admin` (登录成功后自行通过右上下拉弹窗进行密码修改)
- 具体情况请自行搭建访问
- 默认访问地址: `服务器IP:默认端口/manage [http://127.0.0.1:3600/manage]`

**更新说明**

- 后台功能完善阶段时不会同步更新到演示站点, 需自行使用服务器搭建体验
- 使用中出现问题可在项目 Issues 中进行描述, 有需要添加的新功能和好的建议也可以提供
- 新版本安装方法以及使用说明请查看本项目 film 文件夹下的说明文件

## 目录结构

- client 客户端项目目录 [Client简介](./client/README.md)
- server 服务端接口项目目录 [Server简介](./client/README.md)
- film 项目部署相关配置目录 [film 项目安装](./film/README.md)
- 详细说明请查看具体目录中的README文件

```text
GoFilm-main                            
├─ client                              
│  ├─ public                           
│  │  └─ favicon.ico                   
│  ├─ src                              
│  │  ├─ assets                        
│  │  │  ├─ css                        
│  │  │  │  ├─ classify.css            
│  │  │  │  ├─ film.css                
│  │  │  │  └─ pagination.css          
│  │  │  └─ image                      
│  │  │     ├─ 404.png                 
│  │  │     └─ play.png                
│  │  ├─ components                    
│  │  │  ├─ Loading                    
│  │  │  │  ├─ index.ts                
│  │  │  │  └─ Loading.vue             
│  │  │  ├─ FilmList.vue               
│  │  │  ├─ Footer.vue                 
│  │  │  ├─ Header.vue                 
│  │  │  ├─ RelateList.vue             
│  │  │  └─ Util.vue                   
│  │  ├─ router                        
│  │  │  └─ router.ts                  
│  │  ├─ utils                         
│  │  │  ├─ cookie.ts                  
│  │  │  └─ request.ts                 
│  │  ├─ views                         
│  │  │  ├─ error                      
│  │  │  │  └─ Error404.vue            
│  │  │  ├─ index                      
│  │  │  │  ├─ FilmClassify.vue        
│  │  │  │  ├─ FilmClassifySearch.vue  
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
│  ├─ package.json                     
│  ├─ README.md                        
│  ├─ tsconfig.json                    
│  ├─ tsconfig.node.json               
│  └─ vite.config.ts                   
├─ film                                
│  ├─ data                             
│  │  ├─ nginx                         
│  │  │  ├─ html                       
│  │  │  │  ├─ assets                  
│  │  │  │  │  ├─ 404-b813c94a.png     
│  │  │  │  │  ├─ index-984712d6.js    
│  │  │  │  │  ├─ index-de4c7ff5.css   
│  │  │  │  │  └─ play-bb9c8990.png    
│  │  │  │  ├─ favicon.ico             
│  │  │  │  └─ index.html              
│  │  │  └─ nginx.conf                 
│  │  └─ redis                         
│  │     └─ redis.conf                 
│  ├─ server                           
│  │  ├─ config                        
│  │  │  └─ DataConfig.go              
│  │  ├─ controller                    
│  │  │  ├─ IndexController.go         
│  │  │  └─ SpiderController.go        
│  │  ├─ logic                         
│  │  │  ├─ IndexLogic.go              
│  │  │  └─ SpiderLogic.go             
│  │  ├─ model                         
│  │  │  ├─ Categories.go              
│  │  │  ├─ Movies.go                  
│  │  │  ├─ RequestParams.go           
│  │  │  ├─ ResponseJson.go            
│  │  │  └─ Search.go                  
│  │  ├─ plugin                        
│  │  │  ├─ common                     
│  │  │  │  ├─ dp                      
│  │  │  │  │  ├─ ProcessCategory.go   
│  │  │  │  │  └─ ProcessMovies.go     
│  │  │  │  └─ param                   
│  │  │  │     └─ SimpleParam.go       
│  │  │  ├─ db                         
│  │  │  │  ├─ mysql.go                
│  │  │  │  └─ redis.go                
│  │  │  └─ spider                     
│  │  │     ├─ Spider.go               
│  │  │     ├─ SpiderCron.go           
│  │  │     └─ SpiderRequest.go        
│  │  ├─ router                        
│  │  │  └─ router.go                  
│  │  ├─ go.mod                        
│  │  ├─ go.sum                        
│  │  ├─ main.go                       
│  │  └─ README.md                     
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
│  │  ├─ IndexLogic.go                 
│  │  └─ SpiderLogic.go                
│  ├─ model                            
│  │  ├─ Categories.go                 
│  │  ├─ Movies.go                     
│  │  ├─ RequestParams.go              
│  │  ├─ ResponseJson.go               
│  │  └─ Search.go                     
│  ├─ plugin                           
│  │  ├─ common                        
│  │  │  ├─ dp                         
│  │  │  │  ├─ ProcessCategory.go      
│  │  │  │  └─ ProcessMovies.go        
│  │  │  ├─ param                      
│  │  │  │  └─ SimpleParam.go          
│  │  │  └─ util                       
│  │  │     ├─ FileDownload.go         
│  │  │     └─ Request.go              
│  │  ├─ db                            
│  │  │  ├─ mysql.go                   
│  │  │  └─ redis.go                   
│  │  └─ spider                        
│  │     ├─ Spider.go                  
│  │     └─ SpiderCron.go              
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



## 更新迭代计划

- 目前用户界面的一些功能有待开发和完善, 大家也可以继续提供一些好的建议
- 目前pc端的历史记录写了一个简单的测试版, 后面有时间会同步完善pc和wrap端的历史记录和收藏功能
- 前台功能目前基本满足观看的需求, 后续考虑切入一些登录和账户以及管理后台的功能,慢慢完善这个项目.



## JetBrains 开源证书

感谢Jetbrains提供的免费开源许可, GoLang 和 WebStorm 为编程开发带来了良好的体验.



<a href="https://www.jetbrains.com/?from=GoFilm" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo."></a>



