package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type LoginLogPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	Username string `form:"username" search:"type:contains;column:username;table:login_logs"`
	IpAddr string `form:"ipaddr" search:"type:contains;column:ipaddr;table:login_logs"`
	Status   string `form:"status" search:"type:exact;column:status;table:login_logs"`
}

func LoginLogPagination(data interface{}, c *gin.Context) map[string][]string {
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



type OperaLogPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	BeginTime string `form:"beginTime" search:"type:gte;column:created_at;table:opera_logs"`
	EndTime   string `form:"endTime" search:"type:lte;column:created_at;table:opera_logs"`
	Status   string `form:"status" search:"type:exact;column:status;table:opera_logs"`
}

func OperaLogPagination(data interface{}, c *gin.Context) map[string][]string {
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

