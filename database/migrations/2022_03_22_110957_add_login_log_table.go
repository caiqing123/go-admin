package migrations

import (
	"database/sql"
	"time"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type LoginLog struct {
		models.BaseModel

		Username      string    `gorm:"type:varchar(128);comment:用户名"`
		Status        string    `gorm:"type:varchar(4);comment:状态"`
		Ipaddr        string    `gorm:"type:varchar(255);comment:ip地址"`
		LoginLocation string    `gorm:"type:varchar(255);comment:归属地"`
		Browser       string    `gorm:"type:varchar(255);comment:浏览器"`
		Os            string    `gorm:"type:varchar(255);comment:系统"`
		Platform      string    `gorm:"type:varchar(255);comment:固件"`
		LoginTime     time.Time `gorm:"column:login_time;index;comment:登录时间"`
		Remark        string    `gorm:"type:varchar(255);comment:备注"`
		Msg           string    `gorm:"type:varchar(255);comment:信息"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&LoginLog{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&LoginLog{})
	}

	migrate.Add("2022_03_22_110957_add_login_log_table", up, down)
}
