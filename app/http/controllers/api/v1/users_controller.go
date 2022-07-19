package v1

import (
	"strconv"

	"github.com/spf13/cast"

	"api/app/models/menu"
	"api/app/models/role"
	"api/app/models/user"
	"api/app/requests"
	"api/pkg/auth"
	"api/pkg/database"
	"api/pkg/excelize"
	"api/pkg/file"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UsersController struct {
	BaseAPIController
}

// CurrentUser 当前登录用户信息
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	data := make(map[string]interface{})
	var roles = make([]string, 1)
	roles[0] = userModel.RoleID
	data["roles"] = roles
	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	if userModel.RoleID == "1" {
		data["permissions"] = permissions
	} else {
		menuIds := role.GetRoleMenuId(userModel.RoleID)
		list := menu.GetPermissions(menuIds)
		data["permissions"] = list
	}
	data["introduction"] = userModel.Introduction
	data["avatar"] = "https://img2.baidu.com/it/u=2436646421,1026055356&fm=253&fmt=auto&app=138&f=PNG?w=275&h=275"
	if userModel.Avatar != "" {
		data["avatar"] = userModel.Avatar
	}
	data["userId"] = userModel.ID
	data["name"] = userModel.Name
	response.Data(c, data)
}

// GetProfile 个人信息
func (ctrl *UsersController) GetProfile(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	data := make(map[string]interface{})
	data["phone"] = userModel.Phone
	data["email"] = userModel.Email
	data["city"] = userModel.City
	data["introduction"] = userModel.Introduction
	data["nick_name"] = userModel.NickName
	data["username"] = userModel.Name
	data["createdAt"] = userModel.CreatedAt
	data["roleName"] = role.Get(userModel.RoleID).RoleName
	response.Data(c, data)
}

// Index 所有用户
func (ctrl *UsersController) Index(c *gin.Context) {
	request := requests.UserPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}
	data, pager := user.Paginate(c, 10, request)
	if c.GetHeader("Http-Download") == "download" {
		dataKey, dataList := excelize.FormatDataExport(user.User{}, data)
		excel := excelize.NewMyExcel()
		excel.ExportToWeb(dataKey, dataList, c, "用户数据")
		return
	}
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Import 文件导入创建用户
func (ctrl *UsersController) Import(c *gin.Context) {
	rows, err := excelize.ImportToWeb("file", user.User{}, c)
	if err != nil {
		response.Error(c, err)
		return
	}
	isCover := c.PostForm("isCover")
	tx := database.DB.Begin()
	for _, row := range rows {
		userModel := user.User{}
		userModel.Name = row[0]
		userModel.RoleID = row[1]
		userModel.City = row[2]
		userModel.Introduction = row[3]
		userModel.Avatar = row[4]
		userModel.NickName = row[5]
		userModel.Email = row[6]
		userModel.Phone = row[7]
		userModel.Status = row[8]
		if isCover == "true" {
			users := user.GetByName(userModel.Name)
			userModel.ID = users.ID
		}
		err := error(nil)

		//验证
		request := requests.UserRequest{}
		//model复制到request验证
		err = copier.Copy(&request, userModel)
		if err != nil {
			tx.Rollback()
			response.Error(c, err)
			return
		}
		//由于复制的role_id类型不一致
		request.RoleID = cast.ToInt(row[1])
		if ok := requests.Validate(c, &request, requests.UserSave); !ok {
			tx.Rollback()
			return
		}
		if isCover == "true" {
			err = tx.Save(&userModel).Error
		} else {
			userModel.Password = "123456"
			err = tx.Create(&userModel).Error
		}
		if err != nil {
			tx.Rollback()
			response.Error(c, err)
			return
		}
	}
	tx.Commit()
	response.Success(c)
}

