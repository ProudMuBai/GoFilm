package conver

import (
	"encoding/xml"
	"log"
	"server/config"
	"server/model/collect"
	"server/model/system"
	"strings"
)

/*
	处理 不同结构体数据之间的转化
	统一转化为内部结构体
*/

// GenCategoryTree 解析处理 filmListPage数据 生成分类树形数据
func GenCategoryTree(list []collect.FilmClass) *system.CategoryTree {
	// 遍历所有分类进行树形结构组装
	tree := &system.CategoryTree{Category: &system.Category{Id: 0, Pid: -1, Name: "分类信息", Show: true}}
	temp := make(map[int64]*system.CategoryTree)
	temp[tree.Id] = tree
	for _, c := range list {
		// 判断当前节点ID是否存在于 temp中
		category, ok := temp[c.TypeID]
		if ok {
			// 将当前节点信息保存
			category.Category = &system.Category{Id: c.TypeID, Pid: c.TypePid, Name: c.TypeName, Show: true}
		} else {
			// 如果不存在则将当前分类存放到 temp中
			category = &system.CategoryTree{Category: &system.Category{Id: c.TypeID, Pid: c.TypePid, Name: c.TypeName, Show: true}}
			temp[c.TypeID] = category
		}
		// 根据 pid获取父节点信息
		parent, ok := temp[category.Pid]
		if !ok {
			// 如果不存在父节点存在, 则将父节点存放到temp中
			temp[c.TypePid] = parent
		}
		// 将当前节点存放到父节点的Children中
		parent.Children = append(parent.Children, category)
	}

	return tree
}

// ConvertCategoryList 将分类树形数据转化为list类型
func ConvertCategoryList(tree system.CategoryTree) []system.Category {
	var cl = []system.Category{system.Category{Id: tree.Id, Pid: tree.Pid, Name: tree.Name, Show: tree.Show}}
	for _, c := range tree.Children {
		cl = append(cl, system.Category{Id: c.Id, Pid: c.Pid, Name: c.Name, Show: c.Show})
		if c.Children != nil && len(c.Children) > 0 {
			for _, subC := range c.Children {
				cl = append(cl, system.Category{Id: subC.Id, Pid: subC.Pid, Name: subC.Name, Show: subC.Show})
			}
		}
	}
	return cl
}

// ConvertFilmDetails 批量处理影片详情信息
func ConvertFilmDetails(details []collect.FilmDetail) []system.MovieDetail {
	var dl []system.MovieDetail
	for _, d := range details {
		dl = append(dl, ConvertFilmDetail(d))
	}
	return dl

}

// ConvertFilmDetail 将影片详情数据处理转化为 system.MovieDetail
func ConvertFilmDetail(detail collect.FilmDetail) system.MovieDetail {
	md := system.MovieDetail{
		Id:       detail.VodID,
		Cid:      detail.TypeID,
		Pid:      detail.TypeID1,
		Name:     detail.VodName,
		Picture:  detail.VodPic,
		DownFrom: detail.VodDownFrom,
		MovieDescriptor: system.MovieDescriptor{
			SubTitle:    detail.VodSub,
			CName:       detail.TypeName,
			EnName:      detail.VodEn,
			Initial:     detail.VodLetter,
			ClassTag:    detail.VodClass,
			Actor:       detail.VodActor,
			Director:    detail.VodDirector,
			Writer:      detail.VodWriter,
			Blurb:       detail.VodBlurb,
			Remarks:     detail.VodRemarks,
			ReleaseDate: detail.VodPubDate,
			Area:        detail.VodArea,
			Language:    detail.VodLang,
			Year:        detail.VodYear,
			State:       detail.VodState,
			UpdateTime:  detail.VodTime,
			AddTime:     detail.VodTimeAdd,
			DbId:        detail.VodDouBanID,
			DbScore:     detail.VodDouBanScore,
			Hits:        detail.VodHits,
			Content:     detail.VodContent,
		},
	}
	// 通过分割符切分播放源信息  PlaySeparator $$$
	md.PlayFrom = strings.Split(detail.VodPlayFrom, detail.VodPlayNote)
	// v2 只保留m3u8播放源
	md.PlayList = GenFilmPlayList(detail.VodPlayURL, detail.VodPlayNote)
	md.DownloadList = GenFilmPlayList(detail.VodDownURL, detail.VodPlayNote)

	return md
}

// GenFilmPlayList 处理影片播放地址数据, 只保留m3u8与mp4格式的链接,生成playList
func GenFilmPlayList(playUrl, separator string) [][]system.MovieUrlInfo {
	var res [][]system.MovieUrlInfo
	if separator != "" {
		// 1. 通过分隔符切分播放源地址
		for _, l := range strings.Split(playUrl, separator) {
			// 2.只对m3u8播放源 和 .mp4下载地址进行处理
			if strings.Contains(l, ".m3u8") || strings.Contains(l, ".mp4") {
				// 2. 将每组播放源对应的播放列表信息存储到列表中
				res = append(res, ConvertPlayUrl(l))

			}
		}
	} else {
		// 1.只对m3u8播放源 和 .mp4下载地址进行处理
		if strings.Contains(playUrl, ".m3u8") || strings.Contains(playUrl, ".mp4") {
			// 2. 将每组播放源对应的播放列表信息存储到列表中
			res = append(res, ConvertPlayUrl(playUrl))
		}
	}
	return res
}

