package logic

import (
	"errors"
	"server/model/system"
)

type ManageLogic struct {
}

var ML *ManageLogic

// GetFilmSourceList 获取采集站列表数据
func (ml *ManageLogic) GetFilmSourceList() []system.FilmSource {
	// 返回当前已添加的采集站列表信息
	return system.GetCollectSourceList()
}

// GetFilmSource 获取ID对应的采集源信息
func (ml *ManageLogic) GetFilmSource(id string) *system.FilmSource {
	return system.FindCollectSourceById(id)
}

// UpdateFilmSource 更新采集源信息
func (ml *ManageLogic) UpdateFilmSource(s system.FilmSource) error {
	return system.UpdateCollectSource(s)
}

// SaveFilmSource  保存采集源信息
func (ml *ManageLogic) SaveFilmSource(s system.FilmSource) error {
	return system.AddCollectSource(s)
}

// DelFilmSource  删除采集源信息
func (ml *ManageLogic) DelFilmSource(id string) error {
	// 先查找是否存在对应ID的站点信息
	s := system.FindCollectSourceById(id)
	if s == nil {
		return errors.New("当前资源站信息不存在, 请勿重复操作")
	}
	//  如果是主站点则返回提示禁止直接删除
	if s.Grade == system.MasterCollect {
		return errors.New("主站点无法直接删除, 请先降级为附属站点再进行删除")
	}
	system.DelCollectResource(id)
	return nil
}

// GetSiteBasicConfig 获取网站基本配置信息
func (ml *ManageLogic) GetSiteBasicConfig() system.BasicConfig {
	return system.GetSiteBasic()
}

// UpdateSiteBasic 更新网站配置信息
func (ml *ManageLogic) UpdateSiteBasic(c system.BasicConfig) error {
	return system.SaveSiteBasic(c)
}

// GetBanners 获取轮播组件信息
func (ml *ManageLogic) GetBanners() system.Banners {
	return system.GetBanners()
}

// SaveBanners 保存轮播信息
func (ml *ManageLogic) SaveBanners(bl system.Banners) error {
	return system.SaveBanners(bl)
}
