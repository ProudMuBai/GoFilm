package logic

import (
	"errors"
	"server/model/system"
	"server/plugin/spider"
)

type CollectLogic struct {
}

var CollectL *CollectLogic

// ------------------------------------------------------ 采集站点管理 ------------------------------------------------------

// GetFilmSourceList 获取采集站列表数据
func (cl *CollectLogic) GetFilmSourceList() []system.FilmSource {
	// 返回当前已添加的采集站列表信息
	return system.GetCollectSourceList()
}

// GetFilmSource 获取ID对应的采集源信息
func (cl *CollectLogic) GetFilmSource(id string) *system.FilmSource {
	return system.FindCollectSourceById(id)
}

// UpdateFilmSource 更新采集源信息
func (cl *CollectLogic) UpdateFilmSource(s system.FilmSource) error {
	return system.UpdateCollectSource(s)
}

// SaveFilmSource  保存采集源信息
func (cl *CollectLogic) SaveFilmSource(s system.FilmSource) error {
	return system.AddCollectSource(s)
}

// DelFilmSource  删除采集源信息
func (cl *CollectLogic) DelFilmSource(id string) error {
	// 先查找是否存在对应ID的站点信息
	s := system.FindCollectSourceById(id)
	if s == nil {
		return errors.New("当前资源站信息不存在, 请勿重复操作")
	}
	//  如果是主站点则返回提示禁止直接删除
	if s.Grade == system.MasterCollect {
		return errors.New("主站点无法直接删除, 请先降级为附属站点再进行删除")
	}
	system.DelCollectResource(id)
	return nil
}

// ------------------------------------------------------ 采集记录管理 ------------------------------------------------------

// GetRecordList 获取采集记录列表
func (cl *CollectLogic) GetRecordList(params system.RecordRequestVo) []system.FailureRecord {
	return system.FailureRecordList(params)
}

// GetRecordOptions 获取采集记录筛选参数
func (cl *CollectLogic) GetRecordOptions() system.OptionGroup {
	var options = make(system.OptionGroup)
	// 获取筛选参数, 采集源(ID:name) | 采集类型 | 状态
	options["collectType"] = []system.Option{{"全部", -1}, {"影片详情", 0}, {"文章", 1}, {"演员", 2}, {"角色", 3}, {"网站", 4}}
	options["status"] = []system.Option{{"全部", -1}, {"待重试", 1}, {"已处理", 0}}
	// 获取全部采集站
	var originOptions = []system.Option{{"全部", ""}}
	for _, v := range system.GetCollectSourceList() {
		originOptions = append(originOptions, system.Option{Name: v.Name, Value: v.Id})
	}
	options["origin"] = originOptions
	return options
}

// CollectRecover 恢复采集
func (cl *CollectLogic) CollectRecover(id int) error {
	// 通过ID获取完整的失败记录信息
	fr := system.FindRecordById(uint(id))
	// 如果获取失败记录信息为空, 则不进行后续操作
	if fr == nil {
		return errors.New("采集重试执行失败: 失败记录信息获取异常")
	}
	// 执行恢复采集, 恢复对应的采集数据
	go spider.SingleRecoverSpider(fr)

	return nil
}

// RecoverAll 恢复重新对所有失效记录进行重新采集处理
func (cl *CollectLogic) RecoverAll() {
	// 是否进行身份验证, 暂定无需处理

	// 对数据表中的所有待处理记录进行恢复采集
	go spider.FullRecoverSpider()
}

// ClearDoneRecord 清除已处理完成的记录信息	(将记录表中已经完成处理的记录删除)
func (cl *CollectLogic) ClearDoneRecord() {
	//  <逻辑删除 or 真实删除> 为避免ID中断暂定逻辑删除
	system.DelDoneRecord()
}

// ClearAllRecord 清除所有记录信息	(直接对记录表直接进行截断处理)
func (cl *CollectLogic) ClearAllRecord() {
	// 重置记录表状态, 删除所有数据并将自增ID归零
	system.TruncateRecordTable()
}
