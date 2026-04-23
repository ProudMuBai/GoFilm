package system

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"regexp"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Movie 影片基本信息
type Movie struct {
	Id       int64  `json:"id"`       // 影片ID
	Name     string `json:"name"`     // 影片名
	Cid      int64  `json:"cid"`      // 所属分类ID
	CName    string `json:"CName"`    // 所属分类名称
	EnName   string `json:"enName"`   // 英文片名
	Time     string `json:"time"`     // 更新时间
	Remarks  string `json:"remarks"`  // 备注 | 清晰度
	PlayFrom string `json:"playFrom"` // 播放来源
}

// MovieBasicInfo 影片基本信息
type MovieBasicInfo struct {
	Id       int64  `json:"id"`       //影片Id
	Cid      int64  `json:"cid"`      //分类ID
	Pid      int64  `json:"pid"`      //一级分类ID
	Name     string `json:"name"`     //片名
	SubTitle string `json:"subTitle"` //子标题
	CName    string `json:"cName"`    //分类名称
	State    string `json:"state"`    //影片状态 正片|预告...
	Picture  string `json:"picture"`  //简介图片
	Actor    string `json:"actor"`    //主演
	Director string `json:"director"` //导演
	Blurb    string `json:"blurb"`    //简介, 不完整
	Remarks  string `json:"remarks"`  // 更新情况
	Area     string `json:"area"`     // 地区
	Year     string `json:"year"`     //年份
}

// PlayItem 影视资源url信息
type PlayItem struct {
	Episode string `json:"episode"` // 集数
	Link    string `json:"link"`    // 播放地址
}

// MoviePlayList 播放列表信息, 二维切片
type MoviePlayList [][]PlayItem

// FromList 播放来源切片
type FromList []string

// MovieDetail 影片详情信息
type MovieDetail struct {
	Id          int64    `json:"id" gorm:"primaryKey"`      //影片Id
	Mid         int64    `json:"mid"`                       //影片Id
	Cid         int64    `json:"cid"`                       //分类ID
	Pid         int64    `json:"pid"`                       //一级分类ID
	Name        string   `json:"name"`                      //片名
	Picture     string   `json:"picture"`                   //简介图片
	SubTitle    string   `json:"subTitle"`                  //子标题
	CName       string   `json:"cName"`                     //分类名称
	EnName      string   `json:"enName"`                    //英文名
	Initial     string   `json:"initial"`                   //首字母
	ClassTag    string   `json:"classTag"`                  //分类标签
	Actor       string   `json:"actor"`                     //主演
	Director    string   `json:"director"`                  //导演
	Writer      string   `json:"writer"`                    //作者
	Blurb       string   `json:"blurb" gorm:"type:text"`    //简介, 残缺,不建议使用
	Remarks     string   `json:"remarks"`                   // 更新情况
	ReleaseDate string   `json:"releaseDate"`               //上映时间
	Area        string   `json:"area"`                      // 地区
	Language    string   `json:"language"`                  //语言
	Year        string   `json:"year"`                      //年份
	State       string   `json:"state"`                     //影片状态 正片|预告...
	UpdateTime  string   `json:"updateTime"`                //更新时间
	AddTime     int64    `json:"addTime"`                   //资源添加时间戳
	DbId        int64    `json:"dbId"`                      //豆瓣id
	DbScore     string   `json:"dbScore"`                   // 豆瓣评分
	Hits        int64    `json:"hits"`                      //影片热度
	Content     string   `json:"content" gorm:"type:text"`  //内容简介
	PlayFrom    FromList `json:"playFrom" gorm:"type:json"` // 播放来源
	DownFrom    string   `json:"DownFrom"`                  //下载来源 例: http
	//PlaySeparator   string              `json:"playSeparator"` // 播放信息分隔符
	PlayList     MoviePlayList `json:"playList" gorm:"type:json"`     //播放地址url
	DownloadList MoviePlayList `json:"downloadList" gorm:"type:json"` // 下载url地址
}

type SlaveMovieInfo struct {
	Id  int64  `json:"id" gorm:"primaryKey"` // 自增ID
	Sid string `json:"sid"`                  // 采集站标识ID
	//Name      string  `json:"name"`	// 影片名称
	Mid      string        `json:"mid"`  // 归一匹配ID
	DbId     int64         `json:"dbId"` //豆瓣ID 可能为空
	PlayList MoviePlayList `json:"playList" gorm:"type:json"`
}

// TableName 设置MovieDetail表的表名
func (m *MovieDetail) TableName() string {
	return config.MovieDetailName
}

// TableName 设置slaveMovieInfo 表名
func (m *SlaveMovieInfo) TableName() string {
	return config.SlaveMovieInfo
}

// ================================= 数据表处理 =================================

