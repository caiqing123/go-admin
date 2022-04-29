//Package job 模型
package job

import (
	"api/app/models"
	"api/pkg/database"
)

type Job struct {
	JobId int `gorm:"column:job_id;primaryKey;autoIncrement;" json:"job_id,omitempty"`

	JobName        string `json:"job_name"`        // 名称
	JobGroup       string `json:"job_group"`       // 任务分组
	JobType        int    `json:"job_type"`        // 任务类型
	CronExpression string `json:"cron_expression"` // cron表达式
	InvokeTarget   string `json:"invoke_target"`   // 调用目标
	Args           string `json:"args"`            // 目标参数
	MisfirePolicy  int    `json:"misfire_policy"`  // 执行策略
	Concurrent     int    `json:"concurrent"`      // 是否并发
	Status         int    `json:"status"`          // 状态
	EntryId        int    `json:"entry_id"`        // job启动时返回的id

	models.CommonTimestampsField
}

func (job *Job) Create() {
	database.DB.Create(&job)
}

func Save(job Job) (rowsAffected int64) {
	result := database.DB.Save(&job)
	return result.RowsAffected
}

func (job *Job) Save() (rowsAffected int64) {
	result := database.DB.Save(&job)
	return result.RowsAffected
}


func (job *Job) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&job)
	return result.RowsAffected
}
