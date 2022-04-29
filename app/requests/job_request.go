package requests

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"

	jobs "api/pkg/job"
)

type JobPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	JobName  string `form:"job_name" search:"type:contains;column:job_name;table:jobs"`
	JobGroup string `form:"job_group" search:"type:contains;column:job_group;table:jobs"`
	Status   string `form:"status" search:"type:exact;column:status;table:jobs"`
}

func JobPagination(data interface{}, c *gin.Context) map[string][]string {
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

type JobRequest struct {
	JobName        string `valid:"job_name" json:"job_name" form:"job_name"`
	JobGroup       string `valid:"job_group" json:"job_group" form:"job_group"`
	InvokeTarget   string `valid:"invoke_target" json:"invoke_target" form:"invoke_target"`
	CronExpression string `valid:"cron_expression" json:"cron_expression" form:"cron_expression"`
	JobType        int    `valid:"job_type" json:"job_type" form:"job_type"`
	Status         int    `valid:"status" json:"status" form:"status"`
	MisfirePolicy  int    `valid:"misfire_policy" json:"misfire_policy" form:"misfire_policy"`

	Concurrent int    `json:"concurrent" form:"concurrent"`
	Args       string `json:"args" form:"args"`

	JobId int `valid:"job_id" json:"job_id" form:"job_id"`
}

func JobSave(data interface{}, c *gin.Context) map[string][]string {
	_data := data.(*JobRequest)
	rules := govalidator.MapData{
		"job_name":        []string{"required", "not_exists:jobs,job_name," + strconv.Itoa(_data.JobId) + ",job_id"},
		"job_group":       []string{"required"},
		"invoke_target":   []string{"required", "not_exists:jobs,invoke_target," + strconv.Itoa(_data.JobId) + ",job_id"},
		"cron_expression": []string{"required"},
		"misfire_policy":  []string{"required", "in:1,2"},
		"job_type":        []string{"required", "in:1,2"},
		"status":          []string{"required", "in:1,2"},
		"job_id":          []string{"exists:jobs,job_id"},
	}

	messages := govalidator.MapData{
		"job_name": []string{
			"required:名称为必填项",
			"not_exists:名称已被占用",
		},
		"job_group": []string{
			"required:任务分组为必填项",
		},
		"invoke_target": []string{
			"required:调用目标为必填项",
			"not_exists:调用目标已被占用",
		},
		"cron_expression": []string{
			"required:cron表达式为必填项",
		},
		"job_type": []string{
			"required:调用类型为必填项",
			"in:调用类型格式错误",
		},
		"misfire_policy": []string{
			"required:执行策略为必填项",
			"in:执行策略格式错误",
		},
		"status": []string{
			"required:状态为必填项",
			"in:状态格式错误",
		},
		"job_id": []string{
			"exists:任务id不存在",
		},
	}

	errs := validate(data, rules, messages)

	errs = ValidateCron(_data, errs)
	return errs
}

// ValidateCron 自定义规则，cron验证
func ValidateCron(cron *JobRequest, errs map[string][]string) map[string][]string {
	c, _ := regexp.MatchString("(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\\d+(ns|us|µs|ms|s|m|h))+)|((((\\d+,)+\\d+|(\\d+([/\\-])\\d+)|\\d+|\\*) ?){5,7})", cron.CronExpression)
	if c == false {
		errs["cron_expression"] = append(errs["cron_expression"], "cron表达式格式错误")
	}

	u, _ := regexp.MatchString("^(http://|https://)?((?:[A-Za-z\\d]+-[A-Za-z\\d]+|[A-Za-z\\d]+)\\.)+([A-Za-z]+)[/?:]?.*$", cron.InvokeTarget)
	if cron.JobType == 1 && u == false {
		errs["invoke_target"] = append(errs["invoke_target"], "调用目标不是url")
	}

	if cron.JobType == 2 && jobs.JobList[cron.InvokeTarget] == nil {
		errs["invoke_target"] = append(errs["invoke_target"], "调用目标 函数不存在")
	}
	return errs
}

type JobDeleteRequest struct {
	Ids []int `valid:"ids" form:"ids"`
}

func JobDelete(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"ids": []string{"required"},
	}
	messages := govalidator.MapData{
		"ids": []string{
			"required:id为空",
		},
	}
	errs := validate(data, rules, messages)

	return errs
}