// CreateMovieDetailTable 创建存储检索信息的数据表
func CreateMovieDetailTable() {
	// 如果不存在则创建表 并设置自增ID初始值为10000
	if !ExistMovieDetailTable() {
		err := db.Mdb.AutoMigrate(&MovieDetail{})
		if err != nil {
			log.Println("Create Table MovieDetailsTable Failed: ", err)
		}
	}
}

// ExistMovieDetailTable 检测是否存在 MovieDetails表
func ExistMovieDetailTable() bool {
	return db.Mdb.Migrator().HasTable(&MovieDetail{})
}

// ResetMovieDetailTable 重置 MovieDetailTable
func ResetMovieDetailTable() {
	var m MovieDetail
	err := db.Mdb.Exec(fmt.Sprintf("TRUNCATE TABLE %s", m.TableName())).Error
	if err != nil {
		log.Println("MovieDetailTable Reset Error: ", err)
	}
}

// CreateSlaveMovieInfoTable 创建存储检索信息的数据表
func CreateSlaveMovieInfoTable() {
	// 如果不存在则创建表 并设置自增ID初始值为10000
	if !ExistSlaveMovieInfoTable() {
		err := db.Mdb.AutoMigrate(&SlaveMovieInfo{})
		if err != nil {
			log.Println("Create Table SlaveMovieInfoTable Failed: ", err)
		}
	}
}

// ExistSlaveMovieInfoTable 检测是否存在 MovieDetails表
func ExistSlaveMovieInfoTable() bool {
	return db.Mdb.Migrator().HasTable(&SlaveMovieInfo{})
}

// ResetSlaveMovieInfoTable 重置 SlaveMovieInfoTable (附属站点数据表一般不会单独重置)
func ResetSlaveMovieInfoTable() {
	var s SlaveMovieInfo
	err := db.Mdb.Exec(fmt.Sprintf("TRUNCATE TABLE %s", s.TableName())).Error
	if err != nil {
		log.Println("SlaveMovieInfoTable Reset Error: ", err)
	}
}

// DelSlaveMovieInfos 删除表中对应站点的数据信息
func DelSlaveMovieInfos(id string) {
	// 一次删除过多数据会锁表, 因此直接截断表

	//if err := db.Mdb.Where("sid = ?", id).Delete(&SlaveMovieInfo{}).Error; err != nil {
	//	log.Println("Delete SlaveMovieInfos  Error: ", err)
	//}
}

// AddMovieDetailIndex 添加详情表索引
func AddMovieDetailIndex() {
	var m MovieDetail
	tableName := m.TableName()
	// 添加索引
	db.Mdb.Exec(fmt.Sprintf("CREATE UNIQUE INDEX idx_mid ON %s (mid)", tableName))
}

// AddSlaveMovieInfoIndex 添加附属站点信息表索引
func AddSlaveMovieInfoIndex() {
	var s SlaveMovieInfo
	tableName := s.TableName()
	// 如果不存在索引则创建对应索引
	if !db.Mdb.Migrator().HasIndex(&s, "idx_mid") {
		// 添加索引
		//db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_sid ON %s (sid DESC)", tableName))
		//db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_mid ON %s (mid DESC)", tableName))
		//db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_dbId ON %s (db_id DESC, sid DESC)", tableName))
		db.Mdb.Exec(fmt.Sprintf("ALTER TABLE %s ADD INDEX idx_mid (mid DESC),ADD INDEX idx_ds (db_id DESC, sid DESC)", tableName))
	}
}

// =================================== column序列化 接口========================================================

func (m *MoviePlayList) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("MoviePlayList serialization failed, value is not []byte")
	}
	return json.Unmarshal(b, m)
}

func (m MoviePlayList) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (fl *FromList) Scan(value interface{}) error {
	if value == nil {
		*fl = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("FromList serialization failed, value is not []byte")
	}
	return json.Unmarshal(b, fl)
}

func (fl FromList) Value() (driver.Value, error) {
	if fl == nil {
		return nil, nil
	}
	return json.Marshal(fl)
}

// =================================== Spider数据处理 ========================================================

// SaveDetails 保存影片详情信息到redis中 格式: MovieDetail:Cid?:Id?
func SaveDetails(ml []MovieDetail) (err error) {
	// 1. 先将详情信息存入 MovieDetail表中
	if err = db.Mdb.Create(&ml).Error; err != nil {
		log.Println("影片详情信息保存失败: ", err)
	}
	// 2. 将详情信息转化为SearchInfo并保存
	BatchSaveSearchInfo(ml)
	return err
}

