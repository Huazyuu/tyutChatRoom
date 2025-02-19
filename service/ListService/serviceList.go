package ListService

import (
	"fmt"
	"gin-gorilla/global"
	"gin-gorilla/model"
)

type Option struct {
	model.PageInfo
	Likes []string // 模糊匹配
}

// ComList ComList用于显示分页(可排序)
func ComList[T any](model T, option Option) (list []T, count int64, err error) {
	// model 用来推断 T 的 类型
	DB := global.DB
	// sort
	if option.Sort == "desc" {
		option.Sort = "created_at desc" // 时间倒序
	} else if option.Sort == "asc" {
		option.Sort = "created_at asc"
	} else {
		option.Sort = "created_at desc"
	}
	// cnt
	DB = DB.Where(model)
	for idx, like := range option.Likes {
		if idx == 0 {
			DB.Where(fmt.Sprintf("%s like ?", like), fmt.Sprintf("%%%s%%", option.Key))
			continue
		} else {
			DB.Or(fmt.Sprintf("%s like ?", like), fmt.Sprintf("%%%s%%", option.Key))
		}
	}
	count = DB.Find(&list).RowsAffected
	query := DB.Where(model)
	// 分页
	if option.Limit == 0 {
		option.Limit = 20
	}
	offset := (option.Page - 1) * option.Limit
	if offset < 0 {
		offset = 0
	}
	err = query.Limit(option.Limit).Offset(offset).Order(option.Sort).Find(&list).Error
	return list, count, err
}
