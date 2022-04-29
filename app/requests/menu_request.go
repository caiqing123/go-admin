package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type MenuPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	Title   string `form:"title" search:"type:contains;column:title;table:menus"`
	Visible string `form:"visible" search:"type:exact;column:visible;table:menus"`
}

func MenuPagination(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"sort":     []string{"in:id,created_at,sort"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:2,2000"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 id,created_at,sort",
		},
		"order": []string{
			"in:排序规则仅支持 asc（正序）,desc（倒序）",
		},
		"per_page": []string{
			"numeric_between:每页条数的值介于 2~2000 之间",
		},
	}
	return validate(data, rules, messages)
}

type MenuRequest struct {
	IsFrame  string `valid:"is_frame" json:"is_frame" form:"is_frame"`
	Visible  string `valid:"visible" json:"visible" form:"visible"`
	Title    string `valid:"title" json:"title" form:"title"`
	MenuType string `valid:"menu_type" json:"menu_type" form:"menu_type"`

	Component  string `json:"component" form:"component"`
	Path       string `json:"path" form:"path"`
	MenuName   string `json:"menu_name,omitempty" form:"menu_name"`
	Permission string `json:"permission" form:"permission"`
	Action     string `json:"action,omitempty" form:"action"`
	ApiUrl     string `json:"api_url" form:"api_url"`
	Icon       string `json:"icon" form:"icon"`
	Sort       int    `json:"sort" form:"sort"`
	ParentId   int    `json:"parent_id" form:"parent_id"`
	Id         int    `valid:"id" json:"id" form:"id"`
}

func MenuSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"title":    []string{"required"},
		"menu_type": []string{"required", "in:M,C,F"},
		"visible":  []string{"required", "in:0,1"},
		"is_frame":  []string{"required", "in:0,1"},
		"id":       []string{"exists:menus,id"},
	}

	messages := govalidator.MapData{
		"title": []string{
			"required:菜单标题为必填项",
		},
		"menu_type": []string{
			"required:菜单类型为必填项",
			"in:菜单类型格式错误",
		},
		"visible": []string{
			"required:菜单状态为必填项",
			"in:菜单状态格式错误",
		},
		"is_frame": []string{
			"required:是否外链为必填项",
			"in:是否外链格式错误",
		},
		"id": []string{
			"exists:菜单不存在",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
