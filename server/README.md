# Film Server

## 简介

- server 是本项目的后端项目
- 主要用于提供前端项目需要的 API数据接口, 以及数据搜集和更新
- 实现思路 : 
  - 使用 gocolly 获取公开的影视资源, 
  - 将请求数据通过程序处理整合成统一格式后使用redis进行暂存
  - 使用 mysql 存储收录的影片的检索信息, 用于影片检索, 分类
  - 使用 gin 作为web服务, 提供相应api接口
- 项目依赖

```go
# gin web服务框架, 用于处理与前端工程的交互
github.com/gin-gonic/gin v1.9.0
# gocolly go语言爬虫框架, 用于搜集公共影视资源
github.com/gocolly/colly/v2 v2.1.0
# go-redis redis交互程序
github.com/redis/go-redis/v9 v9.0.2
# gorm 用于处理与mysql数据库的交互
gorm.io/gorm v1.24.6
gorm.io/driver/mysql v1.4.7
```



## 项目结构

> 项目主要目录结构

- config				用于存放项目中使用的配置信息和静态常量
- controller	      请求处理控制器
- logic                请求处理逻辑实现
- model              数据模型结构体以及与数据库交互
- plugin              项目所需的插件工具集合
  - common     公共依赖
  - db              数据库配置信息
  - spider         gocolly配置, 执行逻辑, 数据前置处理等

```text
server                       
├─ config                    
│  └─ DataConfig.go          
├─ controller                
│  ├─ IndexController.go     
│  └─ SpiderController.go    
├─ logic                     
│  └─ IndexLogic.go          
├─ model                     
│  ├─ Categories.go          
│  ├─ LZJson.go              
│  ├─ Movies.go              
│  └─ Search.go              
├─ plugin                    
│  ├─ common                 
│  │  ├─ JsonUtils.go        
│  │  ├─ ProcessCategory.go  
│  │  └─ ProcessMovies.go    
│  ├─ db                     
│  │  ├─ mysql.go            
│  │  └─ redis.go            
│  └─ spider                 
│     ├─ Spider.go           
│     ├─ SpiderCron.go       
│     └─ SpiderRequest.go    
├─ router                    
│  └─ router.go              
├─ go.mod                    
├─ go.sum                    
├─ main.go                   
└─ README.md                 
```



## 启动方式

### 本地运行

1.  修改 /server/plugin/db 目录下的 mysql.go 和 redis.go 中的连接地址和用户名密码
2. 在 server 目录下执行 `go run main.go`

