package model

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"reflect"
	"regexp"
	"server/config"
	"server/plugin/db"
	"strings"
)

// SearchInfo 存储用于检索的信息
type SearchInfo struct {
	gorm.Model
	Mid         int64   `json:"mid" gorm:"uniqueIndex:idx_mid"` //影片ID
	Cid         int64   `json:"cid"`                            //分类ID
	Pid         int64   `json:"pid"`                            //上级分类ID
	Name        string  `json:"name"`                           // 片名
	SubTitle    string  `json:"subTitle"`                       // 影片子标题
	CName       string  `json:"CName"`                          // 分类名称
	ClassTag    string  `json:"classTag"`                       //类型标签
	Area        string  `json:"area"`                           // 地区
	Language    string  `json:"language"`                       // 语言
	Year        int64   `json:"year"`                           // 年份
	Initial     string  `json:"initial"`                        // 首字母
	Score       float64 `json:"score"`                          //评分
	Time        int64   `json:"time"`                           // 更新时间
	Rank        int64   `json:"rank"`                           // 热度排行id
	State       string  `json:"state"`                          //状态 正片|预告
	Remarks     string  `json:"remarks"`                        // 完结 | 更新至x集
	ReleaseDate int64   `json:"releaseDate"`                    //上映时间 时间戳
}

// Page 分页信息结构体
type Page struct {
	PageSize  int `json:"pageSize"`  // 每页大小
	Current   int `json:"current"`   // 当前页
	PageCount int `json:"pageCount"` // 总页数
	Total     int `json:"total"`     // 总记录数
	//List      []interface{} `json:"list"`      // 数据
}

func (s *SearchInfo) TableName() string {
	return config.SearchTableName
}

// ================================= Spider 数据处理(redis) =================================

// RdbSaveSearchInfo 批量保存检索信息到redis
func RdbSaveSearchInfo(list []SearchInfo) {
	// 1.整合一下zset数据集
	var members []redis.Z
	for _, s := range list {
		member, _ := json.Marshal(s)
		members = append(members, redis.Z{Score: float64(s.Mid), Member: member})
	}
	// 2.批量保存到zset集合中
	db.Rdb.ZAdd(db.Cxt, config.SearchInfoTemp, members...)
}

// ScanSearchInfo 批量扫描处理详情检索信息, 返回检索信息列表和下次开始的游标
func ScanSearchInfo(cursor uint64, count int64) ([]SearchInfo, uint64) {
	// 1.从redis中批量扫描详情信息
	list, nextCursor := db.Rdb.ZScan(db.Cxt, config.SearchInfoTemp, cursor, "*", count).Val()
	// 2. 处理数据
	var resList []SearchInfo
	for i, s := range list {
		// 3. 判断当前是否是元素
		if i%2 == 0 {
			info := SearchInfo{}
			_ = json.Unmarshal([]byte(s), &info)
			info.Model = gorm.Model{}
			resList = append(resList, info)
		}
	}
	return resList, nextCursor
}

// RemoveAll 删除所有库存数据
func RemoveAll() {
	// 删除redis中当前库存储的所有数据
	db.Rdb.FlushDB(db.Cxt)
	// 删除mysql中留存的检索表
	var s *SearchInfo
	db.Mdb.Exec(fmt.Sprintf(`drop table if exists %s`, s.TableName()))
}

// ================================= Spider 数据处理(mysql) =================================

// CreateSearchTable 创建存储检索信息的数据表
func CreateSearchTable() {
	// 1. 判断表中是否存在当前表
	isExist := db.Mdb.Migrator().HasTable(&SearchInfo{})
	// 如果不存在则创建表
	if !isExist {
		err := db.Mdb.AutoMigrate(&SearchInfo{})
		if err != nil {
			log.Println("Create Table SearchInfo Failed: ", err)
		}
	}
}

// BatchSave 批量保存影片search信息
func BatchSave(list []SearchInfo) {
	tx := db.Mdb.Begin()
	// 防止程序异常终止
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.CreateInBatches(list, len(list)).Error; err != nil {
		// 插入失败则回滚事务, 重新进行插入
		tx.Rollback()
		return
	}
	// 插入成功后输出一下成功信息
	//log.Println("BatchSave SearchInfo Successful, Count: ", len(list))
	tx.Commit()
}

// BatchSaveOrUpdate 判断数据库中是否存在对应mid的数据, 如果存在则更新, 否则插入
func BatchSaveOrUpdate(list []SearchInfo) {
	tx := db.Mdb.Begin()
	// 失败则回滚事务
	//defer func() {
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//	}
	//}()
	for _, info := range list {
		var count int64
		// 通过当前影片id 对应的记录数
		tx.Model(&SearchInfo{}).Where("mid", info.Mid).Count(&count)
		// 如果存在对应数据则进行更新, 否则进行删除
		if count > 0 {
			// 记录已经存在则执行更新部分内容
			err := tx.Model(&SearchInfo{}).Where("mid", info.Mid).Updates(SearchInfo{Time: info.Time, Rank: info.Rank, State: info.State,
				Remarks: info.Remarks, Score: info.Score, ReleaseDate: info.ReleaseDate}).Error
			if err != nil {
				tx.Rollback()
			}
		} else {
			// 执行插入操作
			if err := tx.Create(&info).Error; err != nil {
				tx.Rollback()
			}
		}
	}
	// 提交事务
	tx.Commit()
}

