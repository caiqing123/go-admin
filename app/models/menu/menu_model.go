//Package menu 模型
package menu

import (
	"strconv"

	"api/app/models"
	"api/pkg/database"
)

type Menu struct {
	models.BaseModel

	MenuName   string `json:"menu_name,omitempty"`
	Title      string `json:"title,omitempty"`
	Path       string `json:"path,omitempty"`
	Paths      string `json:"paths,omitempty"`
	MenuType   string `json:"menu_type,omitempty"`
	Action     string `json:"action,omitempty"`
	Permission string `json:"permission,omitempty"`
	NoCache    string `json:"no_cache"`
	Breadcrumb string `json:"breadcrumb,omitempty"`
	Component  string `json:"component,omitempty"`
	Visible    string `json:"visible"`
	IsFrame    string `json:"is_frame"`
	ParentId   uint64 `json:"parent_id"`
	Sort       int    `json:"sort"`
	Icon       string `json:"icon"`
	ApiUrl     string `json:"api_url"`
	Children   []Menu `json:"children,omitempty" gorm:"-"`

	models.CommonTimestampsField
}

type MenuLabel struct {
	Id       uint64      `json:"id,omitempty" gorm:"-"`
	Label    string      `json:"label,omitempty" gorm:"-"`
	Children []MenuLabel `json:"children,omitempty" gorm:"-"`
}

func (menu *Menu) Create() {
	database.DB.Create(&menu)
}

func (menu *Menu) Save() (rowsAffected int64) {
	result := database.DB.Save(&menu)
	return result.RowsAffected
}

func (menu *Menu) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&menu)
	return result.RowsAffected
}

func (menu *Menu) InitPaths() {
	var data Menu
	parentMenu := new(Menu)
	if menu.ParentId != 0 {
		database.DB.Model(&data).First(parentMenu, menu.ParentId)
		menu.Paths = parentMenu.Paths + "/" + strconv.FormatUint(menu.ID, 10)
	} else {
		menu.Paths = "/0/" + strconv.FormatUint(menu.ID, 10)
	}
	database.DB.Model(&data).Where("id = ?", menu.ID).Update("paths", menu.Paths)
}
