package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type DictDataPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	DictLabel string `form:"dict_label" search:"type:contains;column:dict_label;table:dict_data"`
	DictType  string `form:"dict_type" search:"type:contains;column:dict_type;table:dict_data"`
	Status    string `form:"status" search:"type:exact;column:status;table:dict_data"`
}

func DictDataPagination(data interface{}, c *gin.Context) map[string][]string {
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

type DictDataRequest struct {
	DictLabel string `valid:"dict_label" json:"dict_label" form:"dict_label"`
	DictType  string `valid:"dict_type" json:"dict_type" form:"dict_type"`
	DictValue string `valid:"dict_value" json:"dict_value" form:"dict_value"`
	Status    int    `valid:"status" json:"status" form:"status"`

	DictSort int    `json:"dict_sort" form:"dict_sort"`
	Remark   string `json:"remark" form:"remark"`
	DictCode int    `valid:"dict_code" json:"dict_code" form:"dict_code"`
}

func DictDataSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"dict_label": []string{"required"},
		"dict_value": []string{"required"},
		"dict_type":  []string{"required", "exists:dict_types,dict_type"},
		"status":     []string{"required", "in:1,2"},
		"dict_code":  []string{"exists:dict_data,dict_code"},
	}

	messages := govalidator.MapData{
		"dict_label": []string{
			"required:数据标签为必填项",
		},
		"dict_value": []string{
			"required:数据键值为必填项",
		},
		"dict_type": []string{
			"required:字典类型为必填项",
			"exists:字典类型不存在",
		},
		"status": []string{
			"required:状态为必填项",
			"in:状态格式错误",
		},
		"dict_code": []string{
			"exists:字典id不存在",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
