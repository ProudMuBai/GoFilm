package logic

import (
	"errors"
	"fmt"
	"server/model/system"
	"server/plugin/common/conver"
	"time"
)

/*
处理影片管理相关业务
*/

type FilmLogic struct {
}

var FL *FilmLogic

//----------------------------------------------------影片管理业务逻辑----------------------------------------------------

func (fl *FilmLogic) GetFilmPage(s system.SearchVo) []system.SearchInfo {
	// 获取影片检索信息分页数据
	sl := system.GetSearchPage(s)
	//
	return sl
}

// GetSearchOptions 获取影片检索的select的选项options
func (fl *FilmLogic) GetSearchOptions() map[string]any {
	var options = make(map[string]any)
	// 获取分类 options
	tree := system.GetCategoryTree()
	tree.Name = "全部分类"
	options["class"] = conver.ConvertCategoryList(tree)
	options["remarks"] = []map[string]string{{"Name": `全部`, "Value": ``}, {"Name": `完结`, "Value": `完结`}, {"Name": `未完结`, "Value": `未完结`}}
	// 获取 剧情,地区,语言, 年份 组信息 (每个分类对应的检索信息并不相同)
	var tagGroup = make(map[int64]map[string]any)
	// 遍历一级分类获取对应的标签组信息
	for _, t := range tree.Children {
		option := system.GetSearchOptions(t.Id)
		if len(option) > 0 {
			tagGroup[t.Id] = system.GetSearchOptions(t.Id)
			// 如果年份信息不存在则独立一份年份信息
			if _, ok := options["year"]; !ok {
				options["year"] = tagGroup[t.Id]["Year"]
			}
		}

	}
	options["tags"] = tagGroup
	return options
}

// SaveFilmDetail 自定义上传保存影片信息
func (fl *FilmLogic) SaveFilmDetail(fd system.FilmDetailVo) error {
	// 补全影片信息
	now := time.Now()
	fd.UpdateTime = now.Format(time.DateTime)
	fd.AddTime = fd.UpdateTime
	// 生成ID, 由于是自定义上传的影片, 避免和采集站点的影片冲突, 以当前时间时间戳作为ID
	fd.Id = now.Unix()
	// 生成影片详情信息
	detail, err := conver.CovertFilmDetailVo(fd)
	if err != nil || detail.PlayList == nil {
		return errors.New("影片参数格式异常或缺少关键信息")
	}

	// 保存影片信息
	return system.SaveDetail(detail)
}

//----------------------------------------------------影片分类业务逻辑----------------------------------------------------

// GetFilmClassTree 获取影片分类信息
func (fl *FilmLogic) GetFilmClassTree() system.CategoryTree {
	// 获取原本的影片分类信息
	return system.GetCategoryTree()
}

func (fl *FilmLogic) GetFilmClassById(id int64) *system.CategoryTree {
	tree := system.GetCategoryTree()
	for _, c := range tree.Children {
		// 如果是一级分类, 则相等时直接返回
		if c.Id == id {
			return c
		}
		// 如果当前分类含有子分类, 则继续遍历匹配
		if c.Children != nil {
			for _, subC := range c.Children {
				if subC.Id == id {
					return subC
				}
			}
		}
	}
	return nil
}

// UpdateClass 更新分类信息
func (fl *FilmLogic) UpdateClass(class system.CategoryTree) error {
	// 遍历影片分类信息
	tree := system.GetCategoryTree()
	for _, c := range tree.Children {
		// 如果是一级分类, 则相等时直接修改对应的name和show属性
		if c.Id == class.Id {
			if class.Name != "" {
				c.Name = class.Name
			}
			c.Show = class.Show
			if err := system.SaveCategoryTree(&tree); err != nil {
				return fmt.Errorf("影片分类信息更新失败: %s", err.Error())
			}
			return nil
		}
		// 如果当前分类含有子分类, 则继续遍历匹配
		if c.Children != nil {
			for _, subC := range c.Children {
				if subC.Id == class.Id {
					if class.Name != "" {
						subC.Name = class.Name
					}
					subC.Show = class.Show
					if err := system.SaveCategoryTree(&tree); err != nil {
						return fmt.Errorf("影片分类信息更新失败: %s", err.Error())
					}
					return nil
				}
			}
		}
	}
	return errors.New("需要更新的分类信息不存在")
}

// DelClass 删除分类信息
func (fl *FilmLogic) DelClass(id int64) error {
	tree := system.GetCategoryTree()
	for i, c := range tree.Children {
		// 如果是一级分类, 则相等时直接返回
		if c.Id == id {
			tree.Children = append(tree.Children[:i], tree.Children[i+1:]...)
			if err := system.SaveCategoryTree(&tree); err != nil {
				return fmt.Errorf("影片分类信息删除失败: %s", err.Error())
			}
			return nil
		}
		// 如果当前分类含有子分类, 则继续遍历匹配
		if c.Children != nil {
			for j, subC := range c.Children {
				if subC.Id == id {
					c.Children = append(c.Children[:j], c.Children[j+1:]...)
					if err := system.SaveCategoryTree(&tree); err != nil {
						return fmt.Errorf("影片分类信息删除失败: %s", err.Error())
					}
					return nil
				}
			}
		}
	}
	return errors.New("需要删除的分类信息不存在")
}
