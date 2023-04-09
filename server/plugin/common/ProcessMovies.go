package common

import (
	"server/model"
	"strings"
)

// ProcessMovieListInfo 处理影片列表中的信息
func ProcessMovieListInfo(list []model.MovieInfo) []model.Movie {
	var movies []model.Movie
	for _, info := range list {
		movies = append(movies, model.Movie{
			Id:       info.Id,
			Name:     info.Name,
			Cid:      info.Cid,
			CName:    info.CName,
			EnName:   info.EnName,
			Time:     info.Time,
			Remarks:  info.Remarks,
			PlayFrom: info.PlayFrom,
		})
	}
	return movies
}

// ProcessMovieDetailList 处理影片详情列表数据
func ProcessMovieDetailList(list []model.MovieDetailInfo) []model.MovieDetail {
	var detailList []model.MovieDetail
	for _, d := range list {
		detailList = append(detailList, ProcessMovieDetail(d))
	}
	return detailList
}

// ProcessMovieDetail 处理单个影片详情信息
func ProcessMovieDetail(detail model.MovieDetailInfo) model.MovieDetail {
	md := model.MovieDetail{
		Id:       detail.Id,
		Cid:      detail.Cid,
		Pid:      detail.Pid,
		Name:     detail.Name,
		Picture:  detail.Pic,
		DownFrom: detail.DownFrom,
		MovieDescriptor: model.MovieDescriptor{
			SubTitle:    detail.SubTitle,
			CName:       detail.CName,
			EnName:      detail.EnName,
			Initial:     detail.Initial,
			ClassTag:    detail.ClassTag,
			Actor:       detail.Actor,
			Director:    detail.Director,
			Writer:      detail.Writer,
			Blurb:       detail.Blurb,
			Remarks:     detail.Remarks,
			ReleaseDate: detail.PubDate,
			Area:        detail.Area,
			Language:    detail.Language,
			Year:        detail.Year,
			State:       detail.State,
			UpdateTime:  detail.UpdateTime,
			AddTime:     detail.AddTime,
			DbId:        detail.DbId,
			DbScore:     detail.DbScore,
			Content:     detail.Content,
		},
	}
	// 通过分割符切分播放源信息  PlaySeparator $$$
	md.PlayFrom = strings.Split(detail.PlayFrom, detail.PlaySeparator)
	// v2 只保留m3u8播放源
	md.PlayList = ProcessPlayInfoV2(detail.PlayUrl, detail.PlaySeparator)
	md.DownloadList = ProcessPlayInfoV2(detail.DownUrl, detail.PlaySeparator)
	return md
}

// ProcessPlayInfo 处理影片播放数据信息
func ProcessPlayInfo(info, sparator string) [][]model.MovieUrlInfo {
	var res [][]model.MovieUrlInfo
	// 1. 通过分隔符区分多个片源数据
	for _, l := range strings.Split(info, sparator) {
		// 2.对每个片源的集数和播放地址进行分割
		var item []model.MovieUrlInfo
		for _, p := range strings.Split(l, "#") {
			// 3. 处理 Episode$Link 形式的播放信息
			if strings.Contains(p, "$") {
				item = append(item, model.MovieUrlInfo{
					Episode: strings.Split(p, "$")[0],
					Link:    strings.Split(p, "$")[1],
				})
			} else {
				item = append(item, model.MovieUrlInfo{
					Episode: "O(∩_∩)O",
					Link:    p,
				})
			}
		}
		// 3. 将每组播放源对应的播放列表信息存储到列表中
		res = append(res, item)
	}
	return res
}

// ProcessPlayInfoV2 处理影片信息方案二 只保留m3u8播放源
func ProcessPlayInfoV2(info, sparator string) [][]model.MovieUrlInfo {
	var res [][]model.MovieUrlInfo
	if sparator != "" {
		// 1. 通过分隔符切分播放源地址
		for _, l := range strings.Split(info, sparator) {
			// 只对m3u8播放源 和 .mp4下载地址进行处理
			if strings.Contains(l, ".m3u8") || strings.Contains(l, ".mp4") {
				// 2.对每个片源的集数和播放地址进行分割
				var item []model.MovieUrlInfo
				for _, p := range strings.Split(l, "#") {
					// 3. 处理 Episode$Link 形式的播放信息
					if strings.Contains(p, "$") {
						item = append(item, model.MovieUrlInfo{
							Episode: strings.Split(p, "$")[0],
							Link:    strings.Split(p, "$")[1],
						})
					} else {
						item = append(item, model.MovieUrlInfo{
							Episode: "O(∩_∩)O",
							Link:    p,
						})
					}
				}
				// 3. 将每组播放源对应的播放列表信息存储到列表中
				res = append(res, item)
			}
		}
	} else {
		// 只对m3u8播放源 和 .mp4下载地址进行处理
		if strings.Contains(info, ".m3u8") || strings.Contains(info, ".mp4") {
			// 2.对每个片源的集数和播放地址进行分割
			var item []model.MovieUrlInfo
			for _, p := range strings.Split(info, "#") {
				// 3. 处理 Episode$Link 形式的播放信息
				if strings.Contains(p, "$") {
					item = append(item, model.MovieUrlInfo{
						Episode: strings.Split(p, "$")[0],
						Link:    strings.Split(p, "$")[1],
					})
				} else {
					item = append(item, model.MovieUrlInfo{
						Episode: "O(∩_∩)O",
						Link:    p,
					})
				}
			}
			// 3. 将每组播放源对应的播放列表信息存储到列表中
			res = append(res, item)
		}
	}
	return res
}
