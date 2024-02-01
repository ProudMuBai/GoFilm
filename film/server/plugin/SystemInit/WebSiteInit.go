package SystemInit

import "server/model/system"

// SiteConfigInit 网站配置初始化
func SiteConfigInit() {

}

// BasicConfigInit 初始化网站基本配置信息
func BasicConfigInit() {
	var bc = system.BasicConfig{
		SiteName: "GoFilm",
		Domain:   "http://127.0.0.1:3600",
		Logo:     "https://s2.loli.net/2023/12/05/O2SEiUcMx5aWlv4.jpg",
		Keyword:  "在线视频, 免费观影",
		Describe: "自动采集, 多播放源集成,在线观影网站",
		State:    true,
		Hint:     "网站升级中, 暂时无法访问 !!!",
	}
	_ = system.SaveSiteBasic(bc)
}
