package excelize

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
)

type lkExcelExport struct {
	file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewMyExcel() *lkExcelExport {
	return &lkExcelExport{file: createFile(), sheetName: defaultSheetName}
}

func FormatDataExport(key interface{}, data interface{}) (dataKey []map[string]string, dataList []map[string]interface{}) {
	//dataKey := make([]map[string]string, 0)
	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		title := fi.Tag.Get("title")
		key := strings.Split(fi.Tag.Get("json"), ",")[0]
		if title == "" {
			title = key
		}
		if key != "" {
			dataKey = append(dataKey, map[string]string{"key": key, "title": title})
		}
	}
	//dataList := make([]map[string]interface{}, 0)
	resByre, _ := json.Marshal(data)
	_ = json.Unmarshal(resByre, &dataList)
	return
}

// ExportToPath 导出基本的表格
func (l *lkExcelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, path string, title string) (string, error) {
	l.export(params, data)
	name := createFileName(title)
	filePath := path + "/" + name
	// 确保目录存在，不存在创建
	_ = os.MkdirAll(path, 0755)
	err := l.file.SaveAs(filePath)
	return filePath, err
}

// ExportToWeb 导出到浏览器。此处使用的gin框架 其他框架可自行修改ctx
func (l *lkExcelExport) ExportToWeb(params []map[string]string, data []map[string]interface{}, ctx *gin.Context, title string) {
	l.export(params, data)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	ctx.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName(title)))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	_, _ = ctx.Writer.Write(buffer.Bytes())
}

//设置首行
func (l *lkExcelExport) writeTop(params []map[string]string) {
	topStyle, _ := l.file.NewStyle(`{"font":{"bold":true},"alignment":{"horizontal":"center","vertical":"center"}}`)
	var word = 'A'
	//首行写入
	for _, conf := range params {
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)
		line := fmt.Sprintf("%c1", word)
		//设置标题
		_ = l.file.SetCellValue(l.sheetName, line, title)
		//列宽
		_ = l.file.SetColWidth(l.sheetName, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)
		//设置样式
		_ = l.file.SetCellStyle(l.sheetName, line, line, topStyle)
		word++
	}
}

//写入数据
func (l *lkExcelExport) writeData(params []map[string]string, data []map[string]interface{}) {
	lineStyle, _ := l.file.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	//数据写入
	var j = 2 //数据开始行数
	for i, val := range data {
		//设置行高
		_ = l.file.SetRowHeight(l.sheetName, i+1, defaultHeight)
		//逐列写入
		var word = 'A'
		for _, conf := range params {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)
			valNum := fmt.Sprintf("%v", val[valKey])

			//设置值
			if IsNum(valNum) {
				valNum, _ := strconv.Atoi(valNum)
				_ = l.file.SetCellValue(l.sheetName, line, valNum)
			} else {
				_ = l.file.SetCellValue(l.sheetName, line, val[valKey])
			}

			//设置样式
			_ = l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			word++
		}
		j++
	}
	//设置行高 尾行
	_ = l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (l *lkExcelExport) export(params []map[string]string, data []map[string]interface{}) {
	l.writeTop(params)
	l.writeData(params, data)
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := defaultSheetName
	index := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName(title string) string {
	name := time.Now().Format("2006-01-02-15-04-05")
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%v-%v-%v.xlsx", title, name, rand.Int63n(time.Now().Unix()))
}
