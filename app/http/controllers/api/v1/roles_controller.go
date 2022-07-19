package v1

import (
	"strconv"

	"api/app/models/role"
	"api/app/models/role_menu"
	"api/app/models/user"
	"api/app/requests"
	"api/pkg/excelize"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type RolesController struct {
	BaseAPIController
}

// GetRole 获取角色列表
func (ctrl *RolesController) GetRole(c *gin.Context) {
	request := requests.RolePaginationRequest{}
	if ok := requests.Validate(c, &request, requests.RolePagination); !ok {
		return
	}
	data, pager := role.Paginate(c, 10, request)
	if c.GetHeader("Http-Download") == "download" {
		dataKey, dataList := excelize.FormatDataExport(role.Role{}, data)
		excel := excelize.NewMyExcel()
		excel.ExportToWeb(dataKey, dataList, c, "角色数据")
		return
	}
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Update 修改角色
func (ctrl *RolesController) Update(c *gin.Context) {
	request := requests.RoleRequest{}
	if ok := requests.Validate(c, &request, requests.RoleSave); !ok {
		return
	}

	if request.Id == 0 {
		response.NormalVerificationError(c, "用户id为空")
		return
	}
	if ok := user.IsAdmin(request.Id); !ok {
		response.NormalVerificationError(c, "无法修改")
		return
	}
	roles := role.Get(strconv.Itoa(request.Id))
	status, _ := strconv.Atoi(request.Status)
	roles.RoleSort = request.RoleSort
	roles.Remark = request.Remark
	roles.Status = status

	rowsAffected := roles.Save()
	if rowsAffected > 0 {
		if len(request.MenuIds) > 0 {
			if err := role_menu.ReloadRule(int(roles.ID), request.MenuIds); err != nil {
				response.Abort500(c, "分配菜单权限失败")
				return
			}
		}
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加角色
func (ctrl *RolesController) Add(c *gin.Context) {
	// 1. 验证表单
	request := requests.RoleRequest{}
	if ok := requests.Validate(c, &request, requests.RoleSave); !ok {
		return
	}
	// 2. 验证成功，创建数据
	status, _ := strconv.Atoi(request.Status)
	roleModel := role.Role{
		RoleName: request.RoleName,
		RoleKey:  request.RoleKey,
		RoleSort: request.RoleSort,
		Status:   status,
		Remark:   request.Remark,
	}

	roleModel.Create()

	if roleModel.ID > 0 {
		if len(request.MenuIds) > 0 {
			if err := role_menu.ReloadRule(int(roleModel.ID), request.MenuIds); err != nil {
				response.Abort500(c, "分配菜单权限失败")
				return
			}
		}
		response.Success(c)
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

func (ctrl *RolesController) GetRoles(c *gin.Context) {
	roleId := c.DefaultQuery("role_id", "")
	if roleId == "" {
		response.NormalVerificationError(c, "角色id为空")
		return
	}
	data := role.Get(roleId)
	menuIds := make([]int, 0)
	menuIds = role.GetRoleMenuId(roleId)
	response.JSON(c, gin.H{
		"data":    data,
		"menuIds": menuIds,
	})
}

func (ctrl *RolesController) Delete(c *gin.Context) {
	request := requests.UserDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.UserDelete); !ok {
		return
	}
	for _, v := range request.Ids {
		err := role_menu.DeleteRoleMenu(v)
		if err != nil {
			response.NormalVerificationError(c, "ids错误")
			return
		}
	}
	rowsAffected := role.DeleteIds(request.Ids, role.Role{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
