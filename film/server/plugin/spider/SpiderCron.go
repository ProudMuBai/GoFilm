package spider

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"server/config"
	"server/model/system"
	"time"
)

var (
	CronCollect *cron.Cron = CreateCron()
)

// RegularUpdateMovie 定时更新, 每半小时获取一次站点的最近x小时数据
func RegularUpdateMovie() {
	//创建一个定时任务对象
	c := cron.New(cron.WithSeconds())
	// 添加定时任务每x 分钟更新一次最近x小时的影片数据
	taskId, err := c.AddFunc(config.CornMovieUpdate, func() {
		// 执行更新最近x小时影片的Spider
		log.Println("执行一次影片更新任务...")
		UpdateMovieDetail()
		// 执行更新任务后清理redis中的相关API接口数据缓存
		clearCache()
	})

	// 开启定时任务每月最后一天凌晨两点, 执行一次清库重取数据
	taskId2, err := c.AddFunc(config.CornUpdateAll, func() {
		StartSpiderRe()
	})

	if err != nil {
		log.Println("Corn Start Error: ", err)
	}

	log.Println(taskId, "------", taskId2)
	log.Printf("%v", c.Entries())

	//c.Start()
}

// StartCrontab 启动定时任务
func StartCrontab() {
	// 从redis中读取待启动的定时任务列表

	// 影片更新定时任务列表
	CronCollect.Start()
}

func CreateCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

// AddFilmUpdateCron 添加影片更新定时任务
func AddFilmUpdateCron(id, spec string) (cron.EntryID, error) {
	// 校验 spec 表达式的有效性
	if err := ValidSpec(spec); err != nil {
		return -99, errors.New(fmt.Sprint("定时任务添加失败,Cron表达式校验失败: ", err.Error()))
	}
	return CronCollect.AddFunc(spec, func() {
		// 通过创建任务时生成的 Id 获取任务相关数据
		ft, err := system.GetFilmTaskById(id)
		if err != nil {
			log.Println("FilmCollectCron Exec Failed: ", err)
		}
		// 如果当前定时任务状态为开启则执行对应的采集任务
		if ft.State && ft.Model == 1 {
			// 对指定ids的资源站数据进行更新操作
			BatchCollect(ft.Time, ft.Ids...)
		}
		// 任务执行完毕
		log.Printf("执行一次定时任务: Task[%s]\n", ft.Id)
	})
}

// AddAutoUpdateCron 自动更新定时任务
func AddAutoUpdateCron(id, spec string) (cron.EntryID, error) {
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
		if ft.State && ft.Model == 0 {
			AutoCollect(ft.Time)
			log.Println("执行一次自动更新任务")
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
	log.Printf("%+v\n", CronCollect.Entries())
	log.Println("", CronCollect.Entry(id).Next.Format(time.DateTime))
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

// 清理API接口数据缓存
func clearCache() {
	system.RemoveCache(config.IndexCacheKey)
}
