package system

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"regexp"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
	"strings"
)

// FileInfo 图片信息对象
type FileInfo struct {
	gorm.Model
	Link        string `json:"link"`        // 图片链接
	Uid         int    `json:"uid"`         // 上传人ID
	RelevanceId int64  `json:"relevanceId"` // 关联资源ID
	Type        int    `json:"type"`        // 文件类型 (0 影片封面, 1 用户头像)
	Fid         string `json:"fid"`         // 图片唯一标识, 通常为文件名
	FileType    string `json:"fileType"`    // 文件类型, txt, png, jpg
	//Size        int    `json:"size"`        // 文件大小
}

// VirtualPicture 采集入站,待同步的图片信息
type VirtualPicture struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

//------------------------------------------------本地图库------------------------------------------------

// TableName 设置图片存储表的表名
func (f *FileInfo) TableName() string {
	return config.FileTableName
}

// StoragePath 获取文件的保存路径
func (f *FileInfo) StoragePath() string {
	var storage string
	switch f.FileType {
	case "jpeg", "jpg", "png", "webp":
		storage = strings.Replace(f.Link, config.FilmPictureAccess, fmt.Sprint(config.FilmPictureUploadDir, "/"), 1)
	default:
	}
	return storage
}

// CreateFileTable 创建图片关联信息存储表
func CreateFileTable() {
	// 如果不存在则创建表 并设置自增ID初始值为10000
	if !ExistFileTable() {
		err := db.Mdb.AutoMigrate(&FileInfo{})
		if err != nil {
			log.Println("Create Table FileInfo Failed: ", err)
		}
	}
}

// ExistFileTable 是否存在Picture表
func ExistFileTable() bool {
	// 1. 判断表中是否存在当前表
	return db.Mdb.Migrator().HasTable(&FileInfo{})
}

// SaveGallery 保存图片关联信息
func SaveGallery(f FileInfo) {
	db.Mdb.Create(&f)
}

// ExistFileInfoByRid 查找图片信息是否存在
func ExistFileInfoByRid(rid int64) bool {
	var count int64
	db.Mdb.Model(&FileInfo{}).Where("relevance_id = ?", rid).Count(&count)
	return count > 0
}

// GetFileInfoByRid 通过关联的资源id获取对应的图片信息
func GetFileInfoByRid(rid int64) FileInfo {
	var f FileInfo
	db.Mdb.Where("relevance_id = ?", rid).First(&f)
	return f
}

// GetFileInfoById 通过ID获取对应的图片信息
func GetFileInfoById(id uint) FileInfo {
	var f = FileInfo{}
	db.Mdb.First(&f, id)
	return f
}

// GetFileInfoPage 获取文件关联信息分页数据
func GetFileInfoPage(tl []string, page *Page) []FileInfo {
	var fl []FileInfo
	query := db.Mdb.Model(&FileInfo{}).Where("file_type IN ?", tl).Order("id DESC")
	// 获取分页相关参数
	GetPage(query, page)
	// 获取分页数据
	if err := query.Limit(page.PageSize).Offset((page.Current - 1) * page.PageSize).Find(&fl).Error; err != nil {
		log.Println(err)
		return nil
	}
	return fl
}

func DelFileInfo(id uint) {
	db.Mdb.Unscoped().Delete(&FileInfo{}, id)
}

//------------------------------------------------图片同步------------------------------------------------

// SaveVirtualPic 保存待同步的图片信息
func SaveVirtualPic(pl []VirtualPicture) error {
	// 保存对应的待同步图片信息
	var zl []redis.Z
	for _, p := range pl {
		// 首先查询 Gallery 表中是否存在当前ID对应的图片信息, 如果不存在则保存
		//if !ExistPictureByRid(p.Id) {
		//	m, _ := json.Marshal(p)
		//	zl = append(zl, redis.Z{Score: float64(p.Id), Member: m})
		//}

		// 只要开启图片同步则将图片信息存入待同步图片信息集合中, 是否同步图片交由真正同步到本地时进行决断
		m, _ := json.Marshal(p)
		zl = append(zl, redis.Z{Score: float64(p.Id), Member: m})
	}
	return db.Rdb.ZAdd(db.Cxt, config.VirtualPictureKey, zl...).Err()
}

// SyncFilmPicture 同步新采集入栈还未同步的图片
func SyncFilmPicture() {
	// 获取集合中的元素数量, 如果集合中没有元素则直接返回
	count := db.Rdb.ZCard(db.Cxt, config.VirtualPictureKey).Val()
	if count <= 0 {
		return
	}
	// 扫描待同步图片的信息, 每次扫描count条
	sl := db.Rdb.ZPopMax(db.Cxt, config.VirtualPictureKey, config.MaxScanCount).Val()
	if len(sl) <= 0 {
		return
	}
	// 获取 VirtualPicture
	for _, s := range sl {
		// 获取图片信息
		vp := VirtualPicture{}
		_ = json.Unmarshal([]byte(s.Member.(string)), &vp)
		// 判断当前影片是否已经同步过图片, 如果已经同步则直接跳过后续逻辑
		if ExistFileInfoByRid(vp.Id) {
			continue
		}
		// 将图片同步到服务器中
		fileName, err := util.SaveOnlineFile(vp.Link, config.FilmPictureUploadDir)
		if err != nil {
			continue
		}
		// 完成同步后将图片信息保存到 Gallery 中
		SaveGallery(FileInfo{
			Link:        fmt.Sprint(config.FilmPictureAccess, fileName),
			Uid:         config.UserIdInitialVal,
			RelevanceId: vp.Id,
			Type:        0,
			Fid:         regexp.MustCompile(`\.[^.]+$`).ReplaceAllString(fileName, ""),
			FileType:    strings.TrimPrefix(filepath.Ext(fileName), "."),
		})
	}
	// 递归执行直到图片暂存信息为空
	SyncFilmPicture()
}

// ReplaceDetailPic 将影片详情中的图片地址替换为自己的
func ReplaceDetailPic(d *MovieDetail) {
	// 查询影片对应的本地图片信息
	if ExistFileInfoByRid(d.Id) {
		// 如果存在关联的本地图片, 则查询对应的图片信息
		f := GetFileInfoByRid(d.Id)
		// 替换采集站的图片链接为本地链接
		d.Picture = f.Link
	}
}

// ReplaceBasicDetailPic 替换影片基本数据中的封面图为本地图片
func ReplaceBasicDetailPic(d *MovieBasicInfo) {
	// 查询影片对应的本地图片信息
	if ExistFileInfoByRid(d.Id) {
		// 如果存在关联的本地图片, 则查询对应的图片信息
		f := GetFileInfoByRid(d.Id)
		// 替换采集站的图片链接为本地链接
		d.Picture = f.Link
	}
}
