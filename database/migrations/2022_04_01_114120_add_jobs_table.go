package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type Jobs struct {
		JobId uint64 `gorm:"column:job_id;primaryKey;autoIncrement;"`

		JobName        string `gorm:"type:varchar(255);comment:名称"`
		JobGroup       string `gorm:"type:varchar(255);comment:任务分组"`
		JobType        int    `gorm:"type:tinyint(1);default:1;comment:任务类型"`
		CronExpression string `gorm:"type:varchar(255);comment:cron表达式"`
		InvokeTarget   string `gorm:"type:varchar(255);comment:调用目标"`
		Args           string `gorm:"type:varchar(255);comment:目标参数"`
		MisfirePolicy  string `gorm:"type:varchar(255);comment:执行策略"`
		Concurrent     int    `gorm:"type:tinyint(1);default:1;comment:是否并发"`
		Status         int    `gorm:"type:tinyint(1);default:1;comment:状态"`
		EntryId        int    `gorm:"type:int(11);comment:job启动时返回的id"`

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Jobs{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&Jobs{})
	}

	migrate.Add("2022_04_01_114120_add_jobs_table", up, down)
}
