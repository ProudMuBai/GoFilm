package system

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
	"server/plugin/common/param"
	"server/plugin/db"
	"strconv"
	"strings"
	"time"
)

// SearchInfo 存储用于检索的信息
type SearchInfo struct {
	gorm.Model
	Mid          int64   `json:"mid"`          //影片ID gorm:"uniqueIndex:idx_mid"
	Cid          int64   `json:"cid"`          //分类ID
	Pid          int64   `json:"pid"`          //上级分类ID
	Name         string  `json:"name"`         // 片名
	SubTitle     string  `json:"subTitle"`     // 影片子标题
	CName        string  `json:"cName"`        // 分类名称
	ClassTag     string  `json:"classTag"`     //类型标签
	Area         string  `json:"area"`         // 地区
	Language     string  `json:"language"`     // 语言
	Year         int64   `json:"year"`         // 年份
	Initial      string  `json:"initial"`      // 首字母
	Score        float64 `json:"score"`        //评分
	UpdateStamp  int64   `json:"updateStamp"`  // 更新时间
	Hits         int64   `json:"hits"`         // 热度排行
	State        string  `json:"state"`        //状态 正片|预告
	Remarks      string  `json:"remarks"`      // 完结 | 更新至x集
	ReleaseStamp int64   `json:"releaseStamp"` //上映时间 时间戳
}

// Tag 影片分类标签结构体
type Tag struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
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

// FilmZero 删除所有库存数据
func FilmZero() {
	// 删除redis中当前库存储的所有数据
	//db.Rdb.FlushDB(db.Cxt)
	db.Rdb.Del(db.Cxt, db.Rdb.Keys(db.Cxt, "MovieBasicInfoKey*").Val()...)
	db.Rdb.Del(db.Cxt, db.Rdb.Keys(db.Cxt, "MovieDetail*").Val()...)
	db.Rdb.Del(db.Cxt, db.Rdb.Keys(db.Cxt, "MultipleSource*").Val()...)
	db.Rdb.Del(db.Cxt, db.Rdb.Keys(db.Cxt, "OriginalResource*").Val()...)
	db.Rdb.Del(db.Cxt, db.Rdb.Keys(db.Cxt, "Search*").Val()...)
	// 删除mysql中留存的检索表
	var s *SearchInfo
	//db.Mdb.Exec(fmt.Sprintf(`drop table if exists %s`, s.TableName()))
	// 截断数据表 truncate table users
	if ExistSearchTable() {
		db.Mdb.Exec(fmt.Sprintf(`TRUNCATE table %s`, s.TableName()))
	}
}

// ResetSearchTable 重置Search表
func ResetSearchTable() {
	// 删除 Search 表
	var s *SearchInfo
	db.Mdb.Exec(fmt.Sprintf(`drop table if exists %s`, s.TableName()))
	// 重新创建 Search 表
	CreateSearchTable()
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
		searchMap = make(map[string]string)
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
				for _, t := range GetChildrenTree(search.Pid) {
					db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k),
						redis.Z{Score: float64(-t.Id), Member: fmt.Sprintf("%v:%v", t.Name, t.Id)})
				}
			}
		case "Year":
			// 获取 Year 数据, 如果不存在则缓存一份
			if tagCount == 0 {
				currentYear := time.Now().Year()
				for i := 0; i < 12; i++ {
					db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k),
						redis.Z{Score: float64(currentYear - i), Member: fmt.Sprintf("%v:%v", currentYear-i, currentYear-i)})
				}
			}
		case "Initial":
			// 如果不存在 首字母 Tag 数据, 则缓存一份
			if tagCount == 0 {
				for i := 65; i <= 90; i++ {
					db.Rdb.ZAdd(db.Cxt, fmt.Sprintf(config.SearchTag, search.Pid, k),
						redis.Z{Score: float64(90 - i), Member: fmt.Sprintf("%c:%c", i, i)})
				}
			}
		case "Sort":
			if tagCount == 0 {
				tags := []redis.Z{
					{3, "时间排序:update_stamp"},
					{2, "人气排序:hits"},
					{1, "评分排序:score"},
					{0, "最新上映:release_stamp"},
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
			score := db.Rdb.ZScore(db.Cxt, k, fmt.Sprintf("%v:%v", t, t)).Val()
			// 在原score的基础上+1 重新存入redis中
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: score + 1, Member: fmt.Sprintf("%v:%v", t, t)})
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
		if len(preTags) == 0 {
			// 如果没有 tag信息则不进行缓存
			//db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: 0, Member: fmt.Sprintf("%v:%v", "未知", "未知")})
		} else if preTags == "其它" {
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: 0, Member: fmt.Sprintf("%v:%v", preTags, preTags)})
		} else {
			score := db.Rdb.ZScore(db.Cxt, k, fmt.Sprintf("%v:%v", preTags, preTags)).Val()
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{Score: score + 1, Member: fmt.Sprintf("%v:%v", preTags, preTags)})
		}
	}
}

