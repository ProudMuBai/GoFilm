package SystemInit

import (
	"server/model/system"
	"server/plugin/common/util"
)

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

func BannersInit() {
	var bl = system.Banners{
		system.Banner{Id: util.GenerateSalt(), Name: "樱花庄的宠物女孩", Year: 2020, CName: "日韩动漫", Poster: "https://s2.loli.net/2024/02/21/Wt1QDhabdEI7HcL.jpg", Picture: "https://img.bfzypic.com/upload/vod/20230424-43/06e79232a4650aea00f7476356a49847.jpg", Remark: "已完结"},
		system.Banner{Id: util.GenerateSalt(), Name: "从零开始的异世界生活", Year: 2020, CName: "日韩动漫", Poster: "https://s2.loli.net/2024/02/21/UkpdhIRO12fsy6C.jpg", Picture: "https://img.bfzypic.com/upload/vod/20230424-43/06e79232a4650aea00f7476356a49847.jpg", Remark: "已完结"},
		system.Banner{Id: util.GenerateSalt(), Name: "五等分的花嫁", Year: 2020, CName: "日韩动漫", Poster: "https://s2.loli.net/2024/02/21/wXJr59Zuv4tcKNp.jpg", Picture: "https://img.bfzypic.com/upload/vod/20230424-43/06e79232a4650aea00f7476356a49847.jpg", Remark: "已完结"},
		system.Banner{Id: util.GenerateSalt(), Name: "我的青春恋爱物语果然有问题", Year: 2020, CName: "日韩动漫", Poster: "https://s2.loli.net/2024/02/21/oMAGzSliK2YbhRu.jpg", Picture: "https://img.bfzypic.com/upload/vod/20230424-43/06e79232a4650aea00f7476356a49847.jpg", Remark: "已完结"},
	}
	_ = system.SaveBanners(bl)
}