// SaveDetail 保存单部影片信息
func SaveDetail(m MovieDetail) (err error) {
	// 1. 转换 detail信息 searchInfo
	searchInfo := ConvertSearchInfo(m)
	// 2. 保存 Search tag 到 redis中 只存储用于检索对应影片的关键字信息
	SaveSearchTag(searchInfo)
	// 3. 将影片详情信息保存到 MovieDetails表中

	// 4. 先查询数据库中是否存在对应记录 ,如果不存在对应记录则 保存当前记录
	tx := db.Mdb.Begin()
	if !ExistMovieDetailByMid(m.Mid) {
		// 执行插入操作
		if err := tx.Create(&m).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 只对会变化的字段进行更新
		err := tx.Model(&MovieDetail{}).Where("mid", m.Mid).Updates(MovieDetail{PlayList: m.PlayList, DownloadList: m.DownloadList,
			Remarks: m.Remarks, State: m.State, UpdateTime: m.UpdateTime, AddTime: m.AddTime, DbScore: m.DbScore, Hits: m.Hits}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 提交事务
	tx.Commit()

	// 保存影片检索信息到searchTable
	err = SaveSearchInfo(searchInfo)
	return err

}

// BatchUpdateDetails 保存或更新detail数据
func BatchUpdateDetails(ml []MovieDetail) (err error) {
	// 先将details批量保存或更新
	for _, m := range ml {
		if !ExistMovieDetailByMid(m.Mid) {
			// 执行插入操作
			if err := db.Mdb.Create(&m).Error; err != nil {
				return err
			}
		} else {
			// 只对会变化的字段进行更新
			err := db.Mdb.Model(&MovieDetail{}).Where("mid", m.Mid).Updates(MovieDetail{PlayList: m.PlayList, DownloadList: m.DownloadList,
				Remarks: m.Remarks, State: m.State, UpdateTime: m.UpdateTime, AddTime: m.AddTime, DbScore: m.DbScore, Hits: m.Hits}).Error
			if err != nil {
				return err
			}
		}
		// 转化处理searchInfo信息s
		s := ConvertSearchInfo(m)
		// 保存searchInfo信息
		if err := SaveSearchInfo(s); err != nil {
			return err
		}
	}
	return err
}

// ExistMovieDetailByMid 通过mid判断是否存在对应信息
func ExistMovieDetailByMid(mid int64) bool {
	var count int64
	db.Mdb.Model(&MovieDetail{}).Where("mid", mid).Count(&count)
	return count > 0
}

// SaveSitePlayList 保存附属站点影片信息
func SaveSitePlayList(sl []SlaveMovieInfo) (err error) {
	if len(sl) <= 0 {
		return nil
	}
	if err = db.Mdb.Create(&sl).Error; err != nil {
		log.Println("附属站点影片信息保存失败: ", err)
	}
	return err
}

// UpdateSitePlayList 仅保存播放url列表信息到当前站点
func UpdateSitePlayList(id string, ml []MovieDetail) (err error) {
	// 如果ml 为空则直接返回
	if len(ml) <= 0 {
		return nil
	}
	var sl []SlaveMovieInfo
	for _, m := range ml {
		s := SlaveMovieInfo{Sid: id, Mid: GenerateHashKey(m.Name), DbId: m.DbId, PlayList: m.PlayList}
		// 查询表中是否已经存在对应的数据记录, 如果有则更新, 无则追加到切片中统一处理, id =-1 表示不存在对应数据
		if id := ExistSlaveMovieInfo(s); id > 0 {
			if err = db.Mdb.Model(&s).Where("id", id).Updates(s).Error; err != nil {
				log.Println("附属站点影片信息更新失败: ", err)
			}
			continue
		}
		sl = append(sl, s)
	}
	// 将处理后的结果存储到 SalveMovieInfo表中
	if len(sl) > 0 {
		if err = db.Mdb.Create(&sl).Error; err != nil {
			log.Println("附属站点影片信息保存失败: ", err)
		}
	}
	return
}

// BatchUpdateSlaveInfo 批量更新SlaveMovieInfo
func BatchUpdateSlaveInfo(sl []SlaveMovieInfo) (err error) {
	// 如果ml 为空则直接返回
	if len(sl) <= 0 {
		return nil
	}
	//
	var rl []SlaveMovieInfo
	for _, s := range sl {
		if id := ExistSlaveMovieInfo(s); id > 0 {
			if err = db.Mdb.Model(&s).Where("id", id).Updates(s).Error; err != nil {
				log.Println("附属站点影片信息更新失败: ", err)
			}
			continue
		}
		rl = append(rl, s)
	}
	if len(sl) > 0 {
		if err = db.Mdb.Create(&sl).Error; err != nil {
			log.Println("附属站点影片信息保存失败: ", err)
		}
	}

	return err
}

// DelSlaveInfoBySid 删除sid对应的采集站的所有数据
func DelSlaveInfoBySid(id string) {
	// 查询表中是否存在对应采集站的数据信息
	var count int64
	db.Mdb.Model(&SlaveMovieInfo{}).Count(&count).Where("sid = ?", id)
	// 如果存在对应数据,则进行后续操作
	if count > 0 {
		for {
			res := db.Mdb.Where("sid = ?", id).Limit(5000).Delete(&SlaveMovieInfo{})
			if res.Error != nil {
				log.Println("Delete SlaveMovieInfo Failed: ", res.Error)
				break
			}
			if res.RowsAffected == 0 {
				log.Println("Delete SlaveMovieInfo Over !!!")
				break
			}
			// 短暂休眠, 防止mysql紊乱
			time.Sleep(100 * time.Millisecond)
		}
	}

}

// ExistSlaveMovieInfo 查询对应记录, 如果存在则返还id, 不存在则返还 -1
func ExistSlaveMovieInfo(s SlaveMovieInfo) int64 {
	var id int64
	if err := db.Mdb.Model(&SlaveMovieInfo{}).Select("id").Where("sid = ? AND (mid = ? OR db_id = ?)", s.Sid, s.Mid, s.DbId).First(&id).Error; err != nil {
		// 如果错误类型为gorm.ErrRecordNotFound, 直接返回 0
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0
		}
		// 如果是其他异常则输出异常信息并返回 -1
		log.Println("Find SlaveMovieInfo Failed: ", err)
		return -1
	}
	return id
}

// =================================== Spider数据处理--Redis转存 ========================================================

// MovieDetailCache 主站点数据采集先缓存到redis
func MovieDetailCache(ml []MovieDetail) error {
	// 以mid为key将数据存储到redis的hash中
	var data = make(map[string]string)
	for _, m := range ml {
		r, _ := json.Marshal(m)
		data[strconv.FormatInt(m.Mid, 10)] = string(r)
	}
	return db.Rdb.HSet(db.Cxt, config.MovieDetailKey, data).Err()
}

// SlaveDetailCache 附属站点影片信息缓存
func SlaveDetailCache(id string, ml []MovieDetail) error {
	// 以mid为key将数据存储到redis的hash中
	var data = make(map[string]string)
	for _, m := range ml {
		// 只执行保存操作, 不考虑更新情况
		s := SlaveMovieInfo{Sid: id, Mid: GenerateHashKey(m.Name), DbId: m.DbId, PlayList: m.PlayList}
		r, _ := json.Marshal(s)
		// redis中的存储key优先db_id
		if s.DbId > 0 {
			data[fmt.Sprintf("%d", s.DbId)] = string(r)
		} else {
			data[s.Mid] = string(r)
		}
	}
	// 使用 Sid:Mid为key, 用以区分不同站点数据
	return db.Rdb.HSet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetailKey, id), data).Err()
}

