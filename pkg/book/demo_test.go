package book

import (
	"fmt"
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
	//	result, err := s.Search("搜")
	//	fmt.Println(err)
	//	fmt.Println(result)
	//
	//}

	//详情
	result, err := site.BookInfo("http://www.b520.cc/12_12376/")
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
	//err = store.Conv(*result, "demo.txt")
	err = store.EPUBConv(*result, "demo.epub")
	fmt.Println(time.Since(start))
	fmt.Println(err)

	//多个并发处理
	//var url = make(map[int]string)
	//url[0] = "https://www.qb5.la/book_13659/"
	//url[1] = "https://www.qb5.la/book_37141/"
	//url[2] = "https://www.qb5.la/book_85367/"
	////保存txt 14s
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
