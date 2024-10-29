package system

import (
	"encoding/json"
	"log"
	"server/config"
	"server/plugin/db"
	"sort"
)

// BasicConfig 网站基本信息
type BasicConfig struct {
	SiteName string `json:"siteName"` // 网站名称
	Domain   string `json:"domain"`   // 网站域名
	Logo     string `json:"logo"`     // 网站logo
	Keyword  string `json:"keyword"`  // seo关键字
	Describe string `json:"describe"` // 网站描述信息
	State    bool   `json:"state"`    // 网站状态 开启 || 关闭
	Hint     string `json:"hint"`     // 网站关闭提示
}

// Banner 首页横幅信息
type Banner struct {
	Id      int64  `json:"id"`      // 绑定所属影片Id
	Name    string `json:"name"`    // 影片名称
	Year    int64  `json:"year"`    // 上映年份
	CName   string `json:"cName"`   // 分类名称
	Poster  string `json:"poster"`  // 海报图片链接
	Picture string `json:"picture"` // 横幅大图链接
	Remark  string `json:"remark"`  // 更新状态描述信息
	Sort    int64  `json:"sort"`    // 排序分值
}

type Banners []Banner

func (bl Banners) Len() int {
	return len(bl)
}
func (bl Banners) Less(i, j int) bool {
	return bl[i].Sort < bl[j].Sort
}
func (bl Banners) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}

// ------------------------------------------------------ Redis ------------------------------------------------------

// SaveSiteBasic 保存网站基本配置信息
func SaveSiteBasic(c BasicConfig) error {
	data, _ := json.Marshal(c)
	return db.Rdb.Set(db.Cxt, config.SiteConfigBasic, data, config.ManageConfigExpired).Err()
}

// GetSiteBasic 获取网站基本配置信息
func GetSiteBasic() BasicConfig {
	c := BasicConfig{}
	data := db.Rdb.Get(db.Cxt, config.SiteConfigBasic).Val()
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		log.Println("GetSiteBasic Err", err)
	}
	return c
}

// GetBanners 获取轮播配置信息
func GetBanners() Banners {
	var bl Banners
	data := db.Rdb.Get(db.Cxt, config.BannersKey).Val()
	if err := json.Unmarshal([]byte(data), &bl); err != nil {
		log.Println("Get Banners Error", err)
	}
	// 通过 sort 对banners进行排序
	sort.Sort(bl)
	return bl
}

// SaveBanners 保存轮播配置信息
func SaveBanners(bl Banners) error {
	data, _ := json.Marshal(bl)
	return db.Rdb.Set(db.Cxt, config.BannersKey, data, config.ManageConfigExpired).Err()
}
