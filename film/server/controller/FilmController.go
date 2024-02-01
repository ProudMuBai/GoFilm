package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/model/system"
	"strconv"
	"time"
)

// FilmSearchPage 获取影视分页数据
func FilmSearchPage(c *gin.Context) {
	var s = system.SearchVo{Paging: &system.Page{}}
	var err error
	// 检索参数
	s.Name = c.DefaultQuery("name", "")
	s.Pid, err = strconv.ParseInt(c.DefaultQuery("pid", "0"), 10, 64)
	if err != nil {
		system.Failed("影片分页数据获取失败, 请求参数异常", c)
		return
	}
	s.Cid, err = strconv.ParseInt(c.DefaultQuery("cid", "0"), 10, 64)
	if err != nil {
		system.Failed("影片分页数据获取失败, 请求参数异常", c)
		return
	}
	s.Plot = c.DefaultQuery("plot", "")
	s.Area = c.DefaultQuery("area", "")
	s.Language = c.DefaultQuery("language", "")
	year := c.DefaultQuery("year", "")
	if year == "" {
		s.Year = 0
	} else {
		s.Year, err = strconv.ParseInt(year, 10, 64)
		if err != nil {
			system.Failed("影片分页数据获取失败, 请求参数异常", c)
			return
		}
	}

	s.Remarks = c.DefaultQuery("remarks", "")
	// 处理时间参数
	begin := c.DefaultQuery("beginTime", "")
	if begin == "" {
		s.BeginTime = 0
	} else {
		beginTime, e := time.ParseInLocation(time.DateTime, begin, time.Local)
		if e != nil {
			system.Failed("影片分页数据获取失败, 请求参数异常", c)
			return
		}
		s.BeginTime = beginTime.Unix()
	}
	end := c.DefaultQuery("endTime", "")
	if end == "" {
		s.EndTime = 0
	} else {
		endTime, e := time.ParseInLocation(time.DateTime, end, time.Local)
		if e != nil {
			system.Failed("影片分页数据获取失败, 请求参数异常", c)
			return
		}
		s.EndTime = endTime.Unix()
	}

	// 分页参数
	s.Paging.Current, err = strconv.Atoi(c.DefaultQuery("current", "1"))
	s.Paging.PageSize, err = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// 如果分页数据超出指定范围则设置为默认值
	if s.Paging.PageSize <= 0 || s.Paging.PageSize > 500 {
		s.Paging.PageSize = 10
	}
	if err != nil {
		system.Failed("影片分页数据获取失败, 请求参数异常", c)
		return
	}
	// 提供检索tag options
	options := logic.FL.GetSearchOptions()
	// 检索条件
	sl := logic.FL.GetFilmPage(s)
	system.Success(gin.H{
		"params":  s,
		"list":    sl,
		"options": options,
	}, "影片分页信息获取成功", c)
}

// FilmAdd 手动添加影片
func FilmAdd(c *gin.Context) {
	// 获取请求参数
	var fd = system.FilmDetailVo{}
	if err := c.ShouldBindJSON(&fd); err != nil {
		system.Failed("影片添加失败, 影片参数提交异常", c)
		return
	}

	// 如果绑定成功则执行影片信息处理保存逻辑
	if err := logic.FL.SaveFilmDetail(fd); err != nil {
		system.Failed(fmt.Sprint("影片添加失败, 影片信息保存错误: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("影片信息添加成功", c)
}

//----------------------------------------------------影片分类处理----------------------------------------------------

// FilmClassTree 影片分类树数据
func FilmClassTree(c *gin.Context) {
	// 获取影片分类树信息
	tree := logic.FL.GetFilmClassTree()
	system.Success(tree, "影片分类信息获取成功", c)
	return
}

// FindFilmClass 获取指定ID对应的影片分类信息
func FindFilmClass(c *gin.Context) {
	idStr := c.DefaultQuery("id", "")
	if idStr == "" {
		system.Failed("影片分类信息获取失败, 分类Id不能为空", c)
		return
	}
	// 转化id类型为int
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		system.Failed("影片分类信息获取失败, 参数分类Id格式异常", c)
		return
	}
	// 通过Id返回对应的分类信息
	class := logic.FL.GetFilmClassById(id)
	if class == nil {
		system.Failed("影片分类信息获取失败, 分类信息不存在", c)
		return
	}
	system.Success(class, "分类信息查找成功", c)
}

func UpdateFilmClass(c *gin.Context) {
	// 获取修改后的分类信息
	var class = system.CategoryTree{}
	if err := c.ShouldBindJSON(&class); err != nil {
		system.Failed("更新失败, 请求参数异常", c)
		return
	}
	if class.Id == 0 {
		system.Failed("更新失败, 分类Id缺失", c)
		return
	}
	// 修改分类信息
	if err := logic.FL.UpdateClass(class); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("影片分类信息更新成功", c)
}

// DelFilmClass 删除指定ID对应的影片分类
func DelFilmClass(c *gin.Context) {
	idStr := c.DefaultQuery("id", "")
	if idStr == "" {
		system.Failed("影片分类信息获取失败, 分类Id不能为空", c)
		return
	}
	// 转化id类型为int
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		system.Failed("影片分类信息获取失败, 参数分类Id格式异常", c)
		return
	}
	// 通过ID删除对应分类信息
	if err = logic.FL.DelClass(id); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("当前分类已删除成功", c)
}
