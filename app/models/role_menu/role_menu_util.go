package role_menu

import (
	"api/app/models/menu"
	"api/app/models/role"
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (roleMenu RoleMenu) {
	database.DB.Where("id", idstr).First(&roleMenu)
	return
}

func GetBy(field, value string) (roleMenu RoleMenu) {
	database.DB.Where("? = ?", field, value).First(&roleMenu)
	return
}

func All() (roleMenus []RoleMenu) {
	database.DB.Find(&roleMenus)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(RoleMenu{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int) (roleMenus []RoleMenu, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(RoleMenu{}),
		&roleMenus,
		app.V1URL(database.TableName(&RoleMenu{})),
		perPage,
		"*",
	)
	return
}

func ReloadRule(roleId int, menuId []int) (err error) {
	menus := make([]menu.Menu, 0)
	roleMenu := make([]RoleMenu, len(menuId))
	//先删除所有的
	err = DeleteRoleMenu(roleId)
	if err != nil {
		return
	}
	err = database.DB.Where("id in (?)", menuId).
		Find(&menus).Error
	if err != nil {
		logger.Error("get menu error " + err.Error())
		return
	}
	for i := range menus {
		roleMenu[i] = RoleMenu{
			RoleId: roleId,
			MenuId: int(menus[i].ID),
		}
	}
	err = database.DB.Create(&roleMenu).Error
	if err != nil {
		logger.Error("batch create role's menu error, " + err.Error())
		return
	}
	return
}

func DeleteRoleMenu(roleId int) (err error) {
	err = database.DB.Where("role_id = ?", roleId).
		Delete(&RoleMenu{}).Error
	if err != nil {
		logger.Error("delete role's menu error, " + err.Error())
		return
	}
	var roles role.Role
	err = database.DB.Where("id = ?", roleId).
		First(&roles).Error
	if err != nil {
		logger.Error("get role error, " + err.Error())
		return
	}
	return
}
