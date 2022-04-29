package v1

import (
	"strconv"

	"api/app/models/dict_data"
	"api/app/requests"
	"api/pkg/excelize"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type DictDataController struct {
	BaseAPIController
}

// GetData 获取下拉数据
func (ctrl *DictDataController) GetData(c *gin.Context) {
	data := dict_data.Get(c.DefaultQuery("dictType", ""))
	list := make([]dict_data.DictDataGetAllResp, 0)
	for _, v := range data {
		d := dict_data.DictDataGetAllResp{}
		Transformation(v, &d)
		list = append(list, d)
	}
	response.Created(c, list)
}

// Index 获取字典数据列表
func (ctrl *DictDataController) Index(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id != "" {
		data := dict_data.GetId(id)
		response.JSON(c, gin.H{
			"data": data,
		})
		return
	}

	request := requests.DictDataPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.DictDataPagination); !ok {
		return
	}
	data, pager := dict_data.Paginate(c, 10, request)
	if c.DefaultQuery("export", "") == "1" {
		dataKey, dataList := excelize.FormatDataExport(dict_data.DictData{}, data)
		excel := excelize.NewMyExcel()
		excel.ExportToWeb(dataKey, dataList, c, "字典数据")
		return
	}
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Update 修改字典数据
func (ctrl *DictDataController) Update(c *gin.Context) {
	request := requests.DictDataRequest{}
	if ok := requests.Validate(c, &request, requests.DictDataSave); !ok {
		return
	}

	if request.DictCode == 0 {
		response.NormalVerificationError(c, "id为空")
		return
	}
	dictData := dict_data.GetId(strconv.Itoa(request.DictCode))
	dictData.DictLabel = request.DictLabel
	dictData.DictType = request.DictType
	dictData.DictSort = request.DictSort
	dictData.DictValue = request.DictValue
	dictData.Status = request.Status
	dictData.Remark = request.Remark

	rowsAffected := dictData.Save()
	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加字典数据
func (ctrl *DictDataController) Add(c *gin.Context) {
	request := requests.DictDataRequest{}
	if ok := requests.Validate(c, &request, requests.DictDataSave); !ok {
		return
	}
	dictDataModel := dict_data.DictData{
		DictSort:  request.DictSort,
		DictLabel: request.DictLabel,
		DictValue: request.DictValue,
		DictType:  request.DictType,
		Status:    request.Status,
		Remark:    request.Remark,
	}
	dictDataModel.Create()

	if dictDataModel.DictCode > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

func (ctrl *DictDataController) Delete(c *gin.Context) {
	request := requests.UserDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.UserDelete); !ok {
		return
	}
	rowsAffected := dict_data.DeleteIds(request.Ids, dict_data.DictData{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
