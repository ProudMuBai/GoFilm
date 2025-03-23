package logic

import (
	"server/model/system"
)

type ManageLogic struct {
}

var ML *ManageLogic

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
