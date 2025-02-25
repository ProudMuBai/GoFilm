package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/model/system"
	"server/plugin/SystemInit"
	"server/plugin/common/util"
)

func ManageIndex(c *gin.Context) {
	system.SuccessOnlyMsg("后台管理中心", c)
	return
}

// ------------------------------------------------------ 站点基本配置 ------------------------------------------------------

// SiteBasicConfig  网站基本配置
func SiteBasicConfig(c *gin.Context) {
	system.Success(logic.ML.GetSiteBasicConfig(), "网站基本信息获取成功", c)
}

// UpdateSiteBasic 更新网站配置信息
func UpdateSiteBasic(c *gin.Context) {
	// 获取请求参数 && 校验关键配置项
	bc := system.BasicConfig{}
	if err := c.ShouldBindJSON(&bc); err == nil {
		// 对参数进行校验
		if !util.ValidDomain(bc.Domain) && !util.ValidIPHost(bc.Domain) {
			system.Failed("域名格式校验失败", c)
			return
		}
		if len(bc.SiteName) <= 0 {
			system.Failed("网站名称不能为空", c)
			return
		}
	} else {
		system.Failed(fmt.Sprint("请求参数异常:  ", err), c)
		return
	}

	// 保存更新后的配置信息
	if err := logic.ML.UpdateSiteBasic(bc); err != nil {
		system.Failed(fmt.Sprint("网站配置更新失败:  ", err), c)
		return
	}
	system.SuccessOnlyMsg("更新成功", c)
}

// ResetSiteBasic 重置网站配置信息为初始化状态
func ResetSiteBasic(c *gin.Context) {
	// 执行配置初始化方法直接覆盖当前基本配置信息
	SystemInit.BasicConfigInit()
	system.SuccessOnlyMsg("配置信息重置成功", c)
}

// ------------------------------------------------------ 轮播数据配置 ------------------------------------------------------

// BannerList 获取轮播图数据
func BannerList(c *gin.Context) {
	bl := logic.ML.GetBanners()
	system.Success(bl, "配置信息重置成功", c)
}

// BannerFind 返回ID对应的横幅信息
func BannerFind(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		system.Failed("Banner信息获取失败, ID信息异常", c)
		return
	}
	bl := logic.ML.GetBanners()
	for _, b := range bl {
		if b.Id == id {
			system.Success(b, "Banner信息获取成功", c)
			return
		}
	}
	system.Failed("Banner信息获取失败", c)
}

// BannerAdd  添加海报数据
func BannerAdd(c *gin.Context) {
	var b system.Banner
	if err := c.ShouldBindJSON(&b); err != nil {
		system.Failed("Banner参数提交异常", c)
		return
	}
	// 为新增的banner生成Id
	b.Id = util.GenerateSalt()
	bl := logic.ML.GetBanners()
	if len(bl) > 6 {
		system.Failed("Banners最大阈值为6, 无法添加新的banner信息", c)
		return
	}
	bl = append(bl, b)
	if err := logic.ML.SaveBanners(bl); err != nil {
		system.Failed(fmt.Sprintln("Banners信息添加失败,", err), c)
		return
	}
	system.SuccessOnlyMsg("海报信息添加成功", c)
}

// BannerUpdate  更新海报数据
func BannerUpdate(c *gin.Context) {
	var banner system.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		system.Failed("Banner参数提交异常", c)
		return
	}
	bl := logic.ML.GetBanners()
	for i, b := range bl {
		if b.Id == banner.Id {
			bl[i] = banner
			if err := logic.ML.SaveBanners(bl); err != nil {
				system.Failed("海报信息更新失败", c)
			} else {
				system.SuccessOnlyMsg("海报信息更新成功", c)
				return
			}

		}
	}
	system.Failed("海报信息更新失败, 未匹配对应Banner信息", c)
}

// BannerDel 删除海报数据
func BannerDel(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		system.Failed("Banner信息获取失败, ID信息异常", c)
		return
	}
	bl := logic.ML.GetBanners()
	for i, b := range bl {
		if b.Id == id {
			bl = append(bl[:i], bl[i+1:]...)
			_ = logic.ML.SaveBanners(bl)
			system.SuccessOnlyMsg("海报信息删除成功", c)
			return
		}
	}
	system.Failed("海报信息删除失败", c)
}

// ------------------------------------------------------ 参数校验 ------------------------------------------------------
func validFilmSource(fs system.FilmSource) error {
	// 资源名称不能为空 且长度不能超过20
	if len(fs.Name) <= 0 || len(fs.Name) > 20 {
		return errors.New("资源名称不能为空且长度不能超过20")
	}
	// Uri 采集链接测试格式
	if !util.ValidURL(fs.Uri) {
		return errors.New("资源链接格式异常, 请输入规范的URL链接")
	}
	// 校验接口类型是否是 JSON || XML
	if fs.ResultModel != system.JsonResult && fs.ResultModel != system.XmlResult {
		return errors.New("接口类型异常, 请提交正确的接口类型")
	}
	// 校验采集类型是否符合规范
	switch fs.CollectType {
	case system.CollectVideo, system.CollectArticle, system.CollectActor, system.CollectRole, system.CollectWebSite:
		return nil
	default:
		return errors.New("资源类型异常, 未知的资源类型")
	}
}
