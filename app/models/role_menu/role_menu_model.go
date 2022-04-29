//Package role_menu 模型
package role_menu

import (
    "api/pkg/database"
)

type RoleMenu struct {
    RoleId int `json:"role_id,omitempty"`
    MenuId int `json:"menu_id,omitempty"`
}

func (roleMenu *RoleMenu) Create() {
    database.DB.Create(&roleMenu)
}

func (roleMenu *RoleMenu) Save() (rowsAffected int64) {
    result := database.DB.Save(&roleMenu)
    return result.RowsAffected
}

func (roleMenu *RoleMenu) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&roleMenu)
    return result.RowsAffected
}