// GetSlaveDetailInCache 从redis缓存中获取播放信息
func GetSlaveDetailInCache(sid, mid string) SlaveMovieInfo {
	// 初始化返回值
	var s SlaveMovieInfo
	v, err := db.Rdb.HGet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetailKey, sid), mid).Result()
	if err != nil {
		// 如果没有获取到对应值, 则直接continue
		//log.Println("Get MultipleSiteDetail Failed: ", err)
		return s
	}
	// 如果获取到数据则直接退出本次循环
	_ = json.Unmarshal([]byte(v), &s)
	return s
}

// SyncMovieDetail 同步redis中的影片数据到mysql中
func SyncMovieDetail(sid string, grade SourceGrade, mode int) {
	// 初始化游标
	var cursor uint64 = 0
	// 根据采集站的类型 Master | Slave 进行不同的处理逻辑
	switch grade {
	case MasterCollect:
		// 循环扫描detail信息, 存储完成后进行删除
		for {
			vs, nextCursor, err := db.Rdb.HScan(db.Cxt, config.MovieDetailKey, cursor, "*", config.FilmScanSize).Result()
			if err != nil {
				log.Println("ScanMovieDetail Failed: ", err)
			}
			if len(vs) > 0 {
				var ks []string
				var ml []MovieDetail
				for i := 0; i < len(vs); i += 2 {
					ks = append(ks, vs[i])
					var m MovieDetail
					_ = json.Unmarshal([]byte(vs[i+1]), &m)
					ml = append(ml, m)
				}
				// 批量保存movieDetail
				switch mode {
				case 0:
					// 执行全量保存
					if err := SaveDetails(ml); err != nil {
						log.Println("SyncMovieDetail AllSave Failed: ", err)
					}
				case 1:
					// 执行更新
					if err := BatchUpdateDetails(ml); err != nil {
						log.Println("SyncMovieDetail SaveOrUpdate Failed: ", err)
					}
				default:
					log.Println("Synchronization Mode Exception:", mode)
				}

				// 删除已提取的元素
				if err := db.Rdb.HDel(db.Cxt, config.MovieDetailKey, ks...).Err(); err != nil {
					log.Println("DeleteMovieDetailCache Failed: ", err)
				}
			}
			// 更新游标
			cursor = nextCursor
			// 如果游标归零则结束循环同步
			if cursor <= 0 {
				break
			}
		}
	case SlaveCollect:
		// 循环扫描detail信息, 存储完成后进行删除
		for {
			vs, nextCursor, err := db.Rdb.HScan(db.Cxt, fmt.Sprintf(config.MultipleSiteDetailKey, sid), cursor, "", config.FilmScanSize).Result()
			if err != nil {
				log.Println("ScanSlaveDetail Failed: ", err)
			}
			if len(vs) > 0 {
				var ks []string
				var sl []SlaveMovieInfo
				for i := 0; i < len(vs); i += 2 {
					ks = append(ks, vs[i])
					var s SlaveMovieInfo
					_ = json.Unmarshal([]byte(vs[i+1]), &s)
					sl = append(sl, s)
				}
				// 批量保存movieDetail
				switch mode {
				case 0:
					// 执行全量保存
					if err := SaveSitePlayList(sl); err != nil {
						log.Println("SyncSlaveDetail AllSave Failed: ", err)
					}
				case 1:
					// 执行更新
					if err := BatchUpdateSlaveInfo(sl); err != nil {
						log.Println("SyncSlaveDetail SaveOrUpdate Failed: ", err)
					}
				default:
					log.Println("Synchronization Mode Exception:", mode)
				}
				// 删除已提取的元素
				if err := db.Rdb.HDel(db.Cxt, fmt.Sprintf(config.MultipleSiteDetailKey, sid), ks...).Err(); err != nil {
					log.Println("DeleteSlaveDetailCache Failed: ", err)
				}
			}
			// 更新游标
			cursor = nextCursor
			// 如果游标归零则结束循环同步
			if cursor == 0 {
				break
			}
		}
	}

}

