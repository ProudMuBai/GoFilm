package spider

import (
	"errors"
	"fmt"
	"log"
	"server/config"
	"server/model/system"

	"github.com/robfig/cron/v3"
)

var (
	CronCollect *cron.Cron = CreateCron()
)

// CreateCron 创建定时任务
func CreateCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

// AddCron 添加定时任务
func AddCron(id, spec string) (cron.EntryID, error) {
	// 校验 spec 表达式的有效性
	if err := ValidSpec(spec); err != nil {
		return -99, errors.New(fmt.Sprint("定时任务添加失败,Cron表达式校验失败: ", err.Error()))
	}
	return CronCollect.AddFunc(spec, func() {
		// 通过 Id 获取任务相关数据
		ft, err := system.GetFilmTaskById(id)
		if err != nil {
			log.Println("FilmCollectCron Exec Failed: ", err)
		}
		// 开启对系统中已启用站点的自动更新
		if ft.State {
			switch ft.Model {
			case 0:
				AutoCollect(ft.Time)
				log.Println("执行一次已启用站点的影片自动更新任务")
			case 1:
				// 对指定ids的资源站数据进行更新操作
				BatchCollect(ft.Time, ft.Ids...)
				log.Println("执行一次指定站点的影片自动更新任务")
			case 2:
				FullRecoverSpider()
				log.Println("执行一次采集失败的记录处理任务")
			case 3:
				FullSyncMovieDetail()
				log.Println("执行一次所有站点的影片信息同步任务")
			}
		}
	})

}

// RemoveCron 删除定时任务
func RemoveCron(id cron.EntryID) {
	// 通过定时任务EntryID移出对应的定时任务
	CronCollect.Remove(id)
}

// GetEntryById 返回定时任务的相关时间信息
func GetEntryById(id cron.EntryID) cron.Entry {
	//log.Printf("CronInfo: %+v\n", CronCollect.Entries())
	//log.Println("Corn Next Execute Time:", CronCollect.Entry(id).Next.Format(time.DateTime))
	return CronCollect.Entry(id)
}

// ValidSpec 校验cron表达式是否有效 不能精确到秒
func ValidSpec(spec string) error {
	// 自定义解释器
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	//if _, err := parser.Parse(spec); err != nil {
	//	return err
	//}
	_, err := parser.Parse(spec)
	return err
}

// ClearCache 清理API接口数据缓存
func ClearCache() {
	system.RemoveCache(config.IndexCacheKey)
}
