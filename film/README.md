# Film Build

## 1. 目录结构说明

- data 服务容器相关的数据和配置信息存放目录
  - nginx 
    - html目录用于上传vite项目构建后的dist下的页面相关文件
    - nginx.conf 配置后端接口的代理和端口等相关信息
  - redis
    - redis.conf 配置redis的远程访问和密码等信息
- docker-compose.yml docker 服务配置启动文件
- Dockerfile go程序镜像构建文件

```text
film                   
├─ data                
│  ├─ nginx            
│  │  ├─ html      
│  │  └─ nginx.conf    
│  └─ redis            
│     └─ redis.conf    
├─ docker-compose.yml  
├─ Dockerfile          
└─ README.md           
```







## 2. 程序构建运行

### 1. 环境准备

1.  Linux 服务器
2. 安装 docker, docker compose 服务

  

### 2. 启动流程

> 如果使用默认配置信息,则执行如下流程

- 将本项目中的 film 文件夹完整的上传到服务器的 ` /opt/` 目录下
- 进入服务器中的 `/opt/film/` 目录并执行 `docker compose build` 构建相关docker镜像
- 在 `/opt/film/` 目录下执行命令 `docker compose up -d` (后台运行服务)
- 使用 `docker ps` 命令查看相关服务是否成功启动
- 等待后端程序初始化工作和数据爬取, 大概8分钟左右
- 浏览器中访问nginx服务地址查看效果, 例: [http://xxx.xxx.xxx:3600](http://xxx.xxx.xxx:3600)



### 3.服务配置信息修改

- film 后端接口服务配置, `film/server` 下存放了程序的构建文件, 修改后重新构建镜像即可
- mysql 用户名密码和端口信息直接修改 `docker-compose.yml`  文件中的相关配置即可
- redis 服务信息配置需修改 `/film/data/redis/redis.conf` 文件
- nginx 配置文件 `/film/data/nginx/nginx.conf` 

>注意事项

-  mysql 和 redis 服务配置修改后需要同步修改 `/film/server/config/DataConfig.go` 中的连接地址和账户名信息

```go
const (
	// mysql服务配置信息修改
	mysqlDsn = "用户名:密码$@(服务名:服务端口)/FilmSite?charset=utf8mb4&parseTime=True&loc=Local"

	/*
		redis 配置信息
		RedisAddr host:port
		RedisPassword redis访问密码
		RedisDBNo 使用第几号库
	*/
	RedisAddr = `服务名:服务端口`
	RedisPassword = `密码`
	RedisDBNo = 0
)
```