// ============================ APi接口 ==================================================

// GetDetailByMid 获取影片对应的详情信息
func GetDetailByMid(mid int64) MovieDetail {
	// 初始化返回值
	var m MovieDetail
	// 从redis获取对应的影片信息
	v, err := db.Rdb.HGet(db.Cxt, config.MovieDetailKey, strconv.FormatInt(mid, 10)).Result()
	if err != nil {
		// 如果没有获取到对应值, 则去mysql中进行查找
		if errors.Is(err, redis.Nil) {
			if err := db.Mdb.Where("mid = ?", mid).Find(&m).Error; err != nil {
				log.Println("Find BasicInfo Failed: ", err)
				return m
			}
			// 执行本地图片匹配
			ReplaceDetailPic(&m)
			return m
		}
		log.Println("Find MovieDetail Failed: ", err)
		return m
	}
	// 如果获取到对应值,则进行反序列化
	_ = json.Unmarshal([]byte(v), &m)
	return m
	//var d MovieDetail
	//// 查询mid对应的影片详情信息, 只查询部分字段
	//if err := db.Mdb.Model(&MovieDetail{}).Where("mid = ?", mid).First(&d).Error; err != nil {
	//	log.Println("Find MovieDetail Failed: ", err)
	//	return d
	//}
	//// 执行本地图片匹配
	//ReplaceDetailPic(&d)
	//return d
}

// GetBasicInfoByMid 获取Id对应的影片基本信息
func GetBasicInfoByMid(mid int64) MovieBasicInfo {
	// 通过id查询满足条件的影片基本信息
	var basic MovieBasicInfo
	var d MovieDetail
	// 查询mid对应的影片详情信息, 只查询部分字段
	if err := db.Mdb.Model(&MovieDetail{}).Select("id, mid, cid, pid, name, sub_title, c_name, state, picture, actor, director,"+
		" content, remarks, area, year").Where("mid = ?", mid).First(&d).Error; err != nil {
		log.Println("Find MovieDetail Failed: ", err)
		return basic
	}
	// 匹配本地图片
	ReplaceDetailPic(&d)
	// 将 MovieDetail转化为 BasicInfo
	basic = ConvertBasicInfo(d)
	return basic
}

