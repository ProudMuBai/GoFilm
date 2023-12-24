package system

// SearchTagsVO 搜索标签请求参数
type SearchTagsVO struct {
	Pid      int64  `json:"pid"`
	Cid      int64  `json:"cid"`
	Plot     string `json:"plot"`
	Area     string `json:"area"`
	Language string `json:"language"`
	Year     int64  `json:"year"`
	Sort     string `json:"sort"`
}

// FilmCronVo 影视更新任务请求参数
type FilmCronVo struct {
	Ids    []string `json:"ids"`    // 定时任务关联的资源站Id
	Time   int      `json:"time"`   // 更新最近几小时内更新的影片
	Spec   string   `json:"spec"`   // cron表达式
	Model  int      `json:"model"`  // 任务类型, 0 - 自动更新已启用站点 || 1 - 更新Ids中的资源站数据
	State  bool     `json:"state"`  // 任务状态 开启 | 关闭
	Remark string   `json:"remark"` // 备注信息
}

// CronTaskVo 定时任务数据response
type CronTaskVo struct {
	FilmCollectTask
	PreV string `json:"preV"` // 上次执行时间
	Next string `json:"next"` // 下次执行时间
}

// FilmTaskOptions 影视采集任务添加时需要的options
type FilmTaskOptions struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// CollectParams 数据采集所需要的参数
type CollectParams struct {
	Id    string   `json:"id"`    // 资源站id
	Ids   []string `json:"ids"`   // 资源站id列表
	Time  int      `json:"time"`  // 采集时长
	Batch bool     `json:"batch"` // 是否批量执行
}

// SearchVo 影片信息搜索参数
type SearchVo struct {
	Name     string `json:"name"`     // 影片名
	Pid      int64  `json:"pid"`      // 一级分类ID
	Cid      int64  `json:"cid"`      // 二级分类ID
	Plot     string `json:"plot"`     // 剧情
	Area     string `json:"area"`     // 地区
	Language string `json:"language"` // 语言
	Year     int64  `json:"year"`     // 年份
	//Score    int64  `json:"score"`    // 评分
	Remarks   string `json:"remarks"`   // 完结 | 未完结
	BeginTime int64  `json:"beginTime"` // 更新时间戳起始值
	EndTime   int64  `json:"endTime"`   // 更新时间戳结束值
	Paging    *Page  `json:"paging"`    // 分页参数
}

// FilmDetailVo 添加影片对象
type FilmDetailVo struct {
	Id           int64    `json:"id"`           // 影片id
	Cid          int64    `json:"cid"`          //分类ID
	Pid          int64    `json:"pid"`          //一级分类ID
	Name         string   `json:"name"`         //片名
	Picture      string   `json:"picture"`      //简介图片
	PlayFrom     []string `json:"playFrom"`     // 播放来源
	DownFrom     string   `json:"DownFrom"`     //下载来源 例: http
	PlayLink     string   `json:"playLink"`     //播放地址url
	DownloadLink string   `json:"downloadLink"` // 下载url地址
	SubTitle     string   `json:"subTitle"`     //子标题
	CName        string   `json:"cName"`        //分类名称
	EnName       string   `json:"enName"`       //英文名
	Initial      string   `json:"initial"`      //首字母
	ClassTag     string   `json:"classTag"`     //分类标签
	Actor        string   `json:"actor"`        //主演
	Director     string   `json:"director"`     //导演
	Writer       string   `json:"writer"`       //作者
	Remarks      string   `json:"remarks"`      // 更新情况
	ReleaseDate  string   `json:"releaseDate"`  //上映时间
	Area         string   `json:"area"`         // 地区
	Language     string   `json:"language"`     //语言
	Year         string   `json:"year"`         //年份
	State        string   `json:"state"`        //影片状态 正片|预告...
	UpdateTime   string   `json:"updateTime"`   //更新时间
	AddTime      string   `json:"addTime"`      //资源添加时间戳
	DbId         int64    `json:"dbId"`         //豆瓣id
	DbScore      string   `json:"dbScore"`      // 豆瓣评分
	Hits         int64    `json:"hits"`         //影片热度
	Content      string   `json:"content"`      //内容简介
}

// UserInfoVo 用户信息返回对象
type UserInfoVo struct {
	Id       uint   `json:"id"`
	UserName string `json:"userName"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Gender   int    `json:"gender"`   // 性别
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Status   int    `json:"status"`   // 状态
}

// PlayLinkVo 多站点播放链接数据列表
type PlayLinkVo struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	LinkList []MovieUrlInfo `json:"linkList"`
}

type MovieDetailVo struct {
	MovieDetail
	List []PlayLinkVo `json:"list"`
}
