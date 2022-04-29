package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type Role struct {
		models.BaseModel

		RoleName string `gorm:"type:varchar(60);not null;index"` //角色名称
		Status   int    `gorm:"type:smallint(2);default:2"`      //状态 默认2
		RoleSort int    `gorm:"type:int(11);default:1"`          //排序
		RoleKey  string `gorm:"type:varchar(60)"`                //角色代码
		Remark   string `gorm:"type:varchar(255)"`               //备注
		Admin    bool   `gorm:"type:tinyint(1);default:0"`       //是否管理员

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Role{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&Role{})
	}

	migrate.Add("2022_02_28_203109_add_role_table", up, down)
}