// GetBasicInfoByIds 通过searchInfo 获取影片的基本信息
func GetBasicInfoByIds(ids []int64) []MovieBasicInfo {
	// 初始化返回值
	var l []MovieBasicInfo
	// 首先从redis中获取影片的最新信息, 如果没有则转为去mysql表中获取
	var ks []string
	for _, id := range ids {
		ks = append(ks, strconv.FormatInt(id, 10))
	}
	// 一次性获取所有
	vs, err := db.Rdb.HMGet(db.Cxt, config.MovieDetailKey, ks...).Result()
	if err != nil {
		log.Println("Find MovieDetail Failed: ", err)
		return l
	}
	// 迭代转换 basicInfo, 并将未获取到值的id进行整合
	var newIds []int64
	var ml []MovieDetail
	if len(vs) > 0 {
		for i, v := range vs {
			if v != nil {
				var m MovieDetail
				_ = json.Unmarshal([]byte(v.(string)), &m)
				ReplaceDetailPic(&m)
				l = append(l, ConvertBasicInfo(m))
			} else {
				newIds = append(newIds, ids[i])
			}
		}
	}
	// 如果存在nil值,则去mysql进行补全
	if len(newIds) > 0 {
		if err := db.Mdb.Model(&MovieDetail{}).Select("id, mid, cid, pid, name, sub_title, c_name, state, picture, actor, director,"+
			" content, remarks, area, year").Where("mid IN (?)", newIds).Find(&ml).Error; err != nil {
			log.Println("BatchFind BasicInfo Failed: ", err)
			return nil
		}
		for _, m := range ml {
			// 执行本地图片匹配
			ReplaceDetailPic(&m)
			l = append(l, ConvertBasicInfo(m))
		}
	}

	//var ml []MovieDetail
	//var l []MovieBasicInfo
	// 使用in查询, 一次性拿到满足条件的数据
	//if err := db.Mdb.Model(&MovieDetail{}).Select("id, mid, cid, pid, name, sub_title, c_name, state, picture, actor, director,"+
	//	" content, remarks, area, year").Where("mid IN (?)", ids).Find(&ml).Error; err != nil {
	//	log.Println("BatchFind BasicInfo Failed: ", err)
	//	return nil
	//}
	//// 将查询到的结果批量转化为BasicInfo
	//for _, m := range ml {
	//	// 执行本地图片匹配
	//	ReplaceDetailPic(&m)
	//	l = append(l, ConvertBasicInfo(m))
	//}
	return l
}

// GetMovieListByPid  通过Pid 分类ID 获取对应影片的数据信息
func GetMovieListByPid(pid int64, page *Page) []MovieBasicInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("pid", pid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 通过Search表查询
	var ids []int64
	if err := db.Mdb.Model(&SearchInfo{}).Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Select("mid").Where("pid", pid).Order("update_stamp DESC").Find(&ids).Error; err != nil {
		log.Println(err)
		return nil
	}
	// 通过ids查询影片基本信息并返回
	return GetBasicInfoByIds(ids)
}

// GetMovieListByCid 通过Cid查找对应的影片分页数据, 不适合GetMovieListByPid 糅合
func GetMovieListByCid(cid int64, page *Page) []MovieBasicInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("cid", cid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 进行具体的信息查询
	var ids []int64
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Select("mid").Where("cid", cid).Order("update_stamp DESC").Find(&ids).Error; err != nil {
		log.Println(err)
		return nil
	}
	// 通过影片ID去redis中获取id对应数据信息
	return GetBasicInfoByIds(ids)
}

// GetRelateMovieBasicInfo GetRelateMovie 根据 name, cid, pid, classTag 获取相关影片
func GetRelateMovieBasicInfo(search SearchInfo, page *Page) []MovieBasicInfo {
	/*
		根据当前影片信息匹配相关的影片
		1. 分类Cid,
		2. 如果影片名称含有第x季 则根据影片名进行模糊匹配
		3. class_tag 剧情内容匹配, 切分后使用 or 进行匹配
		4. area 地区
		5. 语言 Language
	*/
	// sql 拼接查询条件
	sql := ""

	// 优先进行名称相似匹配, 先对影片名称进行精简, 只保留主体用于匹配同系列影片
	name := util.CleanFilmName(search.Name)
	//sql = fmt.Sprintf(`select mid from %s where (name LIKE "%%%s%%" or sub_title LIKE "%%%[2]s%%") AND cid=%d AND search.deleted_at IS NULL union`, search.TableName(), name, search.Cid)
	sql = fmt.Sprintf(`select mid from %s where MATCH(name, sub_title) AGAINST('%s') AND cid=%d AND search.deleted_at IS NULL union`, search.TableName(), name, search.Cid)

	// 添加其他相似匹配规则 同属二级分类
	sql = fmt.Sprintf(`%s (select mid from %s where cid=%d AND`, sql, search.TableName(), search.Cid)
	// 根据剧情标签查找相似影片, classTag 使用的分隔符为 , | /首先去除 classTag 中包含的所有空格
	search.ClassTag = strings.ReplaceAll(search.ClassTag, " ", "")
	// 如果 classTag 中包含分割符则进行拆分匹配
	//cl := strings.Split(util.FormatSpecialChar(search.ClassTag), ",")
	// 将
	search.ClassTag = strings.ReplaceAll(util.FormatSpecialChar(search.ClassTag), ",", " ")
	//if len(cl) > 0 {
	//	s := "("
	//	for _, c := range cl {
	//		//s = fmt.Sprintf(`%s class_tag like "%%%s%%" OR`, s, c)
	//		s = fmt.Sprintf(`%s class_tag like "%%%s%%" OR`, s, c)
	//	}
	//	sql = fmt.Sprintf("%s %s)", sql, strings.TrimSuffix(s, "OR"))
	//	sql = fmt.Sprintf("%s %s)", sql, strings.TrimSuffix(s, "OR"))
	//} else {
	//	sql = fmt.Sprintf(`%s class_tag like "%%%s%%"`, sql, search.ClassTag)
	//}
	sql = fmt.Sprintf(`%s MATCH(class_tag) AGAINST('%s')`, sql, search.ClassTag)
	// 除名称外的相似影片使用随机排序
	//sql = fmt.Sprintf("%s ORDER BY RAND() limit %d,%d)", sql, page.Current, page.PageSize)
	sql = fmt.Sprintf("%s AND search.deleted_at IS NULL limit %d,%d)", sql, page.Current, page.PageSize)
	// 条件拼接完成后加上limit参数
	sql = fmt.Sprintf("(%s)  limit %d,%d", sql, page.Current, page.PageSize)
	// 执行sql, 获取满足条件的影片mid切片
	var ids []int64
	db.Mdb.Raw(sql).Scan(&ids)
	// 通过 ids 获取影片基本信息,并返回
	return GetBasicInfoByIds(ids)
}

