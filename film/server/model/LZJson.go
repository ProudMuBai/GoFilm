package model

/*
量子资源JSON解析
*/

// ClassInfo class 分类数据
type ClassInfo struct {
	Id   int64  `json:"type_id"`   //分类ID
	Pid  int64  `json:"type_pid"`  //上级分类ID
	Name string `json:"type_name"` //分类名称
}

// MovieInfo 影片数据
type MovieInfo struct {
	Id       int64  `json:"vod_id"`        // 影片ID
	Name     string `json:"vod_name"`      // 影片名
	Cid      int64  `json:"type_id"`       // 所属分类ID
	CName    string `json:"type_name"`     // 所属分类名称
	EnName   string `json:"vod_en"`        // 英文片名
	Time     string `json:"vod_time"`      // 更新时间
	Remarks  string `json:"vod_remarks"`   // 备注 | 清晰度
	PlayFrom string `json:"vod_play_from"` // 播放来源
}

// MovieListInfo 影视列表响应数据
type MovieListInfo struct {
	Code      int64       `json:"code"`
	Msg       string      `json:"msg"`
	Page      string      `json:"page"`
	PageCount int64       `json:"pagecount"`
	Limit     string      `json:"limit"`
	Total     int64       `json:"total"`
	List      []MovieInfo `json:"list"`
	Class     []ClassInfo `json:"class"`
}

// MovieDetailInfo 影片详情数据 (只保留需要的部分)
type MovieDetailInfo struct {
	Id            int64  `json:"vod_id"`           //影片Id
	Cid           int64  `json:"type_id"`          //分类ID
	Pid           int64  `json:"type_id_1"`        //一级分类ID
	Name          string `json:"vod_name"`         //片名
	SubTitle      string `json:"vod_sub"`          //子标题
	CName         string `json:"type_name"`        //分类名称
	EnName        string `json:"vod_en"`           //英文名
	Initial       string `json:"vod_letter"`       //首字母
	ClassTag      string `json:"vod_class"`        //分类标签
	Pic           string `json:"vod_pic"`          //简介图片
	Actor         string `json:"vod_actor"`        //主演
	Director      string `json:"vod_director"`     //导演
	Writer        string `json:"vod_writer"`       //作者
	Blurb         string `json:"vod_blurb"`        //简介, 残缺,不建议使用
	Remarks       string `json:"vod_remarks"`      // 更新情况
	PubDate       string `json:"vod_pubdate"`      //上映时间
	Area          string `json:"vod_area"`         // 地区
	Language      string `json:"vod_lang"`         //语言
	Year          string `json:"vod_year"`         //年份
	State         string `json:"vod_state"`        //影片状态 正片|预告...
	UpdateTime    string `json:"vod_time"`         //更新时间
	AddTime       int64  `json:"vod_time_add"`     //资源添加时间戳
	DbId          int64  `json:"vod_douban_id"`    //豆瓣id
	DbScore       string `json:"vod_douban_score"` // 豆瓣评分
	Content       string `json:"vod_content"`      //内容简介
	PlayFrom      string `json:"vod_play_from"`    // 播放来源
	PlaySeparator string `json:"vod_play_note"`    // 播放信息分隔符
	PlayUrl       string `json:"vod_play_url"`     //播放地址url
	DownFrom      string `json:"vod_down_from"`    //下载来源 例: http
	DownUrl       string `json:"vod_down_url"`     // 下载url地址
}

// DetailListInfo 影视详情信息
type DetailListInfo struct {
	Code      int64             `json:"code"`
	Msg       string            `json:"msg"`
	Page      int64             `json:"page"`
	PageCount int64             `json:"pagecount"`
	Limit     string            `json:"limit"`
	Total     int64             `json:"total"`
	List      []MovieDetailInfo `json:"list"`
}
