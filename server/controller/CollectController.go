package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/model/system"
	"server/plugin/spider"
	"strconv"
)

// ------------------------------------------------------ 影视采集 ------------------------------------------------------

// FilmSourceList 采集站点信息
func FilmSourceList(c *gin.Context) {
	system.Success(logic.CollectL.GetFilmSourceList(), "影视源站点信息获取成功", c)
	return
}

// FindFilmSource 通过ID返回对应的资源站数据
func FindFilmSource(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	fs := logic.CollectL.GetFilmSource(id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	system.Success(fs, "原站点详情信息查找成功", c)
}

// FilmSourceAdd 添加采集源
func FilmSourceAdd(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		system.Failed("资源接口测试失败, 请确认接口有效再添加", c)
		return
	}
	// 测试通过后将资源站信息添加到list
	if err := logic.CollectL.SaveFilmSource(s); err != nil {
		system.Failed(fmt.Sprint("资源站添加失败: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("添加成功", c)
}

// FilmSourceUpdate 采集站点信息更新
func FilmSourceUpdate(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	// 校验Id信息是否为空
	if s.Id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	fs := logic.CollectL.GetFilmSource(s.Id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	// 如果 uri发生变更则执行spider测试
	if fs.Uri != s.Uri {
		// 执行 spider
		if err := spider.CollectApiTest(s); err != nil {
			system.Failed("资源接口测试失败, 请确认更新的数据接口是否有效", c)
			return
		}
	}
	// 更新资源站信息
	if err := logic.CollectL.UpdateFilmSource(s); err != nil {
		system.Failed(fmt.Sprint("资源站更新失败: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("更新成功", c)
}

// FilmSourceChange 采集站点状态变更
func FilmSourceChange(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	if s.Id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	// 查找对应的资源站点信息
	fs := logic.CollectL.GetFilmSource(s.Id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (fs.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	if s.State != fs.State || s.SyncPictures != fs.SyncPictures {
		// 执行更新操作
		s := system.FilmSource{Id: fs.Id, Name: fs.Name, Uri: fs.Uri, ResultModel: fs.ResultModel,
			Grade: fs.Grade, SyncPictures: s.SyncPictures, CollectType: fs.CollectType, State: s.State}
		// 更新资源站信息
		if err := logic.CollectL.UpdateFilmSource(s); err != nil {
			system.Failed(fmt.Sprint("资源站更新失败: ", err.Error()), c)
			return
		}
	}
	system.SuccessOnlyMsg("更新成功", c)
}

// FilmSourceDel 采集站点删除
func FilmSourceDel(c *gin.Context) {
	id := c.Query("id")
	if len(id) <= 0 {
		system.Failed("资源站ID信息不能为空", c)
		return
	}
	if err := logic.CollectL.DelFilmSource(id); err != nil {
		system.Failed("删除资源站失败", c)
		return
	}
	system.SuccessOnlyMsg("删除成功", c)
}

// FilmSourceTest 测试影视站点数据是否可用
func FilmSourceTest(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("测试成功!!!", c)
}

// GetNormalFilmSource 获取状态为启用的采集站信息
func GetNormalFilmSource(c *gin.Context) {
	// 获取所有的采集站信息
	var l []system.FilmTaskOptions
	for _, v := range logic.CollectL.GetFilmSourceList() {
		if v.State {
			l = append(l, system.FilmTaskOptions{Id: v.Id, Name: v.Name})
		}
	}
	system.Success(l, "影视源信息获取成功", c)
}

// ------------------------------------------------------ 失败采集记录 ------------------------------------------------------

// FailureRecordList 失效采集记录分页数据
func FailureRecordList(c *gin.Context) {
	// 数据返回对象
	var params = system.RecordRequestVo{Paging: &system.Page{}}
	var err error
	// 获取筛选条件
	params.OriginId = c.DefaultQuery("originId", "")
	params.CollectType, err = strconv.Atoi(c.DefaultQuery("collectType", "-1"))
	params.Hour, err = strconv.Atoi(c.DefaultQuery("hour", "0"))
	params.Status, err = strconv.Atoi(c.DefaultQuery("status", "-1"))

	// 分页参数
	params.Paging.Current, err = strconv.Atoi(c.DefaultQuery("current", "1"))
	params.Paging.PageSize, err = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		system.Failed("影片分页数据获取失败, 分页参数异常", c)
		return
	}
	// 如果分页数据超出指定范围则设置为默认值
	if params.Paging.PageSize <= 0 || params.Paging.PageSize > 500 {
		params.Paging.PageSize = 10
	}

	// 获取满足条件的分页数据
	list := logic.CollectL.GetRecordList(params)
	system.Success(gin.H{"params": params, "list": list}, "影片分页信息获取成功", c)
}

// CollectRecover 对失败的采集进行
func CollectRecover(c *gin.Context) {
	// 获取记录id
	id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil && id != 0 {
		system.Failed("采集重试开启失败, 采集记录ID参数异常", c)
		return
	}
	// 通过记录id对失败的采集进行恢复重试操作
	err = logic.CollectL.CollectRecover(id)
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("采集重试已开启, 请勿重复操作", c)
}