// UpdateProfile 修改个人信息
func (ctrl *UsersController) UpdateProfile(c *gin.Context) {

	request := requests.UserUpdateProfileRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateProfile); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.NickName = request.NickName
	currentUser.City = request.City
	currentUser.Introduction = request.Introduction
	rowsAffected := currentUser.Save()
	if rowsAffected > 0 {
		response.Data(c, currentUser)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

//Info 用户详情
func (ctrl *UsersController) Info(c *gin.Context) {
	userId := c.DefaultQuery("id", "")
	if userId == "" {
		response.NormalVerificationError(c, "用户id为空")
		return
	}
	users := user.Get(cast.ToInt(userId))
	response.Data(c, users)
}

// Update 修改用户信息
func (ctrl *UsersController) Update(c *gin.Context) {
	request := requests.UserRequest{}
	if ok := requests.Validate(c, &request, requests.UserSave); !ok {
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
	users := user.Get(request.Id)

	users.NickName = request.NickName
	users.Name = request.Name
	users.Phone = request.Phone
	users.Email = request.Email
	users.Status = request.Status
	if request.Password != "" {
		users.Password = request.Password
	}
	users.RoleID = strconv.Itoa(request.RoleID)
	users.Introduction = request.Introduction
	rowsAffected := users.Save()
	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加用户
func (ctrl *UsersController) Add(c *gin.Context) {
	// 1. 验证表单
	request := requests.UserRequest{}
	if ok := requests.Validate(c, &request, requests.UserSave); !ok {
		return
	}
	// 2. 验证成功，创建数据
	userModel := user.User{
		Name:         request.Name,
		Phone:        request.Phone,
		Email:        request.Email,
		Status:       request.Status,
		Password:     request.Password,
		RoleID:       strconv.Itoa(request.RoleID),
		NickName:     request.NickName,
		Introduction: request.Introduction,
	}
	userModel.Create()

	if userModel.ID > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "创建用户失败，请稍后尝试~")
	}
}

func (ctrl *UsersController) UpdateEmail(c *gin.Context) {

	request := requests.UserUpdateEmailRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateEmail); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Email = request.Email
	rowsAffected := currentUser.Save()

	if rowsAffected > 0 {
		response.Success(c)
	} else {
		// 失败，显示错误提示
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// ResetByPassword 重置密码
func (ctrl *UsersController) ResetByPassword(c *gin.Context) {
	// 1. 验证表单
	userId := cast.ToInt(c.DefaultQuery("id", ""))
	password := c.DefaultQuery("password", "")
	if userId == 0 || password == "" {
		response.NormalVerificationError(c, "参数为空")
		return
	}

	if ok := user.IsAdmin(userId); !ok {
		response.NormalVerificationError(c, "无法修改")
		return
	}

	// 2. 更新密码
	userModel := user.Get(userId)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = password
		userModel.Save()
		response.Success(c)
	}
}

//Status 状态修改
func (ctrl *UsersController) Status(c *gin.Context) {
	userId := cast.ToInt(c.DefaultQuery("id", ""))
	status := c.DefaultQuery("status", "")
	if userId == 0 || status == "" {
		response.NormalVerificationError(c, "参数为空")
		return
	}

	if ok := user.IsAdmin(userId); !ok {
		response.NormalVerificationError(c, "无法修改")
		return
	}

	userModel := user.Get(userId)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Status = status
		userModel.Save()
		response.Success(c)
	}
}

func (ctrl *UsersController) UpdatePhone(c *gin.Context) {

	request := requests.UserUpdatePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdatePhone); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Phone = request.Phone
	rowsAffected := currentUser.Save()

	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

func (ctrl *UsersController) UpdatePassword(c *gin.Context) {

	request := requests.UserUpdatePasswordRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdatePassword); !ok {
		return
	}

	currentUser := auth.CurrentUser(c)
	// 验证原始密码是否正确
	_, err := auth.Attempt(currentUser.Name, request.Password)
	if err != nil {
		// 失败，显示错误提示
		response.NormalVerificationError(c, "原密码不正确")
	} else {
		// 更新密码为新密码
		currentUser.Password = request.NewPassword
		currentUser.Save()

		response.Success(c)
	}
}

func (ctrl *UsersController) UpdateAvatar(c *gin.Context) {

	request := requests.UserUpdateAvatarRequest{}
	if ok := requests.Validate(c, &request, requests.UserUpdateAvatar); !ok {
		return
	}

	avatar, err := file.SaveUploadAvatar(c, request.Avatar)
	if err != nil {
		response.Abort500(c, "上传头像失败，请稍后尝试~")
		return
	}

	currentUser := auth.CurrentUser(c)
	currentUser.Avatar = avatar
	currentUser.Save()

	response.Data(c, currentUser)
}

//Delete 删除用户
func (ctrl *UsersController) Delete(c *gin.Context) {
	request := requests.UserDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.UserDelete); !ok {
		return
	}
	rowsAffected := user.DeleteIds(request.Ids, user.User{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
