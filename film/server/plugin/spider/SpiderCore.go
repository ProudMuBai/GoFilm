package spider

import (
	"encoding/json"
	"errors"
	"log"
	"server/model/collect"
	"server/model/system"
	"server/plugin/common/conver"
	"server/plugin/common/util"
)

/*
	Spider 数据 爬取 & 处理 & 转换
*/

type FilmCollect interface {
	// GetCategoryTree 获取影视分类数据
	GetCategoryTree(r util.RequestInfo) (*system.CategoryTree, error)
	// GetPageCount 获取API接口的分页页数
	GetPageCount(r util.RequestInfo) (count int, err error)
	// GetFilmDetail 获取影片详情信息,返回影片详情列表
	GetFilmDetail(r util.RequestInfo) (list []system.MovieDetail, err error)
}

// ------------------------------------------------- JSON Collect -------------------------------------------------

// JsonCollect 处理返回值为JSON格式的采集数据
type JsonCollect struct {
}

// GetCategoryTree 获取分类树形数据
func (jc *JsonCollect) GetCategoryTree(r util.RequestInfo) (*system.CategoryTree, error) {
	// 设置请求参数信息
	r.Params.Set(`ac`, "list")
	r.Params.Set(`pg`, "1")
	// 执行请求, 获取一次list数据
	util.ApiGet(&r)
	// 解析resp数据
	filmListPage := collect.FilmListPage{}
	if len(r.Resp) <= 0 {
		log.Println("filmListPage 数据获取异常 : Resp Is Empty")
		return nil, errors.New("filmListPage 数据获取异常 : Resp Is Empty")
	}
	err := json.Unmarshal(r.Resp, &filmListPage)
	// 获取分类列表信息
	cl := filmListPage.Class
	// 组装分类数据信息树形结构
	tree := conver.GenCategoryTree(cl)

	// 将分类列表信息存储到redis
	_ = collect.SaveFilmClass(cl)

	return tree, err
}

// GetPageCount 获取分页总页数
func (jc *JsonCollect) GetPageCount(r util.RequestInfo) (count int, err error) {
	// 发送请求获取pageCount, 默认为获取 ac = detail
	if len(r.Params.Get("ac")) <= 0 {
		r.Params.Set("ac", "detail")
	}
	r.Params.Set("pg", "1")
	util.ApiGet(&r)
	//  判断请求结果是否为空, 如果为空直接输出错误并终止
	if len(r.Resp) <= 0 {
		err = errors.New("response is empty")
		return
	}
	// 获取pageCount
	res := collect.CommonPage{}
	err = json.Unmarshal(r.Resp, &res)
	if err != nil {
		return
	}
	count = int(res.PageCount)
	return
}

// GetFilmDetail 通过 RequestInfo 获取并解析出对应的 MovieDetail list
func (jc *JsonCollect) GetFilmDetail(r util.RequestInfo) (list []system.MovieDetail, err error) {
	// 防止json解析异常引发panic
	defer func() {
		if e := recover(); e != nil {
			log.Println("GetMovieDetail Failed : ", e)
		}
	}()
	// 设置分页请求参数
	r.Params.Set(`ac`, `detail`)
	util.ApiGet(&r)
	// 影视详情信息
	detailPage := collect.FilmDetailLPage{}
	//details := system.DetailListInfo{}
	// 如果返回数据为空则直接结束本次循环
	if len(r.Resp) <= 0 {
		err = errors.New("response is empty")
		return
	}
	// 序列化详情数据
	if err = json.Unmarshal(r.Resp, &detailPage); err != nil {
		return
	}

	// 将影视原始详情信息保存到redis中
	// 获取主站点uri
	//mc := system.GetCollectSourceListByGrade(system.MasterCollect)[0]
	//if mc.Uri == r.Uri {
	//	collect.BatchSaveOriginalDetail(detailPage.List)
	//}

	// 处理details信息
	list = conver.ConvertFilmDetails(detailPage.List)
	return
}

// ------------------------------------------------- XML Collect -------------------------------------------------

// XmlCollect 处理返回值为XML格式的采集数据
type XmlCollect struct {
}
