package main

import (
	"fmt"
	"server/config"
	"server/model/system"
	"server/plugin/SystemInit"
	"server/plugin/db"
	"server/router"
	"time"
)

func init() {
	// 执行初始化前等待20s , 让mysql服务完成初始化指令
	time.Sleep(time.Second * 20)
	//初始化redis客户端
	err := db.InitRedisConn()
	if err != nil {
		panic(err)
	}
	// 初始化mysql
	err = db.InitMysql()
	if err != nil {
		panic(err)
	}
}

func main() {
	start()
}

func start() {

	// 启动前先执行数据库内容的初始化工作
	DefaultDataInit()
	// 开启路由监听
	r := router.SetupRouter()
	_ = r.Run(fmt.Sprintf(":%s", config.ListenerPort))
}

func DefaultDataInit() {
	// 如果系统中不存在用户表则进行初始化
	if !system.ExistUserTable() {
		// 初始化数据库相关数据
		SystemInit.TableInIt()
		// 初始化网站基本配置信息
		SystemInit.BasicConfigInit()
	}
	// 初始化影视来源列表信息
	SystemInit.SpiderInit()
}
