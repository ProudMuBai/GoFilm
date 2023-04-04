package main

import (
	"server/model"
	"server/plugin/db"
	"server/plugin/spider"
	"server/router"
	"time"
)

func init() {
	// 执行初始化前等待30s , 让mysql服务完成初始化指令
	time.Sleep(time.Second * 30)
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
	// 开启前先判断是否需要执行Spider
	ExecSpider()
	// web服务启动后开启定时任务, 用于定期更新资源
	spider.RegularUpdateMovie()
	// 开启路由监听
	r := router.SetupRouter()
	_ = r.Run(`:3601`)

}

func ExecSpider() {
	// 判断分类信息是否存在
	isStart := model.ExistsCategoryTree()
	// 如果分类信息不存在则进行一次完整爬取
	if !isStart {
		spider.StartSpider()
	}
}