// GenAllFilmPlayList 处理影片播放地址数据, 保留全部播放链接,生成playList
func GenAllFilmPlayList(playUrl, separator string) [][]system.MovieUrlInfo {
	var res [][]system.MovieUrlInfo
	if separator != "" {
		// 1. 通过分隔符切分播放源地址
		for _, l := range strings.Split(playUrl, separator) {
			// 将playUrl中的所有播放格式链接均进行转换保存
			res = append(res, ConvertPlayUrl(l))
		}
		return res
	}
	// 将playUrl中的所有播放格式链接均进行转换保存
	res = append(res, ConvertPlayUrl(playUrl))
	return res
}

// ConvertPlayUrl 将单个playFrom的播放地址字符串处理成列表形式
func ConvertPlayUrl(playUrl string) []system.MovieUrlInfo {
	// 对每个片源的集数和播放地址进行分割 Episode$Link#Episode$Link
	var l []system.MovieUrlInfo
	for _, p := range strings.Split(playUrl, "#") {
		// 处理 Episode$Link 形式的播放信息
		if strings.Contains(p, "$") {
			l = append(l, system.MovieUrlInfo{
				Episode: strings.Split(p, "$")[0],
				Link:    strings.Split(p, "$")[1],
			})
		} else {
			l = append(l, system.MovieUrlInfo{
				Episode: "(｀・ω・´)",
				Link:    p,
			})
		}
	}
	return l
}

// ConvertVirtualPicture 将影片详情信息转化为虚拟图片信息
func ConvertVirtualPicture(details []system.MovieDetail) []system.VirtualPicture {
	var l []system.VirtualPicture
	for _, d := range details {
		if len(d.Picture) > 0 {
			l = append(l, system.VirtualPicture{Id: d.Id, Link: d.Picture})
		}
	}
	return l
}

// ----------------------------------Provide API---------------------------------------------------

// DetailCovertList 将影视详情信息转化为列表信息
func DetailCovertList(details []collect.FilmDetail) []collect.FilmList {
	var l []collect.FilmList
	for _, d := range details {
		fl := collect.FilmList{
			VodID:       d.VodID,
			VodName:     d.VodName,
			TypeID:      d.TypeID,
			TypeName:    d.TypeName,
			VodEn:       d.VodEn,
			VodTime:     d.VodTime,
			VodRemarks:  d.VodRemarks,
			VodPlayFrom: d.VodPlayFrom,
		}
		l = append(l, fl)
	}
	return l
}

// DetailCovertXml 将影片详情信息转化为Xml格式的对象
func DetailCovertXml(details []collect.FilmDetail) []collect.VideoDetail {
	var vl []collect.VideoDetail
	for _, d := range details {
		vl = append(vl, collect.VideoDetail{
			Last:     d.VodTime,
			ID:       d.VodID,
			Tid:      d.TypeID,
			Name:     collect.CDATA{Text: d.VodName},
			Type:     d.TypeName,
			Pic:      d.VodPic,
			Lang:     d.VodLang,
			Area:     d.VodArea,
			Year:     d.VodYear,
			State:    d.VodState,
			Note:     collect.CDATA{Text: d.VodRemarks},
			Actor:    collect.CDATA{Text: d.VodActor},
			Director: collect.CDATA{Text: d.VodDirector},
			DL:       collect.DL{DD: []collect.DD{collect.DD{Flag: d.VodPlayFrom, Value: d.VodPlayURL}}},
			Des:      collect.CDATA{Text: d.VodContent},
		})
	}
	return vl
}

// DetailCovertListXml 将影片详情信息转化为Xml格式FilmList的对象
func DetailCovertListXml(details []collect.FilmDetail) []collect.VideoList {
	var vl []collect.VideoList
	for _, d := range details {
		vl = append(vl, collect.VideoList{
			Last: d.VodTime,
			ID:   d.VodID,
			Tid:  d.TypeID,
			Name: collect.CDATA{Text: d.VodName},
			Type: d.TypeName,
			Dt:   d.VodPlayFrom,
			Note: collect.CDATA{Text: d.VodRemarks},
		})
	}
	s, _ := xml.Marshal(vl[0])
	log.Println(string(s))
	return vl
}

// ClassListCovertXml 将影片分类列表转化为XML格式
func ClassListCovertXml(cl []collect.FilmClass) collect.ClassXL {
	var l collect.ClassXL
	for _, c := range cl {
		l.ClassX = append(l.ClassX, collect.ClassX{ID: c.TypeID, Value: c.TypeName})
	}
	return l
}

// FilterFilmDetail 对影片详情数据进行处理, t 修饰类型 0-返回m3u8,mp4 | 1 返回 云播链接 | 2 返回全部
func FilterFilmDetail(fd collect.FilmDetail, t int64) collect.FilmDetail {
	// 只保留 mu38 | mp4 格式的播放源, 如果包含多种格式的播放数据
	if strings.Contains(fd.VodPlayURL, fd.VodPlayNote) {
		switch t {
		case 2:
			fd.VodPlayFrom = config.PlayFormAll
		case 1, 0:
			for _, v := range strings.Split(fd.VodPlayURL, fd.VodPlayNote) {
				if t == 0 && (strings.Contains(v, ".m3u8") || strings.Contains(v, ".mp4")) {
					fd.VodPlayFrom = config.PlayForm
					fd.VodPlayURL = v
				} else if t == 1 && !strings.Contains(v, ".m3u8") && !strings.Contains(v, ".mp4") {
					fd.VodPlayFrom = config.PlayFormCloud
					fd.VodPlayURL = v
				}
			}

		}
	} else {
		// 如果只有一种类型的播放链,则默认为m3u8  修改 VodPlayFrom 信息
		fd.VodPlayFrom = config.PlayForm
	}

	return fd
}
