package collect

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"server/config"
	"server/plugin/db"
)

/*
 视频详情接口序列化 struct
*/

//-------------------------------------------------Json 格式-------------------------------------------------

// FilmDetailLPage  视频详情分页数据
type FilmDetailLPage struct {
	Code      int          `json:"code"`      // 响应状态码
	Msg       string       `json:"msg"`       // 数据类型
	Page      any          `json:"page"`      // 页码
	PageCount int          `json:"pagecount"` // 总页数
	Limit     any          `json:"limit"`     // 每页数据量
	Total     int          `json:"total"`     // 总数据量
	List      []FilmDetail `json:"list"`      // 影片详情数据List集合
}

// FilmDetail 视频详情列表
type FilmDetail struct {
	VodID            int64  `json:"vod_id"`             // 影片ID
	TypeID           int64  `json:"type_id"`            // 分类ID
	TypeID1          int64  `json:"type_id_1"`          // 一级分类ID
	GroupID          int    `json:"group_id"`           // 用户组ID
	VodName          string `json:"vod_name"`           // 影片名称
	VodSub           string `json:"vod_sub"`            // 影片别名
	VodEn            string `json:"vod_en"`             //	影片名中文拼音
	VodStatus        int64  `json:"vod_status"`         // 影片状态
	VodLetter        string `json:"vod_letter"`         //	影片名首字母(大写)
	VodColor         string `json:"vod_color"`          //	UI展示颜色
	VodTag           string `json:"vod_tag"`            // 索引标签
	VodClass         string `json:"vod_class"`          // 剧情分类标签
	VodPic           string `json:"vod_pic"`            // 影片封面图
	VodPicThumb      string `json:"vod_pic_thumb"`      // 缩略图
	VodPicSlide      string `json:"vod_pic_slide"`      //  幻灯图片
	VodPicScreenshot string `json:"vod_pic_screenshot"` // ?截图
	VodActor         string `json:"vod_actor"`          // 演员名
	VodDirector      string `json:"vod_director"`       // 导演
	VodWriter        string `json:"vod_writer"`         // 作者
	VodBehind        string `json:"vod_behind"`         // 幕后
	VodBlurb         string `json:"vod_blurb"`          // 内容简介
	VodRemarks       string `json:"vod_remarks"`        // 更新状态 ( 完结 || 更新值 xx集)
	VodPubDate       string `json:"vod_pubdate"`        // 上映日期
	VodTotal         int64  `json:"vod_total"`          // 总集数
	VodSerial        string `json:"vod_serial"`         // 连载数
	VodTv            string `json:"vod_tv"`             // 上映电视台
	VodWeekday       string `json:"vod_weekday"`        // 节目周期
	VodArea          string `json:"vod_area"`           // 地区
	VodLang          string `json:"vod_lang"`           // 语言
	VodYear          string `json:"vod_year"`           // 年代
	VodVersion       string `json:"vod_version"`        // 画质版本 DVD || HD || 720P
	VodState         string `json:"vod_state"`          // 影片类别 正片 || 花絮 || 预告
	VodAuthor        string `json:"vod_author"`         // 编辑人员
	VodJumpUrl       string `json:"vod_jumpurl"`        // 跳转url
	VodTpl           string `json:"vod_tpl"`            // 独立模板
	VodTplPlay       string `json:"vod_tpl_play"`       // 独立播放页模板
	VodTplDown       string `json:"vod_tpl_down"`       // 独立下载页模板
	VodIsEnd         int64  `json:"vod_isend"`          // 是否完结
	VodLock          int64  `json:"vod_lock"`           // 锁定
	VodLevel         int64  `json:"vod_level"`          // 推荐级别
	VodCopyright     int64  `json:"vod_copyright"`      // 版权
	VodPoints        int64  `json:"vod_points"`         // 积分
	VodPointsPlay    int64  `json:"vod_points_play"`    // 点播付费
	VodPointsDown    int64  `json:"vod_points_down"`    // 下载付费
	VodHits          int64  `json:"vod_hits"`           // 总点击量
	VodHitsDay       int64  `json:"vod_hits_day"`       // 日点击量
	VodHitsWeek      int64  `json:"vod_hits_week"`      // 周点击量
	VodHitsMonth     int64  `json:"vod_hits_month"`     // 月点击量
	VodDuration      string `json:"vod_duration"`       // 时长
	VodUp            int64  `json:"vod_up"`             // 顶数
	VodDown          int64  `json:"vod_down"`           // 踩数
	VodScore         string `json:"vod_score"`          // 平均分
	VodScoreAll      int64  `json:"vod_score_all"`      // 总评分
	VodScoreNum      int64  `json:"vod_score_num"`      // 评分次数
	VodTime          string `json:"vod_time"`           // 更新时间
	VodTimeAdd       int64  `json:"vod_time_add"`       // 添加时间
	VodTimeHits      int64  `json:"vod_time_hits"`      // 点击时间
	VodTimeMake      int64  `json:"vod_time_make"`      // 生成时间
	VodTrySee        int64  `json:"vod_trysee"`         // 试看时长
	VodDouBanID      int64  `json:"vod_douban_id"`      // 豆瓣ID
	VodDouBanScore   string `json:"vod_douban_score"`   // 豆瓣评分
	VodReRrl         string `json:"vod_reurl"`          // 来源地址
	VodRelVod        string `json:"vod_rel_vod"`        // 关联视频ids
	VodRelArt        string `json:"vod_rel_art"`        // 关联文章 ids
	VodPwd           string `json:"vod_pwd"`            // 访问内容密码
	VodPwdURL        string `json:"vod_pwd_url"`        // 访问密码连接
	VodPwdPlay       string `json:"vod_pwd_play"`       // 访问播放页密码
	VodPwdPlayURL    string `json:"vod_pwd_play_url"`   // 获取访问密码连接
	VodPwdDown       string `json:"vod_pwd_down"`       // 访问下载页密码
	VodPwdDownURL    string `json:"vod_pwd_down_url"`   // 获取下载密码连接
	VodContent       string `json:"vod_content"`        // 详细介绍
	VodPlayFrom      string `json:"vod_play_from"`      // 播放组
	VodPlayServer    string `json:"vod_play_server"`    // 播放组服务器
	VodPlayNote      string `json:"vod_play_note"`      // 播放组备注 (分隔符)
	VodPlayURL       string `json:"vod_play_url"`       // 播放地址
	VodDownFrom      string `json:"vod_down_from"`      // 下载组
	VodDownServer    string `json:"vod_down_server"`    // 瞎子服务器组
	VodDownNote      string `json:"vod_down_note"`      // 下载备注 (分隔符)
	VodDownURL       string `json:"vod_down_url"`       // 下载地址
	VodPlot          int64  `json:"vod_plot"`           // 是否包含分级剧情
	VodPlotName      string `json:"vod_plot_name"`      // 分类剧情名称
	VodPlotDetail    string `json:"vod_plot_detail"`    // 分集剧情详情
	TypeName         string `json:"type_name"`          // 分类名称
}

