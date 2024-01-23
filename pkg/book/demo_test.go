package book

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"api/pkg/book/site"
	"api/pkg/book/store"
)

func init() {
	InitSites()
}

func TestDemoSearch(t *testing.T) {
	//搜索
	//for _, s := range site.SitePool {
	//	if s.Search == nil {
	//		continue
	//	}
	//	result, err := s.Search("我的诡异人生模拟器")
	//	fmt.Println(err)
	//	fmt.Println(result)
	//
	//}

	//详情
	result, err := site.BookInfo("http://www.b520.cc/159_159394/")
	//fmt.Println(result)
	fmt.Println(err)

	//目录详情
	//result, err := site.Chapter("https://segmentfault.com/a/1190000023240989")
	//fmt.Println(err)
	//fmt.Println(result)

	//保存txt
	//for k, v := range result.Volumes {
	//	for k1, v1 := range v.Chapters {
	//		result.Volumes[k].Chapters[k1].Text, _ = site.Chapter(v1.URL)
	//	}
	//}

	//保存txt
	site.Download(result)
	start := time.Now()
	//err = store.TXTConv(*result, "demo")
	//err = store.EPUBConv(*result, "demo")
	err = store.TXTConv(*result, "demo")
	fmt.Println(time.Since(start))
	fmt.Println(err)

	//多个并发处理
	//var url = make(map[int]string)
	//url[0] = "https://www.qb5.la/book_13659/"
	//url[1] = "https://www.qb5.la/book_37141/"
	//url[2] = "https://www.qb5.la/book_85367/"
	//保存txt
	//for i := 0; i < 3; i++ {
	//	go func(a int) {
	//		fmt.Println(a)
	//		result, err := site.BookInfo(url[a])
	//		site.Download(result)
	//		err = store.Conv(*result, "demo"+strconv.Itoa(a)+".txt")
	//		fmt.Println(err)
	//	}(i)
	//}
	//time.Sleep(1 * time.Minute)

}

////-------------- url资源变文件
//f, err := os.OpenFile("s.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
//fmt.Println(err)
//
//_, err = f.Write(bodyBytes)
//fmt.Println(err)
//
//f.Close()
////--------------

func TestTxt(t *testing.T) {
	//生成文件
	file2Path := "./an-xiao-yao.txt"

	fii, err := os.OpenFile(file2Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
		return
	}
	//读取文件列表
	part_list, err := filepath.Glob("./demo/an-xiao-yao/*.md")
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", part_list, file2Path)
	i := 0
	for _, v := range part_list {
		if v == "README.md" {
			continue
		}
		if v == "SUMMARY.md" {
			continue
		}
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fii.Write([]byte("\n\n" + string(b)))
		f.Close()
		i++
		fmt.Printf("合并%d个\n", i)
	}
	fii.Close()
	fmt.Println("合并成功")
}
