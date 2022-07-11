package requests

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type RolePaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	BeginTime string `form:"beginTime" search:"type:gte;column:created_at;table:roles"`
	EndTime   string `form:"endTime" search:"type:lte;column:created_at;table:roles"`
	RoleName  string `form:"role_name" search:"type:contains;column:role_name;table:roles"`
	RoleKey   string `form:"role_key" search:"type:contains;column:role_key;table:roles"`
	Status    string `form:"status" search:"type:exact;column:status;table:roles"`
}

func RolePagination(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"sort":     []string{"in:id,created_at,role_sort"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:2,100"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 id,created_at,updated_at",
		},
		"order": []string{
			"in:排序规则仅支持 asc（正序）,desc（倒序）",
		},
		"per_page": []string{
			"numeric_between:每页条数的值介于 2~100 之间",
		},
	}
	return validate(data, rules, messages)
}

type RoleRequest struct {
	RoleName string `json:"role_name,omitempty" valid:"role_name" form:"role_name"`
	RoleKey  string `json:"role_key,omitempty" valid:"role_key" form:"role_key"`
	RoleSort int    `json:"role_sort" form:"role_sort"`
	Status   string `valid:"status" json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
	MenuIds  []int  `valid:"menu_ids" json:"menu_ids" form:"menu_ids"`

	Id int `valid:"id" json:"id" form:"id"`
}

func RoleSave(data interface{}, c *gin.Context) map[string][]string {
	_data := data.(*RoleRequest)
	rules := govalidator.MapData{
		"role_name": []string{"required", "not_exists:roles,role_name," + strconv.Itoa(_data.Id)},
		"role_key":  []string{"required", "not_exists:roles,role_key," + strconv.Itoa(_data.Id)},
		"menu_ids":  []string{"required"},
		"status":    []string{"required", "in:1,2"},
		"id":        []string{"exists:roles,id"},
	}

	messages := govalidator.MapData{
		"role_name": []string{
			"required:角色名为必填项",
			"not_exists:角色名 已被占用",
		},
		"role_key": []string{
			"required:权限字符为必填项",
			"not_exists:权限字符 已被占用",
		},
		"menu_ids": []string{
			"required:菜单为必填项",
		},
		"status": []string{
			"required:状态为必填项",
			"in:状态格式错误",
		},
		"id": []string{
			"exists:角色不存在",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
