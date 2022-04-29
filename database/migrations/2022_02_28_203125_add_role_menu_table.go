package migrations

import (
	"database/sql"

	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type RoleMenu struct {
		RoleId int `gorm:"type:int(11);index:idx_role_menu,unique"`
		MenuId int `gorm:"type:int(11);index:idx_role_menu,unique"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&RoleMenu{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&RoleMenu{})
	}

	migrate.Add("2022_02_28_203125_add_role_menu_table", up, down)
}
