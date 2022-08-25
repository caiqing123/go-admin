package hotlist

import (
	"fmt"
	"testing"
)

func TestFileGet(t *testing.T) {

	spider := Spider{DataType: "BaiDu"}
	fmt.Println(spider.GetBaiDu())
	//All()

	//allData := []string{
	//	"GetZhiHu",
	//	"GetWeiBo",
	//	"GetDouBan",
	//	"GetTianYa",
	//	"GetHuPu",
	//	"GetBaiDu",
	//	"Get36Kr",
	//	"GetGuoKr",
	//	"GetHuXiu",
	//	"GetZHDaily",
	//	"GetSegmentfault",
	//	"GetWYNews",
	//	"GetWaterAndWood",
	//	"GetHacPai",
	//	"GetKD",
	//	"GetNGA",
	//	"GetWeiXin",
	//	"GetChiphell",
	//	"GetJianDan",
	//	"GetITHome",
	//}
	//
	//var person = &Spider{}
	//// 获取接口Person的类型对象
	//typeOfPerson := reflect.TypeOf(person)
	//// 打印Person的方法类型和名称
	//for i := 0; i < typeOfPerson.NumMethod(); i++ {
	//	//判断是否有使用对应方法
	//	if !collection.Collect(allData).Contains(typeOfPerson.Method(i).Name) {
	//		fmt.Printf("method is %s.\n", typeOfPerson.Method(i).Name)
	//	}
	//}
}
