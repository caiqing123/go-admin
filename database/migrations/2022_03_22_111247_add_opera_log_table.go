package migrations

import (
	"database/sql"
	"time"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type OperaLog struct {
		models.BaseModel

		Title         string    `gorm:"type:varchar(255);comment:操作模块"`
		BusinessType  string    `gorm:"type:varchar(128);comment:操作类型"`
		Method        string    `gorm:"type:varchar(128);comment:函数"`
		RequestMethod string    `gorm:"type:varchar(128);comment:请求方式"`
		OperatorType  string    `gorm:"type:varchar(128);comment:操作类型"`
		OperaName     string    `gorm:"type:varchar(128);comment:操作者"`
		OperaUrl      string    `gorm:"type:varchar(255);comment:访问地址"`
		OperaIp       string    `gorm:"type:varchar(128);comment:客户端ip"`
		OperaLocation string    `gorm:"type:varchar(128);comment:访问位置"`
		OperaParam    string    `gorm:"type:varchar(255);comment:请求参数"`
		Status        string    `gorm:"type:varchar(4);comment:操作状态"`
		OperaTime     time.Time `gorm:"column:opera_time;index;comment:操作时间"`
		JsonResult    string    `gorm:"type:varchar(255);comment:返回数据"`
		Remark        string    `gorm:"type:varchar(255);comment:备注"`
		LatencyTime   string    `gorm:"type:varchar(128);comment:耗时"`
		UserAgent     string    `gorm:"type:varchar(255);comment:ua"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&OperaLog{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&OperaLog{})
	}

	migrate.Add("2022_03_22_111247_add_opera_log_table", up, down)
}