// GetMultiplePlay 通过影片名的ID值匹配播放源, 不区分站点
func GetMultiplePlay(mIds []string, dbId int64) []SlaveMovieInfo {
	// 初始化返回值
	var l []SlaveMovieInfo
	// 首先从redis进行匹配
	for _, c := range GetCollectSourceListByGrade(SlaveCollect) {
		if !c.State {
			continue
		}
		var s SlaveMovieInfo
		// 优先使用dbID为key去redis中获取
		if s = GetSlaveDetailInCache(c.Id, fmt.Sprintf("%d", dbId)); s.Mid != "" {
			l = append(l, s)
			continue
		}
		// 如果匹配失败则使用name生成的mIds获取数据
		for _, mid := range mIds {
			// 初始化临时变量 SlaveMovieInfo
			if s = GetSlaveDetailInCache(c.Id, mid); s.Mid != "" {
				l = append(l, s)
				break
			}
			//v, err := db.Rdb.HGet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetailKey, c.Id), mid).Result()
			//if err != nil {
			//	// 如果没有获取到对应值, 则直接continue
			//	continue
			//}
			//// 如果获取到数据则直接退出本次循环
			//_ = json.Unmarshal([]byte(v), &s)
			//l = append(l, s)
			//break
		}

		// Redis中没有匹配到对应数据, 则去slave_info表中获取数据
		//如果 dbID 不为0 则优先使用sid 和 dbId 去mysql中锁定对应数据
		if s.Mid == "" {
			//if err := db.Mdb.Select("sid, play_list").Where("sid = ? AND db_id = ?", c.Id, dbId).Last(&s).Error; err != nil {
			//	log.Println("GetMultiplePlay Failed: ", err)
			//} else {
			//	continue
			//}
			if err := db.Mdb.Select("sid, play_list").Where("sid = ? AND db_id = ?", c.Id, dbId).Last(&s).Error; err == nil {
				l = append(l, s)
				continue
			}
		}
		// 如果db_id依旧获取失败, 则使用mIds进行最后的获取
		if s.Mid == "" {
			//if err := db.Mdb.Raw("(SELECT sid, play_list FROM `slave_infos` WHERE sid = ? AND mid IN (?)) UNION ALL (SELECT sid, play_list FROM `slave_infos` WHERE sid = ? AND db_id = ? ORDER BY `slave_infos`.`id` LIMIT 1) LIMIT 1", c.Id, mIds, c.Id, dbId).First(&s).Error; err != nil {
			if err := db.Mdb.Select("sid, play_list").Where("sid = ? AND mid IN (?)", c.Id, mIds).Last(&s).Error; err != nil {
				//log.Println("GetMultiplePlay Failed: ", err)
				continue
			}
			l = append(l, s)
		}
	}
	// 通过siteId, mIds, dbIds 检索满足条件的数据
	//if err := db.Mdb.Model(&SlaveMovieInfo{}).Select("sid, play_list").Where("mid IN (?) OR db_id = ?", mIds, dbId).Find(&l).Error; err != nil {
	//	log.Println("GetMultiplePlay Failed: ", err)
	//	return nil
	//}

	return l
}

// ============================ 数据处理 ==================================================

