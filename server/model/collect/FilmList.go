package collect

import (
	"encoding/json"
	"encoding/xml"
	"server/config"
	"server/plugin/db"
)

/*
 视频列表接口序列化 struct
*/

//-------------------------------------------------Json 格式-------------------------------------------------

// CommonPage 影视列表接口分页数据结构体
type CommonPage struct {
	Code      int    `json:"code"`      // 响应状态码
	Msg       string `json:"msg"`       // 数据类型
	Page      any    `json:"page"`      // 页码
	PageCount int    `json:"pagecount"` // 总页数
	Limit     any    `json:"limit"`     // 每页数据量
	Total     int    `json:"total"`     // 总数据量
}

// FilmListPage 影视列表接口分页数据结构体
type FilmListPage struct {
	Code      int         `json:"code"`      // 响应状态码
	Msg       string      `json:"msg"`       // 数据类型
	Page      any         `json:"page"`      // 页码
	PageCount int         `json:"pagecount"` // 总页数
	Limit     any         `json:"limit"`     // 每页数据量
	Total     int         `json:"total"`     // 总数据量
	List      []FilmList  `json:"list"`      // 影片列表数据List集合
	Class     []FilmClass `json:"class"`     // 影片分类信息
}

// FilmList 影视列表单部影片信息结构体
type FilmList struct {
	VodID       int64  `json:"vod_id"`        // 影片ID
	VodName     string `json:"vod_name"`      // 影片名称
	TypeID      int64  `json:"type_id"`       // 分类ID
	TypeName    string `json:"type_name"`     // 分类名称
	VodEn       string `json:"vod_en"`        // 影片名中文拼音
	VodTime     string `json:"vod_time"`      // 更新时间
	VodRemarks  string `json:"vod_remarks"`   // 更新状态
	VodPlayFrom string `json:"vod_play_from"` // 播放来源
}

// FilmClass 影视分类信息结构体
type FilmClass struct {
	TypeID   int64  `json:"type_id"`   // 分类ID
	TypePid  int64  `json:"type_pid"`  // 父级ID
	TypeName string `json:"type_name"` // 类型名称
}

//-------------------------------------------------Xml 格式-------------------------------------------------

type RssL struct {
	XMLName xml.Name      `xml:"rss"`
	Version string        `xml:"version,attr"`
	List    FilmListPageX `xml:"list"`
	ClassXL ClassXL       `xml:"class"`
}
type FilmListPageX struct {
	XMLName     xml.Name    `xml:"list"`
	Page        any         `xml:"page,attr"`
	PageCount   int         `xml:"pagecount,attr"`
	PageSize    any         `xml:"pagesize,attr"`
	RecordCount int         `xml:"recordcount,attr"`
	Videos      []VideoList `xml:"video"`
}

type VideoList struct {
	Last string `xml:"last"`
	ID   int64  `xml:"id"`
	Tid  int64  `xml:"tid"`
	Name CDATA  `xml:"name"`
	Type string `xml:"type"`
	Dt   string `xml:"dt"`
	Note CDATA  `xml:"note"`
}

type ClassXL struct {
	XMLName xml.Name `xml:"class"`
	ClassX  []ClassX `xml:"ty"`
}

type ClassX struct {
	XMLName xml.Name `xml:"ty"`
	ID      int64    `xml:"id,attr"`
	Value   string   `xml:",chardata"`
}

//-------------------------------------------------redis Func-------------------------------------------------

// SaveFilmClass 保存影片分类列表信息
func SaveFilmClass(list []FilmClass) error {
	data, _ := json.Marshal(list)
	return db.Rdb.Set(db.Cxt, config.FilmClassKey, data, config.ResourceExpired).Err()
}

// GetFilmClass  获取分类列表信息
func GetFilmClass() []FilmClass {
	var l []FilmClass
	data := db.Rdb.Get(db.Cxt, config.FilmClassKey).Val()
	_ = json.Unmarshal([]byte(data), &l)
	return l
}