func BatchHandleSearchTag(infos ...SearchInfo) {
	for _, info := range infos {
		SaveSearchTag(info)
	}
}

// ================================= Spider 数据处理(mysql) =================================

// CreateSearchTable 创建存储检索信息的数据表
func CreateSearchTable() {
	// 如果不存在则创建表
	if !ExistSearchTable() {
		err := db.Mdb.AutoMigrate(&SearchInfo{})
		if err != nil {
			log.Println("Create Table SearchInfo Failed: ", err)
		}
	}
}

func ExistSearchTable() bool {
	// 1. 判断表中是否存在当前表
	return db.Mdb.Migrator().HasTable(&SearchInfo{})
}

// AddSearchIndex search表中数据保存完毕后 将常用字段添加索引提高查询效率
func AddSearchIndex() {
	var s *SearchInfo
	tableName := s.TableName()
	// 添加索引
	db.Mdb.Exec(fmt.Sprintf("CREATE UNIQUE INDEX idx_mid ON %s (mid)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_time ON %s (update_stamp DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_hits ON %s (hits DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_score ON %s (score DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_release ON %s (release_stamp DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_year ON %s (year DESC)", tableName))

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
	}
	// 保存成功后将相应tag数据缓存到redis中
	BatchHandleSearchTag(list...)
	tx.Commit()
}

// BatchSaveOrUpdate 判断数据库中是否存在对应mid的数据, 如果存在则更新, 否则插入
func BatchSaveOrUpdate(list []SearchInfo) {
	tx := db.Mdb.Begin()
	for _, info := range list {
		var count int64
		// 通过当前影片id 对应的记录数
		tx.Model(&SearchInfo{}).Where("mid", info.Mid).Count(&count)
		// 如果存在对应数据则进行更新, 否则保存相应数据
		if count > 0 {
			// 记录已经存在则执行更新部分内容
			err := tx.Model(&SearchInfo{}).Where("mid", info.Mid).Updates(SearchInfo{UpdateStamp: info.UpdateStamp, Hits: info.Hits, State: info.State,
				Remarks: info.Remarks, Score: info.Score, ReleaseStamp: info.ReleaseStamp}).Error
			if err != nil {
				tx.Rollback()
			}
		} else {
			// 执行插入操作
			if err := tx.Create(&info).Error; err != nil {
				tx.Rollback()
			}
			// 插入成功后保存一份tag信息到redis中
			BatchHandleSearchTag(info)
		}
	}
	// 提交事务
	tx.Commit()
}

