package logic

import (
	"fmt"
	"path/filepath"
	"server/config"
	"server/model/system"
	"strings"
)

type FileLogic struct {
}

var FileL FileLogic

func (fl *FileLogic) SingleFileUpload(fileName string, uid int) string {
	// 生成图片信息
	var p = system.Picture{Link: fmt.Sprint(config.FilmPictureAccess, filepath.Base(fileName)), Uid: uid, PicType: 0}
	p.PicUid = strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	// 记录图片信息到系统表中
	system.SaveGallery(p)
	return p.Link
}

// GetPhotoPage 获取系统内的图片分页信息
func (fl *FileLogic) GetPhotoPage(page *system.Page) []system.Picture {
	return system.GetPicturePage(page)

}
