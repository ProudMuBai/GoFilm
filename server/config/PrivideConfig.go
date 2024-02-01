package config

import "time"

/*
	对外开放API相关配置
*/

const (
	// ResourceExpired API所需要的资源有效期
	ResourceExpired = time.Hour * 24 * 90
	// OriginalFilmDetailKey 采集时原始数据存储key
	OriginalFilmDetailKey = "OriginalResource:FilmDetail:Id%d"
	FilmClassKey          = "OriginalResource:FilmClass"
	PlayForm              = "gfm3u8"
	PlayFormCloud         = "gofilm"
	PlayFormAll           = "gofilm$$$gfmu38"
	RssVersion            = "5.1"
)
