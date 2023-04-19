package spider

import (
	"github.com/robfig/cron/v3"
	"log"
	"server/config"
	"server/model"
)

// RegularUpdateMovie 定时更新, 每半小时获取一次站点的最近x小时数据
func RegularUpdateMovie() {
	c := cron.New(cron.WithSeconds())
	// 开启定时任务每x 分钟更新一次最近x小时的影片数据
	_, err := c.AddFunc(config.CornMovieUpdate, func() {
		// 执行更新最近x小时影片的Spider
		log.Println("执行一次影片更新任务...")
		UpdateMovieDetail()
		// 执行更新任务后清理redis中的相关API接口数据缓存
		clearCache()
	})

	// 开启定时任务每月最后一天凌晨两点, 执行一次清库重取数据
	_, err = c.AddFunc(config.CornUpdateAll, func() {
		StartSpiderRe()
	})

	if err != nil {
		log.Println("Corn Start Error: ", err)
	}

	c.Start()
}

// 清理API接口数据缓存
func clearCache() {
	model.RemoveCache(config.IndexCacheKey)
}