// SaveSearchData 添加影片检索信息
func SaveSearchData(s SearchInfo) {
	// 先查询数据库中是否存在对应记录
	isExist := SearchMovieInfo(s.Mid)
	// 如果不存在对应记录则 保存当前记录
	if !isExist {
		db.Mdb.Create(&s)
	}
}

// SearchMovieInfo 通过Mid查询符合条件的数据
func SearchMovieInfo(mid int64) bool {
	search := SearchInfo{}
	db.Mdb.Where("mid", mid).First(&search)
	// reflect.DeepEqual(a, A{})
	return !reflect.DeepEqual(search, SearchInfo{})
}

// TunCateSearchTable 截断SearchInfo数据表
func TunCateSearchTable() {
	var searchInfo *SearchInfo
	err := db.Mdb.Exec(fmt.Sprint("TRUNCATE TABLE ", searchInfo.TableName())).Error
	if err != nil {
		log.Println("TRUNCATE TABLE Error: ", err)
	}
}

// ================================= API 数据接口信息处理 =================================

// GetMovieListByPid  通过Pid 分类ID 获取对应影片的数据信息
func GetMovieListByPid(pid int64, page *Page) []MovieBasicInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("pid", pid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 进行具体的信息查询
	var s []SearchInfo
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Where("pid", pid).Order("year DESC, time DESC").Find(&s).Error; err != nil {
		log.Println(err)
		return nil
	}
	// 通过影片ID去redis中获取id对应数据信息
	var list []MovieBasicInfo
	for _, v := range s {
		// 通过key搜索指定的影片信息 , MovieDetail:Cid6:Id15441
		list = append(list, GetBasicInfoByKey(fmt.Sprintf(config.MovieBasicInfoKey, v.Cid, v.Mid)))
	}
	return list
}

// SearchFilmKeyword 通过关键字搜索库存中满足条件的影片名
func SearchFilmKeyword(keyword string, page *Page) []SearchInfo {
	var searchList []SearchInfo
	// 1. 先统计搜索满足条件的数据量
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("name LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 2. 获取满足条件的数据
	db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).
		Where("name LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Order("year DESC, time DESC").Find(&searchList)
	return searchList
}

// GetMovieListByCid 通过Cid查找对应的影片分页数据, 不适合GetMovieListByPid 糅合
func GetMovieListByCid(cid int64, page *Page) []MovieBasicInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("cid", cid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 进行具体的信息查询
	var s []SearchInfo
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Where("cid", cid).Order("year DESC, time DESC").Find(&s).Error; err != nil {
		log.Println(err)
		return nil
	}
	// 通过影片ID去redis中获取id对应数据信息
	var list []MovieBasicInfo
	for _, v := range s {
		// 通过key搜索指定的影片信息 , MovieDetail:Cid6:Id15441
		list = append(list, GetBasicInfoByKey(fmt.Sprintf(config.MovieBasicInfoKey, v.Cid, v.Mid)))
	}
	return list
}

// GetRelateMovieBasicInfo GetRelateMovie 根据SearchInfo获取相关影片
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
	// 优先进行名称相似匹配
	re := regexp.MustCompile("第.{1,3}季")
	if re.MatchString(search.Name) {
		search.Name = re.ReplaceAllString(search.Name, "")
		sql = fmt.Sprintf(`select * from %s where name LIKE "%%%s%%" union`, search.TableName(), search.Name)
	}
	// 执行后续匹配内容
	//sql = fmt.Sprintf(`%s select * from %s where cid=%d AND area="%s" AND language="%s" AND`, sql, search.TableName(), search.Cid, search.Area, search.Language)

	// 地区限制取消, 过滤掉的影片太多
	sql = fmt.Sprintf(`%s select * from %s where cid=%d AND language="%s" AND`, sql, search.TableName(), search.Cid, search.Language)
	if strings.Contains(search.ClassTag, ",") {
		s := "("
		for _, t := range strings.Split(search.ClassTag, ",") {
			s = fmt.Sprintf(`%s class_tag = "%s" OR`, s, t)
		}
		sql = fmt.Sprintf("%s %s)", sql, strings.TrimSuffix(s, "OR"))
	} else {
		sql = fmt.Sprintf(`%s class_tag = "%s"`, sql, search.ClassTag)
	}
	// 条件拼接完成后加上limit参数
	sql = fmt.Sprintf("(%s)  limit %d,%d", sql, page.Current, page.PageSize)
	// 执行sql
	list := []SearchInfo{}
	db.Mdb.Raw(sql).Scan(&list)
	// 根据list 获取对应的BasicInfo
	var basicList []MovieBasicInfo
	for _, s := range list {
		// 通过key获取对应的影片基本数据
		basicList = append(basicList, GetBasicInfoByKey(fmt.Sprintf(config.MovieBasicInfoKey, s.Cid, s.Mid)))
	}

	return basicList
}

// GetMultiplePlay 通过影片名hash值匹配播放源
func GetMultiplePlay(siteName, key string) []MovieUrlInfo {
	data := db.Rdb.HGet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetail, siteName), key).Val()
	var playList []MovieUrlInfo
	_ = json.Unmarshal([]byte(data), &playList)
	return playList
}
