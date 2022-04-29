package v1

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"api/pkg/app"
	"api/pkg/auth"
	"api/pkg/file"
	"api/pkg/helpers"
	"api/pkg/response"
)

type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

const publicPath = "public"

type FileController struct {
	BaseAPIController
}

func (e *FileController) UploadFile(c *gin.Context) {
	tag, _ := c.GetPostForm("type")
	urlPrefix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse
	var done string
	switch tag {
	case "1": // 单图
		fileResponse, done = e.singleFile(c, fileResponse, urlPrefix)
		if done != "" {
			response.Abort500(c, done)
			return
		}
		response.Data(c, fileResponse)
		return
	case "2": // 多图
		multipartFile, done := e.multipleFile(c, urlPrefix)
		if done != "" {
			response.Abort500(c, done)
			return
		}
		response.Data(c, multipartFile)
		return
	case "3": // base64
		fileResponse, done = e.baseImg(c, fileResponse, urlPrefix)
		if done != "" {
			response.Abort500(c, done)
			return
		}
		response.Data(c, fileResponse)
	default:
		fileResponse, done = e.singleFile(c, fileResponse, urlPrefix)
		if done != "" {
			response.Abort500(c, done)
			return
		}
		response.Data(c, fileResponse)
		return
	}
	response.Abort500(c, "上传失败")
}

func (e FileController) baseImg(c *gin.Context, fileResponse FileResponse, urlPerfix string) (FileResponse, string) {
	var path = fmt.Sprintf("/uploads/uploadfile/%s/%s/", app.TimenowInTimezone().Format("2006/01/02"), auth.CurrentUID(c))
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
	fileName := helpers.RandomString(16) + ".jpg"
	err := file.IsNotExistMkDir(publicPath + path)
	if err != nil {
		return FileResponse{}, "初始化文件路径失败"
	}
	base64File := publicPath + path + fileName
	_ = ioutil.WriteFile(base64File, ddd, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
	fileResponse = FileResponse{
		Size:     file.GetFileSize(base64File),
		Path:     base64File,
		FullPath: urlPerfix + base64File,
		Name:     "",
		Type:     typeStr,
	}
	source, _ := c.GetPostForm("source")
	if source != "1" {
		fileResponse.Path = path + fileName
		fileResponse.FullPath = path + fileName
	}
	return fileResponse, ""
}

func (e FileController) multipleFile(c *gin.Context, urlPerfix string) ([]FileResponse, string) {
	var path = fmt.Sprintf("/uploads/uploadfile/%s/%s/", app.TimenowInTimezone().Format("2006/01/02"), auth.CurrentUID(c))
	files := c.Request.MultipartForm.File["file"]
	source, _ := c.GetPostForm("source")
	var multipartFile []FileResponse
	for _, f := range files {
		fileName := helpers.RandomString(16) + filepath.Ext(f.Filename)
		err := file.IsNotExistMkDir(publicPath + path)
		if err != nil {
			return []FileResponse{}, "初始化文件路径失败"
		}
		multipartFileName := publicPath + path + fileName
		err1 := c.SaveUploadedFile(f, multipartFileName)
		fileType, _ := file.GetType(multipartFileName)
		if err1 == nil {
			fileResponse := FileResponse{
				Size:     file.GetFileSize(multipartFileName),
				Path:     multipartFileName,
				FullPath: urlPerfix + multipartFileName,
				Name:     f.Filename,
				Type:     fileType,
			}
			if source != "1" {
				fileResponse.Path = path + fileName
				fileResponse.FullPath = path + fileName
			}
			multipartFile = append(multipartFile, fileResponse)
		}
	}
	return multipartFile, ""
}

func (e FileController) singleFile(c *gin.Context, fileResponse FileResponse, urlPerfix string) (FileResponse, string) {
	var path = fmt.Sprintf("/uploads/uploadfile/%s/%s/", app.TimenowInTimezone().Format("2006/01/02"), auth.CurrentUID(c))
	files, err := c.FormFile("file")
	if err != nil {
		return FileResponse{}, "图片不能为空"
	}
	// 上传文件至指定目录
	fileName := file.RandomNameFromUploadFile(files)

	err = file.IsNotExistMkDir(publicPath + path)
	if err != nil {
		return FileResponse{}, "初始化文件路径失败"
	}
	singleFile := publicPath + path + fileName
	_ = c.SaveUploadedFile(files, singleFile)
	fileType, _ := file.GetType(singleFile)
	fileResponse = FileResponse{
		Size:     file.GetFileSize(singleFile),
		Path:     singleFile,
		FullPath: urlPerfix + singleFile,
		Name:     files.Filename,
		Type:     fileType,
	}
	fileResponse.Path = path + fileName
	fileResponse.FullPath = path + fileName
	return fileResponse, ""
}
