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
		{Id: util.GenerateSalt(), Name: "HD(lz)", Uri: `https://cj.lziapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		{Id: util.GenerateSalt(), Name: "HD(sn)", Uri: `https://suoniapi.com/api.php/provide/vod/from/snm3u8/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true, Interval: 2000},
		{Id: util.GenerateSalt(), Name: "HD(bf)", Uri: `https://bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true, Interval: 2500},
		{Id: util.GenerateSalt(), Name: "HD(ff)", Uri: `http://cj.ffzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		{Id: util.GenerateSalt(), Name: "HD(kk)", Uri: `https://kuaikan-api.com/api.php/provide/vod/from/kuaikan/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		//{Id: util.GenerateSalt(), Name: "HD(lzBk)", Uri: `https://cj.lzcaiji.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		//{Id: util.GenerateSalt(), Name: "HD(fs)", Uri: `https://www.feisuzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		//{Id: util.GenerateSalt(), Name: "HD(bfApp)", Uri: `http://app.bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: true},
		//Id: util.GenerateSalt(), {Name: "HD(bfBk)", Uri: `http://by.bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false,CollectType:system.CollectVideo, State: false},
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
			switch task.Model {
			case 0:
				cid, err := spider.AddAutoUpdateCron(task.Id, task.Spec)
				// 如果任务添加失败则直接返回错误信息
				if err != nil {
					log.Println("影视自动更新任务添加失败: ", err.Error())
					continue
				}
				// 将新的定时任务Id记录到Task中
				task.Cid = cid
			case 1:
				cid, err := spider.AddFilmUpdateCron(task.Id, task.Spec)
				// 如果任务添加失败则直接返回错误信息
				if err != nil {
					log.Println("影视更新定时任务添加失败: ", err.Error())
					continue
				}
				// 将定时任务Id记录到Task中
				task.Cid = cid
			}
			system.UpdateFilmTask(task)
		}
	} else {
		// 如果系统中不存在任何定时任务信息, 则添加默认的定时任务
		// 1. 添加一条默认任务, 定时更新所有已启用站点的影片信息
		// 生成任务信息
		task := system.FilmCollectTask{Id: util.GenerateSalt(), Time: config.DefaultUpdateTime, Spec: config.DefaultUpdateSpec,
			Model: 0, State: false, Remark: "每20分钟执行一次已启用站点数据的自动更新"}
		// 添加一条定时任务
		cid, err := spider.AddAutoUpdateCron(task.Id, task.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			log.Println("影视更新定时任务添加失败: ", err.Error())
			return
		}
		// 将定时任务Id记录到Task中
		task.Cid = cid
		// 如果没有异常则将当前定时任务信息记录到redis中
		system.SaveFilmTask(task)
	}

	// 完成初始化后启动 Cron
	spider.CronCollect.Start()
}
