package logic

import (
	"fmt"
	"log"
	"server/config"
	"server/model"
	"server/plugin/spider"
)

type SpiderLogic struct {
}

var SL *SpiderLogic

// ReZero 清空所有数据从零开始拉取
func (sl *SpiderLogic) ReZero() {
	// 如果指令正确,则执行重制
	spider.StartSpiderRe()
}

// FixDetail 重新获取主站点数据信息
func (sl *SpiderLogic) FixDetail() {
	spider.MainSiteSpider()
	log.Println("FilmDetail 重制完成!!!")
	// 先截断表中的数据
	model.TunCateSearchTable()
	// 重新扫描完整的信息到mysql中
	spider.SearchInfoToMdb()
	log.Println("SearchInfo 重制完成!!!")
}

// SpiderMtPlayRe 多站点播放数据清空重新获取
func (sl *SpiderLogic) SpiderMtPlayRe() {
	// 先清空有所附加播放源
	var keys []string
	for _, site := range spider.SiteList {
		keys = append(keys, fmt.Sprintf(config.MultipleSiteDetail, site.Name))
	}
	model.DelMtPlay(keys)
	// 如果指令正确,则执行详情数据获取
	spider.MtSiteSpider()
	log.Println("MtSiteSpider 重制完成!!!")
}
