package requests

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type DictTypePaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	DictName string `form:"dict_name" search:"type:contains;column:dict_name;table:dict_types"`
	DictType string `form:"dict_type" search:"type:contains;column:dict_type;table:dict_types"`
	Status   string `form:"status" search:"type:exact;column:status;table:dict_types"`
}

func DictTypePagination(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"sort":     []string{"in:created_at"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:2,100"},
	}
	messages := govalidator.MapData{
		"sort": []string{
			"in:排序字段仅支持 created_at",
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

type DictTypeRequest struct {
	DictName string `valid:"dict_name" json:"dict_name" form:"dict_name"`
	DictType string `valid:"dict_type" json:"dict_type" form:"dict_type"`
	Status   int    `valid:"status" json:"status" form:"status"`

	Remark string `json:"remark" form:"remark"`
	DictId int    `valid:"id" json:"id" form:"id"`
}

func DictTypeSave(data interface{}, c *gin.Context) map[string][]string {
	_data := data.(*DictTypeRequest)
	rules := govalidator.MapData{
		"dict_name": []string{"required", "not_exists:dict_types,dict_name," + strconv.Itoa(_data.DictId) + ",dict_id"},
		"dict_type": []string{"required", "not_exists:dict_types,dict_type," + strconv.Itoa(_data.DictId) + ",dict_id"},
		"status":    []string{"required", "in:1,2"},
		"id":        []string{"exists:dict_types,dict_id"},
	}

	messages := govalidator.MapData{
		"dict_name": []string{
			"required:字典名称为必填项",
		},
		"dict_type": []string{
			"required:字典类型为必填项",
		},
		"status": []string{
			"required:状态为必填项",
			"in:状态格式错误",
		},
		"id": []string{
			"exists:字典id不存在",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
