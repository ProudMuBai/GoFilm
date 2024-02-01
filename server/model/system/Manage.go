package system

import (
	"encoding/json"
	"log"
	"server/config"
	"server/plugin/db"
)

type BasicConfig struct {
	SiteName string `json:"siteName"` // 网站名称
	Domain   string `json:"domain"`   // 网站域名
	Logo     string `json:"logo"`     // 网站logo
	Keyword  string `json:"keyword"`  // seo关键字
	Describe string `json:"describe"` // 网站描述信息
	State    bool   `json:"state"`    // 网站状态 开启 || 关闭
	Hint     string `json:"hint"`     // 网站关闭提示
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
