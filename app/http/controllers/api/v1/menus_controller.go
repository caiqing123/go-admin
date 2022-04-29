package v1

import (
	"strconv"

	"api/app/models/menu"
	"api/app/models/role"
	"api/app/requests"
	"api/pkg/auth"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type MenusController struct {
	BaseAPIController
}

// GetMenuRole 根据登录角色名称获取菜单列表数据
func (ctrl *MenusController) GetMenuRole(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	data := menu.GetByRoleMenu(userModel.RoleID)
	response.Data(c, data)
}

// GetMenuTreeSelect 根据角色id获取菜单列表数据
func (ctrl *MenusController) GetMenuTreeSelect(c *gin.Context) {
	data := menu.GetList()
	roleId := c.DefaultQuery("role_id", "")
	menuIds := make([]int, 0)
	if roleId != "" {
		menuIds = role.GetRoleMenuId(roleId)
	}
	response.JSON(c, gin.H{
		"menus":       data,
		"checkedKeys": menuIds,
	})
}

// Index 获取菜单列表
func (ctrl *MenusController) Index(c *gin.Context) {
	request := requests.MenuPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.MenuPagination); !ok {
		return
	}
	data, pager := menu.Paginate(c, 1000, request)
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Update 修改菜单
func (ctrl *MenusController) Update(c *gin.Context) {
	request := requests.MenuRequest{}
	if ok := requests.Validate(c, &request, requests.MenuSave); !ok {
		return
	}

	if request.Id == 0 {
		response.NormalVerificationError(c, "id为空")
		return
	}
	menus := menu.Get(strconv.Itoa(request.Id))
	menus.MenuName = request.MenuName
	menus.Title = request.Title
	menus.Path = request.Path
	menus.MenuType = request.MenuType
	menus.Action = request.Action
	menus.Permission = request.Permission
	menus.Component = request.Component
	menus.Visible = request.Visible
	menus.IsFrame = request.IsFrame
	menus.ParentId = uint64(request.ParentId)
	menus.Sort = request.Sort
	menus.Icon = request.Icon
	menus.ApiUrl = request.ApiUrl

	rowsAffected := menus.Save()
	if rowsAffected > 0 {
		menus.InitPaths()
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加菜单
func (ctrl *MenusController) Add(c *gin.Context) {
	request := requests.MenuRequest{}
	if ok := requests.Validate(c, &request, requests.MenuSave); !ok {
		return
	}
	menuModel := menu.Menu{
		MenuName:   request.MenuName,
		Title:      request.Title,
		Path:       request.Path,
		MenuType:   request.MenuType,
		Action:     request.Action,
		Permission: request.Permission,
		Component:  request.Component,
		Visible:    request.Visible,
		IsFrame:    request.IsFrame,
		ParentId:   uint64(request.ParentId),
		Sort:       request.Sort,
		Icon:       request.Icon,
		ApiUrl:     request.ApiUrl,
	}
	menuModel.Create()

	if menuModel.ID > 0 {
		menuModel.InitPaths()
		response.Success(c)
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

func (ctrl *MenusController) GetMenus(c *gin.Context) {
	menuId := c.DefaultQuery("id", "")
	if menuId == "" {
		response.NormalVerificationError(c, "菜单id为空")
		return
	}
	data := menu.Get(menuId)
	response.JSON(c, gin.H{
		"data": data,
	})
}

func (ctrl *MenusController) Delete(c *gin.Context) {
	request := requests.UserDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.UserDelete); !ok {
		return
	}
	for _, v := range request.Ids {
		men := menu.Count(v)
		if men {
			response.NormalVerificationError(c, "子菜单未删除")
			return
		}
	}
	rowsAffected := menu.DeleteIds(request.Ids, menu.Menu{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
