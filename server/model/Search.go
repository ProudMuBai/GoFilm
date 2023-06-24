package model

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"math"
	"reflect"
	"regexp"
	"server/config"
	"server/plugin/db"
	"strings"
	"time"
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
	UpdateStamp int64   `json:"updateStamp"`                    // 更新时间
	Hits        int64   `json:"hits"`                           // 热度排行
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

// DelMtPlay 清空附加播放源信息
func DelMtPlay(keys []string) {
	db.Rdb.Del(db.Cxt, keys...)
}

/*
SearchKeyword 设置search关键字集合(影片分类检索类型数据)
	类型, 剧情 , 地区, 语言, 年份, 首字母, 排序
	1. 在影片详情缓存到redis时将影片的相关数据进行记录, 存在相同类型则分值加一
	2. 通过分值对类型进行排序类型展示到页面
*/

func SaveSearchTag(search SearchInfo) {
	// 声明用于存储采集的影片的分类检索信息
	//searchMap := make(map[string][]map[string]int)

	// Redis中的记录形式 Search:SearchKeys:Pid1:Title Hash
	// Redis中的记录形式 Search:SearchKeys:Pid1:xxx Hash

	// 获取redis中的searchMap
	key := fmt.Sprintf(config.SearchTitle, search.Pid)
	searchMap := db.Rdb.HGetAll(db.Cxt, key).Val()
	// 是否存储对应分类的map, 如果不存在则缓存一份
	if len(searchMap) == 0 {
		searchMap["Category"] = "类型"
		searchMap["Plot"] = "剧情"
		searchMap["Area"] = "地区"
		searchMap["Language"] = "语言"
		searchMap["Year"] = "年份"
		searchMap["Initial"] = "首字母"
		searchMap["Sort"] = "排序"
		db.Rdb.HMSet(db.Cxt, key, searchMap)
	}
	// 对searchMap中的各个类型进行处理
	for k, _ := range searchMap {
		tagKey := fmt.Sprintf(config.SearchTag, search.Pid, k)
		tagCount := db.Rdb.ZCard(db.Cxt, tagKey).Val()
		switch k {
		case "Category":
			// 获取 Category 数据, 如果不存在则缓存一份
			if tagCount == 0 {
				var tags []redis.Z
				for _, t := range GetChildrenTree(search.Pid) {
					tags = append(tags, redis.Z{Score: float64(-t.Id), Member: t.Name})
				}
				db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k), tags...)
			}
		case "Year":
			// 获取 Year 数据, 如果不存在则缓存一份
			if tagCount == 0 {
				var tags []redis.Z
				currentYear := time.Now().Year()
				for i := 0; i < 12; i++ {
					tags = append(tags, redis.Z{Score: float64(currentYear - i), Member: currentYear - i})
				}
				db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k), tags...)
			}
		case "Initial":
			// 如果不存在 首字母 Tag 数据, 则缓存一份
			if tagCount == 0 {
				var tags []redis.Z
				for i := 65; i <= 90; i++ {
					tags = append(tags, redis.Z{Score: float64(90 - i), Member: fmt.Sprintf("%c", i)})
				}
				db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k), tags...)
			}
		case "Sort":
			if tagCount == 0 {
				tags := []redis.Z{
					{2, "time"},
					{1, "hits"},
					{0, "score"},
				}
				db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k), tags...)
			}
		case "Plot":
			HandleSearchTags(search.ClassTag, tagKey)
		case "Area":
			HandleSearchTags(search.Area, tagKey)
		case "Language":
			HandleSearchTags(search.Language, tagKey)
		default:
			break
		}
	}

}

