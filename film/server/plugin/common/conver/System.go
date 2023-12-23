package conver

import (
	"server/model/system"
	"time"
)

/*
	系统内部对象想换转换
*/

// CovertFilmDetailVo 将 FilmDetailVo 转化为 MovieDetail
func CovertFilmDetailVo(fd system.FilmDetailVo) (system.MovieDetail, error) {
	t, err := time.ParseInLocation(time.DateTime, fd.AddTime, time.Local)
	md := system.MovieDetail{
		Id:       fd.Id,
		Cid:      fd.Cid,
		Pid:      fd.Pid,
		Name:     fd.Name,
		Picture:  fd.Picture,
		DownFrom: fd.DownFrom,
		MovieDescriptor: system.MovieDescriptor{
			SubTitle:    fd.SubTitle,
			CName:       fd.CName,
			EnName:      fd.EnName,
			Initial:     fd.Initial,
			ClassTag:    fd.ClassTag,
			Actor:       fd.Actor,
			Director:    fd.Director,
			Writer:      fd.Writer,
			Blurb:       fd.Content,
			Remarks:     fd.Remarks,
			ReleaseDate: fd.ReleaseDate,
			Area:        fd.Area,
			Language:    fd.Language,
			Year:        fd.Year,
			State:       fd.State,
			UpdateTime:  fd.UpdateTime,
			AddTime:     t.Unix(),
			DbId:        fd.DbId,
			DbScore:     fd.DbScore,
			Hits:        fd.Hits,
			Content:     fd.Content,
		},
	}
	// 通过分割符切分播放源信息  PlaySeparator $$$
	//md.PlayFrom = strings.Split(fd.VodPlayFrom, fd.VodPlayNote)
	// v2 只保留m3u8播放源
	md.PlayList = GenFilmPlayList(fd.PlayLink, "$$$")
	//md.DownloadList = GenFilmPlayList(fd.DownloadLink, fd.VodPlayNote)

	return md, err
}