// SaveSearchInfo 添加影片检索信息
func SaveSearchInfo(s SearchInfo) error {
	// 先查询数据库中是否存在对应记录
	// 如果不存在对应记录则 保存当前记录
	tx := db.Mdb.Begin()
	if !ExistSearchInfo(s.Mid) {
		// 执行插入操作
		if err := tx.Create(&s).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 执行添加操作时保存一份tag信息
		BatchHandleSearchTag(s)
	} else {
		// 如果已经存在当前记录则将当前记录进行更新
		err := tx.Model(&SearchInfo{}).Where("mid", s.Mid).Updates(SearchInfo{UpdateStamp: s.UpdateStamp, Hits: s.Hits, State: s.State,
			Remarks: s.Remarks, Score: s.Score, ReleaseStamp: s.ReleaseStamp}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 提交事务
	tx.Commit()
	return nil
}

// ExistSearchInfo 通过Mid查询是否存在影片的检索信息
func ExistSearchInfo(mid int64) bool {
	var count int64
	db.Mdb.Model(&SearchInfo{}).Where("mid", mid).Count(&count)
	return count > 0
}

// TunCateSearchTable 截断SearchInfo数据表
func TunCateSearchTable() {
	var searchInfo *SearchInfo
	err := db.Mdb.Exec(fmt.Sprint("TRUNCATE TABLE ", searchInfo.TableName())).Error
	if err != nil {
		log.Println("TRUNCATE TABLE Error: ", err)
	}
}

// SyncSearchInfo 同步影片检索信息
func SyncSearchInfo(model int) {
	switch model {
	case 0:
		// 重置Search表, (恢复为初始状态, 未添加索引)
		ResetSearchTable()
		// 批量添加 SearchInfo
		SearchInfoToMdb(model)
		// 保存完所有 SearchInfo 后添加字段索引
		AddSearchIndex()
	case 1:
		// 批量更新或添加
		SearchInfoToMdb(model)
	}
}

// SearchInfoToMdb 扫描redis中的检索信息, 并批量存入mysql (model 执行模式 0-清空并保存 || 1-更新)
func SearchInfoToMdb(model int) {
	// 获取集合中的元素数量, 如果集合中没有元素则直接返回
	count := db.Rdb.ZCard(db.Cxt, config.SearchInfoTemp).Val()
	if count <= 0 {
		return
	}
	// 1.从redis中批量扫描详情信息
	list := db.Rdb.ZPopMax(db.Cxt, config.SearchInfoTemp, config.MaxScanCount).Val()
	// 如果扫描到的信息为空则直接退出
	if len(list) <= 0 {
		return
	}
	// 2. 处理数据
	var sl []SearchInfo
	for _, s := range list {
		// 解析详情数据
		info := SearchInfo{}
		_ = json.Unmarshal([]byte(s.Member.(string)), &info)
		sl = append(sl, info)
	}
	// 通过model执行对应的保存方法
	switch model {
	case 0:
		// 批量添加 SearchInfo
		BatchSave(sl)
	case 1:
		// 批量更新或添加
		BatchSaveOrUpdate(sl)
	}
	//  如果 SearchInfoTemp 依然存在数据, 则递归执行
	SearchInfoToMdb(model)
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
func GetMultiplePlay(siteId, key string) []MovieUrlInfo {
	data := db.Rdb.HGet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetail, siteId), key).Val()
	var playList []MovieUrlInfo
	_ = json.Unmarshal([]byte(data), &playList)
	return playList
}

// GetSearchTag 通过影片分类 Pid 返回对应分类的tag信息
func GetSearchTag(pid int64) map[string]interface{} {
	// 整合searchTag相关内容
	res := make(map[string]interface{})
	titles := db.Rdb.HGetAll(db.Cxt, fmt.Sprintf(config.SearchTitle, pid)).Val()
	res["titles"] = titles
	// 处理单一分类的数据格式
	tagMap := make(map[string]interface{})
	for t, _ := range titles {
		tagMap[t] = HandleTagStr(t, GetTagsByTitle(pid, t)...)
	}
	res["tags"] = tagMap
	// 分类列表展示的顺序
	res["sortList"] = []string{"Category", "Plot", "Area", "Language", "Year", "Sort"}
	return res
}

// GetTagsByTitle 返回Pid和title对应的用于检索的tag
func GetTagsByTitle(pid int64, t string) []string {
	// 通过 k 获取对应的 tag , 并以score进行排序
	var tags []string
	// 过滤分类tag
	switch t {
	case "Category":
		tags = db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, t), 0, -1).Val()
	case "Plot":
		tags = db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, t), 0, 10).Val()
	case "Area":
		tags = db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, t), 0, 11).Val()
	case "Language":
		tags = db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, t), 0, 6).Val()
	case "Year", "Initial", "Sort":
		tags = db.Rdb.ZRevRange(db.Cxt, fmt.Sprintf(config.SearchTag, pid, t), 0, -1).Val()
	default:
		break
	}
	return tags
}