func HandleSearchTags(preTags string, k string) {
	// 先处理字符串中的空白符 然后对处理前的tag字符串进行分割
	preTags = regexp.MustCompile(`[\s\n\r]+`).ReplaceAllString(preTags, "")
	f := func(sep string) {
		for _, t := range strings.Split(preTags, sep) {
			// 获取 tag对应的score
			score := db.Rdb.ZScore(db.Cxt, k, t).Val()
			// 在原score的基础上+1 重新存入redis中

			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: score + 1, Member: t})
		}
	}
	switch {
	case strings.Contains(preTags, "/"):
		f("/")
	case strings.Contains(preTags, ","):
		f(",")
	case strings.Contains(preTags, "，"):
		f("，")
	case strings.Contains(preTags, "、"):
		f("、")
	default:
		// 获取 tag对应的score
		if len(preTags) == 0 || preTags == "其它" {
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: 0, Member: preTags})
		} else {
			score := db.Rdb.ZScore(db.Cxt, k, preTags).Val()
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: score + 1, Member: preTags})
		}
	}

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
	for _, info := range list {
		var count int64
		// 通过当前影片id 对应的记录数
		tx.Model(&SearchInfo{}).Where("mid", info.Mid).Count(&count)
		// 如果存在对应数据则进行更新, 否则进行删除
		if count > 0 {
			// 记录已经存在则执行更新部分内容
			err := tx.Model(&SearchInfo{}).Where("mid", info.Mid).Updates(SearchInfo{UpdateStamp: info.UpdateStamp, Hits: info.Hits, State: info.State,
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
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Where("pid", pid).Order("year DESC, update_stamp DESC").Find(&s).Error; err != nil {
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

// GetMovieListByCid 通过Cid查找对应的影片分页数据, 不适合GetMovieListByPid 糅合
func GetMovieListByCid(cid int64, page *Page) []MovieBasicInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("cid", cid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 进行具体的信息查询
	var s []SearchInfo
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Where("cid", cid).Order("year DESC, update_stamp DESC").Find(&s).Error; err != nil {
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

// GetHotMovieByPid  获取指定类别的热门影片
func GetHotMovieByPid(pid int64, page *Page) []SearchInfo {
	// 返回分页参数
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("pid", pid).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 进行具体的信息查询
	var s []SearchInfo
	// 当前时间偏移一个月
	t := time.Now().AddDate(0, -1, 0).Unix()
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).Where("pid=? AND update_stamp > ?", pid, t).Order(" year DESC, hits DESC").Find(&s).Error; err != nil {
		log.Println(err)
		return nil
	}
	return s
}

// SearchFilmKeyword 通过关键字搜索库存中满足条件的影片名
func SearchFilmKeyword(keyword string, page *Page) []SearchInfo {
	var searchList []SearchInfo
	// 1. 先统计搜索满足条件的数据量
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("name LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Or("sub_title LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 2. 获取满足条件的数据
	db.Mdb.Limit(page.PageSize).Offset((page.Current-1)*page.PageSize).
		Where("name LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Or("sub_title LIKE ?", fmt.Sprint(`%`, keyword, `%`)).Order("year DESC, update_stamp DESC").Find(&searchList)
	return searchList
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
	//search.Name = regexp.MustCompile("第.{1,3}季").ReplaceAllString(search.Name, "")
	name := regexp.MustCompile(`(第.{1,3}季.*)|([0-9]{1,3})|(剧场版)|(\s\S*$)|(之.*)|([\p{P}\p{S}].*)`).ReplaceAllString(search.Name, "")
	// 如果处理后的影片名称依旧没有改变 且具有一定长度 则截取部分内容作为搜索条件
	if len(name) == len(search.Name) && len(name) > 10 {
		// 中文字符需截取3的倍数,否则可能乱码
		name = name[:int(math.Ceil(float64(len(name)/5))*3)]
	}
	sql = fmt.Sprintf(`select * from %s where (name LIKE "%%%s%%" or sub_title LIKE "%%%[2]s%%") AND cid=%d union`, search.TableName(), name, search.Cid)
	// 执行后续匹配内容, 匹配结果过少,减少过滤条件
	//sql = fmt.Sprintf(`%s select * from %s where cid=%d AND area="%s" AND language="%s" AND`, sql, search.TableName(), search.Cid, search.Area, search.Language)

	// 添加其他相似匹配规则
	sql = fmt.Sprintf(`%s (select * from %s where cid=%d AND `, sql, search.TableName(), search.Cid)
	// 根据剧情标签查找相似影片, classTag 使用的分隔符为 , | /
	// 首先去除 classTag 中包含的所有空格
	search.ClassTag = strings.ReplaceAll(search.ClassTag, " ", "")
	// 如果 classTag 中包含分割符则进行拆分匹配
	if strings.Contains(search.ClassTag, ",") {
		s := "("
		for _, t := range strings.Split(search.ClassTag, ",") {
			s = fmt.Sprintf(`%s class_tag like "%%%s%%" OR`, s, t)
		}
		sql = fmt.Sprintf("%s %s)", sql, strings.TrimSuffix(s, "OR"))
	} else if strings.Contains(search.ClassTag, "/") {
		s := "("
		for _, t := range strings.Split(search.ClassTag, "/") {
			s = fmt.Sprintf(`%s class_tag like "%%%s%%" OR`, s, t)
		}
		sql = fmt.Sprintf("%s %s)", sql, strings.TrimSuffix(s, "OR"))
	} else {
		sql = fmt.Sprintf(`%s class_tag like "%%%s%%"`, sql, search.ClassTag)
	}
	// 除名称外的相似影片使用随机排序
	sql = fmt.Sprintf("%s ORDER BY RAND() limit %d,%d)", sql, page.Current, page.PageSize)
	// 条件拼接完成后加上limit参数
	sql = fmt.Sprintf("(%s)  limit %d,%d", sql, page.Current, page.PageSize)
	// 执行sql
	var list []SearchInfo
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

// GetSearchTag 通过影片分类 Pid 返回对应分类的tag信息
func GetSearchTag(pid int64) map[string]interface{} {
	res := make(map[string]interface{})
	titles := db.Rdb.HGetAll(db.Cxt, fmt.Sprintf(config.SearchTitle, pid)).Val()
	for k, v := range titles {
		// 通过 k 获取对应的 tag , 并以score进行排序
		tags := db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, k), 0, 10).Val()
		res[v] = tags

		// 过滤分类tag
		switch k {
		case "Category", "Year", "Initial", "Sort":
			tags := db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, k), 0, -1).Val()
			res[v] = tags
		case "Plot":
			tags := db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, k), 0, 10).Val()
			res[v] = tags
		case "Area":
			tags := db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, k), 0, 11).Val()
			res[v] = tags
		case "Language":
			tags := db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, k), 0, 6).Val()
			res[v] = tags
		default:
			break
		}

	}
	return res
}

// ================================= 接口数据缓存 =================================

// DataCache  API请求 数据缓存
func DataCache(key string, data map[string]interface{}) {
	val, _ := json.Marshal(data)
	db.Rdb.Set(db.Cxt, key, val, time.Minute*30)
}

// GetCacheData 获取API接口的缓存数据
func GetCacheData(key string) map[string]interface{} {
	data := make(map[string]interface{})
	val, err := db.Rdb.Get(db.Cxt, key).Result()
	if err != nil || len(val) <= 0 {
		return nil
	}
	_ = json.Unmarshal([]byte(val), &data)
	return data
}

// RemoveCache 删除数据缓存
func RemoveCache(key string) {
	db.Rdb.Del(db.Cxt, key)
}
