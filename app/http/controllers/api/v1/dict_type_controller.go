package v1

import (
	"strconv"

	"api/app/models/dict_type"
	"api/app/requests"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type DictTypeController struct {
	BaseAPIController
}

// Index 获取字典列表
func (ctrl *DictTypeController) Index(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id != "" {
		data := dict_type.Get(id)
		response.JSON(c, gin.H{
			"data":    data,
		})
		return
	}

	request := requests.DictTypePaginationRequest{}
	if ok := requests.Validate(c, &request, requests.DictTypePagination); !ok {
		return
	}
	data, pager := dict_type.Paginate(c, 10, request)
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Update 修改字典类型
func (ctrl *DictTypeController) Update(c *gin.Context) {
	request := requests.DictTypeRequest{}
	if ok := requests.Validate(c, &request, requests.DictTypeSave); !ok {
		return
	}

	if request.DictId == 0 {
		response.NormalVerificationError(c, "id为空")
		return
	}
	menus := dict_type.Get(strconv.Itoa(request.DictId))
	menus.DictName = request.DictName
	menus.DictType = request.DictType
	menus.Status = request.Status
	menus.Remark = request.Remark

	rowsAffected := menus.Save()
	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加字典类型
func (ctrl *DictTypeController) Add(c *gin.Context) {
	request := requests.DictTypeRequest{}
	if ok := requests.Validate(c, &request, requests.DictTypeSave); !ok {
		return
	}
	dictTypeModel := dict_type.DictType{
		DictName: request.DictName,
		DictType: request.DictType,
		Status:   request.Status,
		Remark:   request.Remark,
	}
	dictTypeModel.Create()

	if dictTypeModel.DictId > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}


func (ctrl *DictTypeController) Delete(c *gin.Context) {
	request := requests.UserDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.UserDelete); !ok {
		return
	}

	rowsAffected := dict_type.DeleteIds(request.Ids,dict_type.DictType{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
