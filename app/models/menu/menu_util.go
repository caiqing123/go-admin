package menu

import (
	"api/app/models/role"
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (menu Menu) {
	database.DB.Where("id", idstr).First(&menu)
	return
}

func GetBy(value string) (menu Menu) {
	database.DB.Where("api_url = ?", value).First(&menu)
	return
}

func All() (menus []Menu) {
	database.DB.Find(&menus)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Menu{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int, q interface{}) (menus []Menu, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(Menu{}))
	menu := make([]Menu, 0)
	paging = paginator.Paginate(
		c,
		db,
		&menu,
		app.V1URL(database.TableName(&Menu{})),
		perPage,
		"*",
	)

	for i := 0; i < len(menu); i++ {
		if menu[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menu, menu[i])
		menus = append(menus, menusInfo)
	}
	return
}

// GetByRoleMenu 获取用户菜单
func GetByRoleMenu(roleId string) (m []Menu) {
	var MenuList []Menu
	if roleId == "1" {
		var data []Menu
		database.DB.Where(" menu_type in ('M','C')").Order("sort").Find(&data)
		MenuList = data
	} else {
		var dataC []Menu
		menuIds := role.GetRoleMenuId(roleId)
		database.DB.Where("menu_type in ('C')").Where("id in (?)", menuIds).Order("sort").Find(&dataC)
		for _, datum := range dataC {
			MenuList = append(MenuList, datum)
		}
		cIds := make([]int, 0)
		for _, menu := range MenuList {
			if menu.ParentId != 0 {
				cIds = append(cIds, int(menu.ParentId))
			}
		}
		var dataM []Menu
		database.DB.Where("menu_type in ('M')").Where("id in (?)", cIds).Order("sort").Find(&dataM)
		for _, datum := range dataM {
			MenuList = append(MenuList, datum)
		}

	}
	m = make([]Menu, 0)
	for i := 0; i < len(MenuList); i++ {
		if MenuList[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&MenuList, MenuList[i])
		m = append(m, menusInfo)
	}
	return
}

// GetList 获取菜单
func GetList() (m []MenuLabel) {
	var MenuList []Menu
	database.DB.Find(&MenuList)
	m = make([]MenuLabel, 0)
	for i := 0; i < len(MenuList); i++ {
		if MenuList[i].ParentId != 0 {
			continue
		}
		e := MenuLabel{}
		e.Id = MenuList[i].ID
		e.Label = MenuList[i].Title
		menusInfo := menuLabelCall(&MenuList, e)
		m = append(m, menusInfo)
	}
	return
}

// menuCall 构建菜单树
func menuLabelCall(menuList *[]Menu, menu MenuLabel) MenuLabel {
	list := *menuList
	min := make([]MenuLabel, 0)
	for j := 0; j < len(list); j++ {

		if menu.Id != list[j].ParentId {
			continue
		}
		mi := MenuLabel{}
		mi.Id = list[j].ID
		mi.Label = list[j].Title
		mi.Children = []MenuLabel{}
		if list[j].MenuType != "F" {
			ms := menuLabelCall(menuList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	if len(min) > 0 {
		menu.Children = min
	} else {
		menu.Children = nil
	}
	return menu
}

// menuCall 构建菜单树
func menuCall(menuList *[]Menu, menu Menu) Menu {
	list := *menuList
	min := make([]Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.ID != list[j].ParentId {
			continue
		}
		mi := list[j]
		if mi.MenuType != "F" {
			ms := menuCall(menuList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	menu.Children = min
	return menu
}

func GetPermissions(ids []int) (menu []interface{}) {
	database.DB.Table("menus").Where("id in (?)", ids).Where("permission != ''").Select("permission").Find(&menu)
	return
}

func DeleteIds(ids []int, menu Menu) (rowsAffected int64) {
	result := database.DB.Delete(menu, ids)
	return result.RowsAffected
}

func Count(value int) bool {
	var count int64
	database.DB.Model(Menu{}).Where("parent_id = ?", value).Count(&count)
	return count > 0
}
