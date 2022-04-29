package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type Menu struct {
		models.BaseModel

		MenuName   string `gorm:"type:varchar(128);not null;"` //菜单name
		Title      string `gorm:"type:varchar(128);not null;"` //显示名称
		Icon       string `gorm:"type:varchar(128);"`          //图标
		Path       string `gorm:"type:varchar(128)"`           //路径
		Paths      string `gorm:"type:varchar(128)"`           //id路径
		MenuType   string `gorm:"type:varchar(5)"`             //菜单类型
		Action     string `gorm:"type:varchar(16)"`            //请求方式
		ApiUrl     string `gorm:"type:varchar(255)"`           //后台路由格式
		Permission string `gorm:"type:varchar(255)"`           //权限编码
		ParentId   uint64 `gorm:"type:int(11);index"`          //上级菜单
		NoCache    int    `gorm:"type:tinyint(1);default:0"`   //是否缓存
		Breadcrumb string `gorm:"type:varchar(255)"`           //是否面包屑
		Component  string `gorm:"type:varchar(255)"`           //组件
		Sort       int    `gorm:"type:int(11)"`                //排序
		Visible    int    `gorm:"type:tinyint(1);default:1"`   //是否显示
		IsFrame    int    `gorm:"type:tinyint(1);default:0"`   //是否frame

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Menu{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&Menu{})
	}

	migrate.Add("2022_02_28_203149_add_menu_table", up, down)
}
