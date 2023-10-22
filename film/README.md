# Film Build

## 1. 目录结构说明

- data 服务容器相关的数据和配置信息存放目录
  - nginx 
    - html目录用于上传vite项目构建后的dist下的页面相关文件
    - nginx.conf 配置后端接口的代理和端口等相关信息
  - redis
    - redis.conf 配置redis的远程访问和密码等信息
- server 服务端资源文件, Dockerfile 生成镜像时所需
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
├─ server                          
├─ docker-compose.yml              
├─ Dockerfile                      
└─ README.md                       
```

>此目录下的client内容和server内容并不一定与client同步 (小更新可能不会实时同步到运行服务器上)
>
>可自行将根目录下的server和client内容与此目录下的对应文件进行替换

## 2. 程序构建运行

### 1. 环境准备

1.  Linux 服务器
2.  安装 docker, docker compose 服务
    - Centos 安装 Docker Engine  [官方文档链接](https://docs.docker.com/engine/install/centos/)
    - Ubuntu 安装 Docker Engine   [官方文档链接](https://docs.docker.com/engine/install/ubuntu/)

```shell
# Centos 系统安装Docker Engine示例
# 1. 卸载旧版本Docker
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine

#2. 设置存储库
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
#3. 安装最新版本Docker
$ sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
#启动 Docker 服务
sudo systemctl start docker
```

### 2. 启动流程

> 如果使用默认配置信息,则执行如下流程

- 将本项目中的 film 文件夹完整的上传到服务器的 ` /opt/` 目录下 (放在其他目录下时需同步修改 `Dockerfile` 以及 `docker-compose.yml` 文件中的相关路径)
- 进入服务器中的 `/opt/film/` 目录并执行 `docker compose build` 构建相关docker镜像
- 在 `/opt/film/` 目录下执行命令 `docker compose up -d` (后台运行服务)
- 使用 `docker ps` 命令查看相关服务是否成功启动
- 等待后端程序初始化工作和数据爬取, 大概3~8分钟左右
- 停止服务 `docker compose down`
- 查看服务容器运行状态 `docker ps`
- 浏览器中访问nginx服务地址查看效果, 例: [http://xxx.xxx.xxx:3600](http://xxx.xxx.xxx:3600)

### 3.服务配置信息修改

- film 后端接口服务配置, `film/server` 下存放了程序的构建文件, 修改后重新构建镜像即可
- mysql 用户名密码和端口信息直接修改 `docker-compose.yml`  文件中的相关配置即可
- redis 服务信息配置需修改 `/film/data/redis/redis.conf` 文件
- nginx 配置文件 `/film/data/nginx/nginx.conf` 

>注意事项

-  mysql 和 redis 服务配置修改后需要同步修改 `/film/server/config/DataConfig.go` 中的连接地址和账户名信息

```go
## 配置使用的用户名密码信息需和ocker-compose.yml文件中设置的一致
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
## docker-compose.yml (设置服务的启动端口和服务名以及账户密码信息)
mysql:
    container_name: film_mysql
    image: mysql
    ports:
    - 3610:3306
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: FilmSite
 redis:
    container_name: film_redis
    image: redis
    ports:
      - 3620:6379
```