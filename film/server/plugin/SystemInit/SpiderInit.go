package SystemInit

import (
	"log"
	"server/config"
	"server/model/system"
	"server/plugin/common/util"
	"server/plugin/spider"
)

// SpiderInit 数据采集相关信息初始化
func SpiderInit() {
	FilmSourceInit()
	CollectCrontabInit()
}

// FilmSourceInit  初始化预存站点信息 提供一些预存采集连Api链接
func FilmSourceInit() {
	// 首先获取filmSourceList 数据, 如果存在则直接返回
	if system.ExistCollectSourceList() {
		return
	}
	var l []system.FilmSource = []system.FilmSource{
		{Id: util.GenerateSalt(), Name: "HD(LZ)", Uri: `https://cj.lziapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(BF)", Uri: `https://bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false, Interval: 2500},
		{Id: util.GenerateSalt(), Name: "HD(FF)", Uri: `http://cj.ffzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(OK)", Uri: `https://api.okzyw.net/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(LY)", Uri: `https://360zy.com/api.php/provide/vod/at/json`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(SN)", Uri: `https://suoniapi.com/api.php/provide/vod/from/snm3u8/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false, Interval: 2000},
		{Id: util.GenerateSalt(), Name: "HD(DB)", Uri: `https://caiji.dbzy.tv/api.php/provide/vod/from/dbm3u8/at/josn/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(TT)", Uri: `https://caiji.dyttzyapi.com/api.php/provide/vod/at/json/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
	}
	err := system.SaveCollectSourceList(l)
	if err != nil {
		log.Println("SaveSourceApiList Error: ", err)
	}
}

// CollectCrontabInit 初始化系统预定义的定时任务
func CollectCrontabInit() {
	// 如果系统已经存在Task定时任务信息,则将redis中的定时任务信息重新添加到执行队列
	if system.ExistTask() {
		// 将系统中的定时任务重新设置到 CollectCron中
		for _, task := range system.GetAllFilmTask() {
			// 添加任务到 cron中
			cid, err := spider.AddCron(task.Id, task.Spec)
			// 如果任务添加失败则直接返回错误信息
			if err != nil {
				log.Println("自动任务恢复失败: ", err.Error())
				continue
			}
			// 将新的定时任务Id记录到Task中
			task.Cid = cid
			system.UpdateFilmTask(task)
		}
	} else {
		/*
			如果系统中不存在任何定时任务信息, 则添加默认的定时任务
			1. 添加一条默认任务, 定时更新所有已启用站点的影片信息
			2. 添加一条默认任务, 定时处理采集失败的记录
			3.生成任务信息
		*/
		collectTask := system.FilmCollectTask{Id: util.GenerateSalt(), Time: config.DefaultUpdateTime, Spec: config.DefaultUpdateSpec,
			Model: 0, State: false, Remark: "每20分钟执行一次已启用站点数据的自动更新"}
		// 添加一条定时任务-影片定时更新
		cid, err := spider.AddCron(collectTask.Id, collectTask.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			log.Println("影视更新定时任务添加失败: ", err.Error())
			return
		}
		// 将定时任务Id记录到Task中
		collectTask.Cid = cid
		// 如果没有异常则将当前定时任务信息记录到redis中
		system.SaveFilmTask(collectTask)

		// 添加一条定时任务-定期处理失败请求
		recoverTask := system.FilmCollectTask{Id: util.GenerateSalt(), Time: 0, Spec: config.EveryWeekSpec,
			Model: 2, State: false, Remark: "每周日凌晨3点清理一次采集失败的采集记录"}
		// 添加一条定时任务-影片定时更新
		cid, err = spider.AddCron(recoverTask.Id, recoverTask.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			log.Println("失败采集恢复定时任务添加失败: ", err.Error())
			return
		}
		// 将定时任务Id记录到Task中
		recoverTask.Cid = cid
		// 如果没有异常则将当前定时任务信息记录到redis中
		system.SaveFilmTask(recoverTask)

		// 添加定时任务-定期同步redis中的影片信息
		syncTask := system.FilmCollectTask{Id: util.GenerateSalt(), Time: 0, Spec: config.PeriodSyncSpec,
			Model: 3, State: false, Remark: "每周日凌晨4点同步一次影片信息"}
		cid, err = spider.AddCron(syncTask.Id, syncTask.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			log.Println("影片信息同步任务执行失败: ", err.Error())
			return
		}
		// 将定时任务Id记录到Task中
		syncTask.Cid = cid
		// 如果没有异常则将当前定时任务信息记录到redis中
		system.SaveFilmTask(syncTask)

	}

	// 完成初始化后启动 Cron
	spider.CronCollect.Start()
}