// ConvertSearchInfo 将detail信息处理成 searchInfo
func ConvertSearchInfo(m MovieDetail) SearchInfo {
	score, _ := strconv.ParseFloat(m.DbScore, 64)
	stamp, _ := time.ParseInLocation(time.DateTime, m.UpdateTime, time.Local)
	// detail中的年份信息并不准确, 因此采用 ReleaseDate中的年份
	year, err := strconv.ParseInt(regexp.MustCompile(`[1-9][0-9]{3}`).FindString(m.ReleaseDate), 10, 64)
	if err != nil {
		year = 0
	}
	return SearchInfo{
		Mid:         m.Mid,
		Cid:         m.Cid,
		Pid:         m.Pid,
		Name:        m.Name,
		SubTitle:    m.SubTitle,
		CName:       m.CName,
		ClassTag:    m.ClassTag,
		Area:        m.Area,
		Language:    m.Language,
		Year:        year,
		Initial:     m.Initial,
		Score:       score,
		Hits:        m.Hits,
		UpdateStamp: stamp.Unix(),
		State:       m.State,
		Remarks:     m.Remarks,
		// ReleaseDate 部分影片缺失该参数, 所以使用添加时间作为上映时间排序
		ReleaseStamp: m.AddTime,
	}
}

// ConvertBasicInfo 将Detail信息转化为basic信息
func ConvertBasicInfo(m MovieDetail) MovieBasicInfo {
	return MovieBasicInfo{Id: m.Mid, Cid: m.Cid, Pid: m.Pid, Name: m.Name, SubTitle: m.SubTitle,
		CName: m.CName, State: m.State, Picture: m.Picture, Actor: m.Actor, Director: m.Director, Blurb: m.Content,
		Remarks: m.Remarks, Area: m.Area, Year: m.Year}
}

/*
	对附属播放源入库时的name|dbID进行处理,保证唯一性
1. 去除name中的所有空格
2. 去除name中含有的别名～.*～
3. 去除name首尾的标点符号
4. 将处理完成后的name转化为hash值作为存储时的key
*/
// GenerateHashKey 存储播放源信息时对影片名称进行处理-生成id, 提高各站点间同一影片的匹配度
func GenerateHashKey[K string | ~int | int64](key K) string {
	mName := fmt.Sprint(key)
	//1. 去除name中的所有空格
	mName = regexp.MustCompile(`\s`).ReplaceAllString(mName, "")
	//2. 添加常用的名称标准化替换规则
	rules := []string{
		// 中文季数标签统一
		"season", "s", "第", "s", "季", "", "期", "", "画", "",
		// --- 3. 剧场版标准化 ---
		"剧场版", "ovo", "映画", "ovo", "电影版", "ovo", "The Movie", "ovo", "Movie", "ovo", "(Movie)", "ovo", "〔映画〕", "ovo",
		// 特殊数学符号 (用户常用来代替数字，如 ∬ 代表 2)
		"Ⅰ", "1", "Ⅱ", "2", "Ⅲ", "3",
		"∫", "1", "∬", "2", "∮", "3", "Ⅳ", "4", "Ⅴ", "5", "Ⅵ", "6", "Ⅶ", "7", "Ⅷ", "8", "Ⅸ", "9", "Ⅹ", "10", // 用户可能用积分号代表季数
		"一", "1", "二", "2", "三", "3", "四", "4", "五", "5", "六", "6", "七", "7", "八", "8", "九", "9",
		// 移除或替换无意义的装饰符号，这些符号在搜索中通常不仅无用还会阻碍匹配
		"★", "", "☆", "", "◆", "", "◇", "", "●", "", "○", "",
		"【", "", "】", "", "〖", "", "〗", "", "〔", "", "〕", "",
		"「", "", "」", "", "『", "", "』", "",
		"|", "", "｜", "", // 竖线分隔符
		"~", "", "～", "", // 波浪号
		"...", "", "……", "", // 省略号
		"!", "", "！", "", "?", "", "？", "",
		"(", "", ")", "", "（", "", "）", "",
		"[", "", "]", "", "［", "", "］", "",
		"{", "", "}", "", "｛", "", "｝", "",
		"＆", "&", "＋", "+",
		"-", "", "－", "", "—", "", "–", "", // 策略：通常移除所有标点，让 "A-B" 变成 "AB"
		"_", "", "＿", "",
		".", "", "．", "", "。", "",
		",", "", "，", "",
		":", "", ":", "", ":", "",
		";", "", "；", "",
		"'", "", "’", "", "\"", "", "“", "", "”", "",
		"`", "", "｀", "",
	}
	mName = strings.NewReplacer(rules...).Replace(mName)
	//3. 去除name首尾的标点符号
	mName = regexp.MustCompile(`^[[:punct:]]+|[[:punct:]]+$`).ReplaceAllString(mName, "")
	//4. 将处理完成后的name转化为hash值作为存储时的key
	h := fnv.New32a()
	_, err := h.Write([]byte(mName))
	if err != nil {
		return ""
	}
	return fmt.Sprint(h.Sum32())
}
