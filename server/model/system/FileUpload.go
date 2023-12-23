package system

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"regexp"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
)

// Picture 图片信息对象
type Picture struct {
	gorm.Model
	Link        string `json:"link"`        // 图片链接
	Uid         int    `json:"uid"`         // 上传人ID
	RelevanceId int64  `json:"relevanceId"` // 关联资源ID
	PicType     int    `json:"picType"`     // 图片类型 (0 影片封面, 1 用户头像)
	PicUid      string `json:"picUid"`      // 图片唯一标识, 通常为文件名
	//Size        int    `json:"size"`        // 图片大小
}

// VirtualPicture 采集入站,待同步的图片信息
type VirtualPicture struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

//------------------------------------------------本地图库------------------------------------------------

// TableName 设置图片存储表的表名
func (p *Picture) TableName() string {
	return config.PictureTableName
}

// CreatePictureTable 创建图片关联信息存储表
func CreatePictureTable() {
	// 如果不存在则创建表 并设置自增ID初始值为10000
	if !ExistPictureTable() {
		err := db.Mdb.AutoMigrate(&Picture{})
		if err != nil {
			log.Println("Create Table Picture Failed: ", err)
		}
	}
}

// ExistPictureTable 是否存在Picture表
func ExistPictureTable() bool {
	// 1. 判断表中是否存在当前表
	return db.Mdb.Migrator().HasTable(&Picture{})
}

// SaveGallery 保存图片关联信息
func SaveGallery(p Picture) {
	db.Mdb.Create(&p)
}

// ExistPictureByRid 查找图片信息是否存在
func ExistPictureByRid(rid int64) bool {
	var count int64
	db.Mdb.Model(&Picture{}).Where("relevance_id = ?", rid).Count(&count)
	return count > 0
}

// GetPictureByRid 通过关联的资源id获取对应的图片信息
func GetPictureByRid(rid int64) Picture {
	var p Picture
	db.Mdb.Where("relevance_id = ?", rid).First(&p)
	return p
}

func GetPicturePage(page *Page) []Picture {
	var pl []Picture
	query := db.Mdb.Model(&Picture{})
	// 获取分页相关参数
	GetPage(query, page)
	// 获取分页数据
	if err := query.Limit(page.PageSize).Offset((page.Current - 1) * page.PageSize).Find(&pl).Error; err != nil {
		log.Println(err)
		return nil
	}
	return pl
}

//------------------------------------------------图片同步------------------------------------------------

// SaveVirtualPic 保存待同步的图片信息
func SaveVirtualPic(pl []VirtualPicture) error {
	// 保存对应的
	var zl []redis.Z
	for _, p := range pl {
		// 首先查询 Gallery 表中是否存在当前ID对应的图片信息, 如果不存在则保存
		if !ExistPictureByRid(p.Id) {
			m, _ := json.Marshal(p)
			zl = append(zl, redis.Z{Score: float64(p.Id), Member: m})
		}
	}
	return db.Rdb.ZAdd(db.Cxt, config.VirtualPictureKey, zl...).Err()
}

// SyncFilmPicture 同步新采集入栈还未同步的图片
func SyncFilmPicture() {
	// 扫描待同步图片的信息, 每次扫描count条
	sl, cursor := db.Rdb.ZScan(db.Cxt, config.VirtualPictureKey, 0, "*", config.MaxScanCount).Val()
	if len(sl) <= 0 {
		return
	}
	// 获取 VirtualPicture
	for i, s := range sl {
		if i%2 == 0 {
			// 获取图片信息
			vp := VirtualPicture{}
			_ = json.Unmarshal([]byte(s), &vp)
			// 删除已经取出的数据
			db.Rdb.ZRem(db.Cxt, config.VirtualPictureKey, []byte(s))
			// 将图片同步到服务器
			fileName, err := util.SaveOnlineFile(vp.Link, config.FilmPictureUploadDir)
			if err != nil {
				continue
			}
			// 完成同步后将图片信息保存到 Gallery 中
			SaveGallery(Picture{
				Link:        fmt.Sprint(config.FilmPictureAccess, fileName),
				Uid:         config.UserIdInitialVal,
				RelevanceId: vp.Id,
				PicType:     0,
				PicUid:      regexp.MustCompile(`\.[^.]+$`).ReplaceAllString(fileName, ""),
			})
		}
	}
	// 如果 cursor != 0 则继续递归执行
	if cursor > 0 {
		SyncFilmPicture()
	}
}

// ReplaceDetailPic 将影片详情中的图片地址替换为自己的
func ReplaceDetailPic(d *MovieDetail) {
	// 查询影片对应的本地图片信息
	if ExistPictureByRid(d.Id) {
		// 如果存在关联的本地图片, 则查询对应的图片信息
		p := GetPictureByRid(d.Id)
		// 替换采集站的图片链接为本地链接
		d.Picture = p.Link
	}
}

// ReplaceBasicDetailPic 替换影片基本数据中的封面图为本地图片
func ReplaceBasicDetailPic(d *MovieBasicInfo) {
	// 查询影片对应的本地图片信息
	if ExistPictureByRid(d.Id) {
		// 如果存在关联的本地图片, 则查询对应的图片信息
		p := GetPictureByRid(d.Id)
		// 替换采集站的图片链接为本地链接
		d.Picture = p.Link
	}
}
