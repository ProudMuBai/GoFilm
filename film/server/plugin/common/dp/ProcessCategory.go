package dp

import (
	"server/model"
)

// =================Spider数据处理=======================

// CategoryTree 组装树形菜单
func CategoryTree(list []model.ClassInfo) *model.CategoryTree {
	// 遍历所有分类进行树形结构组装
	tree := &model.CategoryTree{Category: &model.Category{Id: 0, Pid: -1, Name: "分类信息"}}
	temp := make(map[int64]*model.CategoryTree)
	temp[tree.Id] = tree

	for _, c := range list {
		// 判断当前节点ID是否存在于 temp中
		category, ok := temp[c.Id]
		if ok {
			// 将当前节点信息保存
			category.Category = &model.Category{Id: c.Id, Pid: c.Pid, Name: c.Name}
		} else {
			// 如果不存在则将当前分类存放到 temp中
			category = &model.CategoryTree{Category: &model.Category{Id: c.Id, Pid: c.Pid, Name: c.Name}}
			temp[c.Id] = category
		}
		// 根据 pid获取父节点信息
		parent, ok := temp[category.Pid]
		if !ok {
			// 如果不存在父节点存在, 则将父节点存放到temp中
			temp[c.Pid] = parent
		}
		// 将当前节点存放到父节点的Children中
		parent.Children = append(parent.Children, category)
	}

	return tree
}