//-------------------------------------------------Xml 格式-------------------------------------------------

type RssD struct {
	XMLName xml.Name        `xml:"rss"`
	Version string          `xml:"version,attr"`
	List    FilmDetailPageX `xml:"list"`
}

type CDATA struct {
	Text string `xml:",cdata"`
}

type FilmDetailPageX struct {
	XMLName     xml.Name      `xml:"list"`
	Page        string        `xml:"page,attr"`
	PageCount   int           `xml:"pagecount,attr"`
	PageSize    string        `xml:"pagesize,attr"`
	RecordCount int           `xml:"recordcount,attr"`
	Videos      []VideoDetail `xml:"video"`
}
type VideoDetail struct {
	XMLName  xml.Name `xml:"video"`
	Last     string   `xml:"last"`
	ID       int64    `xml:"id"`
	Tid      int64    `xml:"tid"`
	Name     CDATA    `xml:"name"`
	Type     string   `xml:"type"`
	Pic      string   `xml:"pic"`
	Lang     string   `xml:"lang"`
	Area     string   `xml:"area"`
	Year     string   `xml:"year"`
	State    string   `xml:"state"`
	Note     CDATA    `xml:"note"`
	Actor    CDATA    `xml:"actor"`
	Director CDATA    `xml:"director"`
	DL       DL       `xml:"dl"`
	Des      CDATA    `xml:"des"`
}

type DL struct {
	XMLName xml.Name `xml:"dl"`
	DD      []DD     `xml:"dd"`
}

type DD struct {
	XMLName xml.Name `xml:"dd"`
	Flag    string   `xml:"flag,attr"`
	Value   string   `xml:",cdata"`
}

//-------------------------------------------------Json 格式-------------------------------------------------

// BatchSaveOriginalDetail 批量保存原始影片详情数据
func BatchSaveOriginalDetail(dl []FilmDetail) {
	for _, d := range dl {
		SaveOriginalDetail(d)
	}
}

// SaveOriginalDetail 保存未处理的完整影片详情信息到redis
func SaveOriginalDetail(fd FilmDetail) {
	data, err := json.Marshal(fd)
	if err != nil {
		log.Println("Json Marshal FilmDetail Error: ", err)
	}
	if err = db.Rdb.Set(db.Cxt, fmt.Sprintf(config.OriginalFilmDetailKey, fd.VodID), data, config.ResourceExpired).Err(); err != nil {
		log.Println("Save Original FilmDetail Error: ", err)
	}
}

// GetOriginalDetailById 获取原始的影片详情数据
func GetOriginalDetailById(id int64) (FilmDetail, error) {
	data, err := db.Rdb.Get(db.Cxt, fmt.Sprintf(config.OriginalFilmDetailKey, id)).Result()
	if err != nil {
		log.Println("Get OriginalDetail Fail: ", err)
	}
	var fd = FilmDetail{}
	err = json.Unmarshal([]byte(data), &fd)
	if err != nil {
		log.Println("json.Unmarshal OriginalDetail Fail: ", err)
		return fd, err
	}

	return fd, nil

}
