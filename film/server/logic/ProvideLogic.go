package logic

import (
	"fmt"
	"log"
	"server/config"
	"server/model/collect"
	"server/model/system"
	"server/plugin/common/conver"
	"strconv"
	"strings"
)

type ProvideLogic struct {
}

var PL *ProvideLogic

// GetFilmDetailPage 处理请求参数, 返回filmDetail数据
func (pl *ProvideLogic) GetFilmDetailPage(params map[string]string, page *system.Page) collect.FilmDetailLPage {
	return filmDetailPage(params, page)
}

// GetFilmListPage 处理请求参数, 返回filmList数据
func (pl *ProvideLogic) GetFilmListPage(params map[string]string, page *system.Page) collect.FilmListPage {
	dp := filmDetailPage(params, page)
	var p collect.FilmListPage = collect.FilmListPage{
		Code:      dp.Code,
		Msg:       dp.Msg,
		Page:      dp.Page,
		PageCount: dp.PageCount,
		Limit:     dp.Limit,
		Total:     dp.Total,
		List:      conver.DetailCovertList(dp.List),
		Class:     collect.GetFilmClass(),
	}
	return p
}

func (pl *ProvideLogic) GetFilmDetailXmlPage(params map[string]string, page *system.Page) collect.RssD {
	dp := filmDetailPage(params, page)
	var dxp = collect.RssD{
		Version: config.RssVersion,
		List: collect.FilmDetailPageX{
			Page:        fmt.Sprint(dp.Page),
			PageCount:   dp.PageCount,
			PageSize:    fmt.Sprint(dp.Limit),
			RecordCount: len(dp.List),
			Videos:      conver.DetailCovertXml(dp.List),
		},
	}
	return dxp
}

func (pl *ProvideLogic) GetFilmListXmlPage(params map[string]string, page *system.Page) collect.RssL {
	dp := filmDetailPage(params, page)
	cl := collect.GetFilmClass()
	var dxp = collect.RssL{
		Version: config.RssVersion,
		List: collect.FilmListPageX{
			Page:        dp.Page,
			PageCount:   dp.PageCount,
			PageSize:    dp.Limit,
			RecordCount: len(dp.List),
			Videos:      conver.DetailCovertListXml(dp.List),
		},
		ClassXL: conver.ClassListCovertXml(cl),
	}
	return dxp
}

func filmDetailPage(params map[string]string, page *system.Page) collect.FilmDetailLPage {
	var p collect.FilmDetailLPage = collect.FilmDetailLPage{
		Code: 1,
		Msg:  "数据列表",
		Page: fmt.Sprint(page.Current),
	}
	// 如果params中的ids不为空, 则直接返回ids对应的数据
	if len(params["ids"]) > 0 {
		var ids []int64
		for _, idStr := range strings.Split(params["ids"], ",") {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				ids = append(ids, id)
			}
		}
		page.Total = len(ids)
		page.PageCount = int((len(ids) + page.PageSize - 1) / page.PageSize)
		// 获取id对应的数据
		for i := 0; i >= (page.Current-1)*page.PageSize && i < page.Total && i < page.Current*page.PageSize; i++ {
			if fd, err := collect.GetOriginalDetailById(ids[i]); err == nil {
				p.List = append(p.List, conver.FilterFilmDetail(fd, 0))
			}
		}
		p.PageCount = page.PageCount
		p.Limit = fmt.Sprint(page.PageSize)
		p.Total = page.Total
		return p
	}

	// 如果请求参数中不包含 ids, 则通过条件进行对应查找
	l, err := system.FindFilmIds(params, page)
	if err != nil {
		log.Println(err)
	}
	for _, id := range l {
		if fd, e := collect.GetOriginalDetailById(id); e == nil {
			p.List = append(p.List, conver.FilterFilmDetail(fd, 0))
		}
	}
	p.PageCount = page.PageCount
	p.Limit = fmt.Sprint(page.PageSize)
	p.Total = page.Total
	return p
}
