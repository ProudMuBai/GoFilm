package logic

import (
	"errors"
	"log"
	"server/model/system"
	"server/plugin/spider"
)

type SpiderLogic struct {
}

var SL *SpiderLogic

// BatchCollect 批量采集
func (sl *SpiderLogic) BatchCollect(time int, ids []string) {
	go spider.BatchCollect(time, ids...)
}

// StartCollect 执行对指定站点的采集任务
func (sl *SpiderLogic) StartCollect(id string, h int) error {
	// 先判断采集站是否存在于系统数据中
	if fs := system.FindCollectSourceById(id); fs == nil {
		return errors.New("采集任务开启失败采集站信息不存在")
	}
	// 存在则开启协程执行采集方法
	go func() {
		err := spider.HandleCollect(id, h)
		if err != nil {
			log.Printf("资源站[%s]采集任务执行失败: %s", id, err)
		}
	}()
	return nil
}

// AutoCollect 自动采集
func (sl *SpiderLogic) AutoCollect(time int) {
	go spider.AutoCollect(time)
}

// ZeroCollect 数据清除从零开始采集
func (sl *SpiderLogic) ZeroCollect(time int) {
	go spider.StarZero(time)
}

// FilmClassCollect 影视分类采集, 直接覆盖当前分类数据
func (sl *SpiderLogic) FilmClassCollect() error {
	l := system.GetCollectSourceListByGrade(system.MasterCollect)
	if l == nil {
		return errors.New("未获取到主采集站信息")
	}
	// 获取主站点信息, 只取第一条有效
	for _, fs := range l {
		if fs.State {
			go spider.CollectCategory(&fs)
			return nil
		}
	}
	return errors.New("未获取到已启用的主采集站信息")
}