// HandleTagStr 处理tag数据格式
func HandleTagStr(title string, tags ...string) []map[string]string {
	var r []map[string]string
	if !strings.EqualFold(title, "Sort") {
		r = append(r, map[string]string{
			"Name":  "全部",
			"Value": "",
		})
	}
	for _, t := range tags {
		if sl := strings.Split(t, ":"); len(sl) > 0 {
			r = append(r, map[string]string{
				"Name":  sl[0],
				"Value": sl[1],
			})
		}
	}
	if !strings.EqualFold(title, "Sort") && !strings.EqualFold(title, "Year") && !strings.EqualFold(title, "Category") {
		r = append(r, map[string]string{
			"Name":  "其它",
			"Value": "其它",
		})
	}
	return r
}

// GetSearchInfosByTags 查询满足searchTag条件的影片分页数据
func GetSearchInfosByTags(st SearchTagsVO, page *Page) []SearchInfo {
	// 准备查询语句的条件
	qw := db.Mdb.Model(&SearchInfo{})
	// 通过searchTags的非空属性值, 拼接对应的查询条件
	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)
	for i := 0; i < t.NumField(); i++ {
		// 如果字段值不为空
		value := v.Field(i).Interface()
		if !param.IsEmpty(value) {
			// 如果value是 其它 则进行特殊处理
			var ts []string
			if v, flag := value.(string); flag && strings.EqualFold(v, "其它") {
				for _, s := range GetTagsByTitle(st.Pid, t.Field(i).Name) {
					ts = append(ts, strings.Split(s, ":")[1])
				}
			}
			k := strings.ToLower(t.Field(i).Name)
			switch k {
			case "pid", "cid", "year":
				qw = qw.Where(fmt.Sprintf("%s = ?", k), value)
			case "area", "language":
				if strings.EqualFold(value.(string), "其它") {
					qw = qw.Where(fmt.Sprintf("%s NOT IN ?", k), ts)
					break
				}
				qw = qw.Where(fmt.Sprintf("%s = ?", k), value)
			case "plot":
				if strings.EqualFold(value.(string), "其它") {
					for _, t := range ts {
						qw = qw.Where("class_tag NOT LIKE ?", fmt.Sprintf("%%%v%%", t))
					}
					break
				}
				qw = qw.Where("class_tag LIKE ?", fmt.Sprintf("%%%v%%", value))
			case "sort":
				if strings.EqualFold(value.(string), "release_stamp") {
					qw.Order(fmt.Sprintf("year DESC ,%v DESC", value))
					break
				}
				qw.Order(fmt.Sprintf("%v DESC", value))
			default:
				break
			}
		}
	}

	// 返回分页参数
	GetPage(qw, page)
	// 查询具体的searchInfo 分页数据
	var sl []SearchInfo
	if err := qw.Limit(page.PageSize).Offset((page.Current - 1) * page.PageSize).Find(&sl).Error; err != nil {
		log.Println(err)
		return nil
	}
	return sl

}

// GetMovieListBySort 通过排序类型返回对应的影片基本信息
func GetMovieListBySort(t int, pid int64, page *Page) []MovieBasicInfo {
	var sl []SearchInfo
	qw := db.Mdb.Model(&SearchInfo{}).Where("pid", pid).Limit(page.PageSize).Offset((page.Current) - 10*page.PageSize)
	// 针对不同排序类型返回对应的分页数据
	switch t {
	case 0:
		// 最新上映 (上映时间)
		qw.Order("year DESC, release_stamp DESC")
	case 1:
		// 排行榜 (暂定为热度排行)
		qw.Order("year DESC, hits DESC")
	case 2:
		// 最近更新 (更新时间)
		qw.Order("year DESC, update_stamp DESC")
	}
	if err := qw.Find(&sl).Error; err != nil {
		log.Println(err)
		return nil
	}
	return GetBasicInfoBySearchInfos(sl...)

}

