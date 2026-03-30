package system

import (
	"encoding/json"
	"log"
	"server/config"
	"server/plugin/db"
)

// Category 分类信息
type Category struct {
	Id   int64  `json:"id"`   // 分类ID
	Pid  int64  `json:"pid"`  // 父级分类ID
	Name string `json:"name"` // 分类名称
	Show bool   `json:"show"` // 是否展示
}

// CategoryTree 分类信息树形结构
type CategoryTree struct {
	*Category
	Children []*CategoryTree `json:"children"` // 子分类信息
}

// 影视分类展示树形结构

// SaveCategoryTree 保存影片分类信息
func SaveCategoryTree(tree *CategoryTree) error {
	data, _ := json.Marshal(tree)
	return db.Rdb.Set(db.Cxt, config.CategoryTreeKey, data, config.FilmExpired).Err()
}

// GetCategoryTree 获取影片分类信息
func GetCategoryTree() CategoryTree {
	data := db.Rdb.Get(db.Cxt, config.CategoryTreeKey).Val()
	tree := CategoryTree{}
	_ = json.Unmarshal([]byte(data), &tree)
	return tree
}

// ExistsCategoryTree 查询分类信息是否存在
func ExistsCategoryTree() bool {
	exists, err := db.Rdb.Exists(db.Cxt, config.CategoryTreeKey).Result()
	if err != nil {
		log.Println("ExistsCategoryTree Error", err)
	}
	return exists == 1
}

// GetChildrenTree 根据影片Id获取对应分类的子分类信息
func GetChildrenTree(id int64) []*CategoryTree {
	tree := GetCategoryTree()
	for _, t := range tree.Children {
		if t.Id == id {
			return t.Children
		}
	}
	return nil

}

// GetRevealCategoryList 获取show=true的二级分类信息
func GetRevealCategoryList() []*CategoryTree {
	// 初始化分类信息切片
	var cl []*CategoryTree
	var tree CategoryTree
	// 从redis获取分类信息树
	data := db.Rdb.Get(db.Cxt, config.CategoryTreeKey).Val()
	_ = json.Unmarshal([]byte(data), &tree)
	// 迭代分类树信息, 将show=true的二级分类信息添加到cl切片中
	for _, t := range tree.Children {
		if t.Show {
			for _, c := range t.Children {
				if c.Show {
					cl = append(cl, c)
				}
			}
		}
	}
	return cl
}
