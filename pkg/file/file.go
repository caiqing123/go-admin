// Package file 文件操作辅助函数
package file

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api/pkg/app"
	"api/pkg/auth"
	"api/pkg/helpers"
	"api/pkg/logger"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func SaveUploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {

	var avatar string
	// 确保目录存在，不存在创建
	publicPath := "public"
	dirName := fmt.Sprintf("/uploads/avatars/%s/%s/", app.TimenowInTimezone().Format("2006/01/02"), auth.CurrentUID(c))
	_ = os.MkdirAll(publicPath+dirName, 0755)

	// 保存文件
	fileName := RandomNameFromUploadFile(file)
	// public/uploads/avatars/2021/12/22/1/nFDacgaWKpWWOmOt.png
	avatarPath := publicPath + dirName + fileName
	if err := c.SaveUploadedFile(file, avatarPath); err != nil {
		return avatar, err
	}

	// 裁切图片
	img, err := imaging.Open(avatarPath, imaging.AutoOrientation(true))
	if err != nil {
		return avatar, err
	}
	resizeAvatar := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
	resizeAvatarName := RandomNameFromUploadFile(file)
	resizeAvatarPath := publicPath + dirName + resizeAvatarName
	err = imaging.Save(resizeAvatar, resizeAvatarPath)
	if err != nil {
		return avatar, err
	}

	// 删除老文件
	err = os.Remove(avatarPath)
	if err != nil {
		return avatar, err
	}

	return dirName + resizeAvatarName, nil
}

func RandomNameFromUploadFile(file *multipart.FileHeader) string {
	return helpers.RandomString(16) + filepath.Ext(file.Filename)
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, 0755)
	if err != nil {
		return err
	}

	return nil
}

// CheckExist 检查文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// IsNotExistMkDir 检查文件夹是否存在
// 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if exist := !CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// GetType 获取文件类型
func GetType(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		logger.Error(err.Error())
	}
	filetype := http.DetectContentType(buff)
	return filetype, nil
}

// GetFileSize 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	_ = filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// GetCurrentPath 获取当前路径
func GetCurrentPath() string {
	dir, _ := os.Getwd()
	return strings.Replace(dir, "\\", "/", -1)
}

// FileMonitoring 监控文件变动执行对应方法
func FileMonitoring(ctx context.Context, filePth string, id string, group string, hookfn func(context.Context, string, string, []byte)) {
	f, err := os.Open(filePth)
	if err != nil {
		logger.Warn(err.Error())
		//log不存在 可不执行或循环
		//time.Sleep(1000 * time.Millisecond)
		//FileMonitoring(ctx, filePth, id, group, hookfn)
		return
	}
	//延迟 需等待websocket注册
	time.Sleep(1000 * time.Millisecond)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	rd := bufio.NewReader(f)
	//Seek将下一次读取或写入文件的偏移量设置为偏移量
	//1直接读取全部，2只读取最新的
	_, _ = f.Seek(0, 1)
	for {
		if ctx.Err() != nil {
			break
		}
		line, err := rd.ReadBytes('\n')
		// 如果是文件末尾不返回
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			logger.Error(err.Error())
		}
		//去除换行符
		line = bytes.TrimRight(line, "\r\n")
		go hookfn(ctx, id, group, line)
	}
}

func CronLog(format string, v ...interface{}) {
	err := IsNotExistMkDir("storage/cron/")
	if err != nil {
		logger.Error("cron dir error " + err.Error())
		return
	}
	logname := "storage/cron/" + time.Now().Format("2006-01-02.log")

	f, _ := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

	if _, err := f.Write([]byte(time.Now().Format("2006-01-02 15:04:05") + fmt.Sprintf(format, v...) + "\n")); err != nil {
		logger.Error("cron write error " + err.Error())
		return
	}
}