// ================================= Manage 管理后台 =================================

func GetSearchPage(s SearchVo) []SearchInfo {
	// 构建 query查询条件
	query := db.Mdb.Model(&SearchInfo{})
	// 如果参数不为空则追加对应查询条件
	if s.Name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", s.Name))
	}
	// 分类ID为负数则默认不追加该条件
	if s.Cid > 0 {
		query = query.Where("cid = ?", s.Cid)
	} else if s.Pid > 0 {
		query = query.Where("pid = ?", s.Pid)
	}
	if s.Plot != "" {
		query = query.Where("class_tag LIKE ?", fmt.Sprintf("%%%s%%", s.Plot))
	}
	if s.Area != "" {
		query = query.Where("area = ?", s.Area)
	}
	if s.Language != "" {
		query = query.Where("language = ?", s.Language)
	}
	if int(s.Year) > time.Now().Year()-12 {
		query = query.Where("year = ?", s.Year)
	}
	switch s.Remarks {
	case "完结":
		query = query.Where("remarks IN ?", []string{"完结", "HD"})
	case "":
	default:
		query = query.Not(map[string]interface{}{"remarks": []string{"完结", "HD"}})
	}
	if s.BeginTime > 0 {
		query = query.Where("update_stamp >= ? ", s.BeginTime)
	}
	if s.EndTime > 0 {
		query = query.Where("update_stamp <= ? ", s.EndTime)
	}

	// 返回分页参数
	GetPage(query, s.Paging)
	// 查询具体的数据
	var sl []SearchInfo
	if err := query.Limit(s.Paging.PageSize).Offset((s.Paging.Current - 1) * s.Paging.PageSize).Find(&sl).Error; err != nil {
		log.Println(err)
		return nil
	}
	return sl

}

// GetSearchOptions 获取全部影片的检索标签信息
func GetSearchOptions(pid int64) map[string]interface{} {
	// 整合searchTag相关内容
	titles := db.Rdb.HGetAll(db.Cxt, fmt.Sprintf(config.SearchTitle, pid)).Val()
	// 处理单一分类的数据格式
	tagMap := make(map[string]interface{})
	for t, _ := range titles {
		switch t {
		// 只获取对应几个类型的标签
		case "Plot", "Area", "Language", "Year":
			tagMap[t] = HandleTagStr(t, GetTagsByTitle(pid, t)...)
		default:
		}
	}
	return tagMap
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

// ================================= OpenApi请求处理 =================================

func FindFilmIds(params map[string]string, page *Page) ([]int64, error) {
	var ids []int64
	query := db.Mdb.Model(&SearchInfo{}).Select("mid")
	for k, v := range params {
		// 如果 v 为空则直接 continue
		if len(v) <= 0 {
			continue
		}
		switch k {
		case "t":
			if cid, err := strconv.ParseInt(v, 10, 64); err == nil {
				query = query.Where("cid = ?", cid)
			}
		case "wd":
			query = query.Where("name like ?", fmt.Sprintf("%%%s%%", v))
		case "h":
			if h, err := strconv.ParseInt(v, 10, 64); err == nil {
				query = query.Where("update_stamp >= ?", time.Now().Unix()-h*3600)
			}
		}
	}
	// 返回分页参数
	var count int64
	query.Count(&count)
	page.Total = int(count)
	page.PageCount = int(page.Total+page.PageSize-1) / page.PageSize
	// 返回满足条件的ids
	err := query.Limit(page.PageSize).Offset(page.Current - 1).Order("update_stamp DESC").Find(&ids).Error
	return ids, err
}
