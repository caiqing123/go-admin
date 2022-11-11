package book

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"

	"api/pkg/book/site"
	"api/pkg/book/site/cc_b520"
	"api/pkg/book/site/cn_bookstack"
	"api/pkg/book/site/me_zxcs"
	"api/pkg/book/site/qb5_la"
	"api/pkg/book/store"
	"api/pkg/file"
	"api/pkg/logger"
)

type siteFunc func() site.SiteA

func addSiteFunc(fn siteFunc) {
	s := fn()
	s.File = runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	site.AddSite(s)
}

func InitSites() {
	addSiteFunc(cc_b520.Site)
	addSiteFunc(me_zxcs.Site)
	addSiteFunc(qb5_la.Site)
	//addSiteFunc(org_wanben.Site)
	addSiteFunc(cn_bookstack.Site)
}

func Download(ctx context.Context, url string, id string, group string, hookfn func(context.Context, string, string, []byte)) {
	err := file.IsNotExistMkDir("public/uploads/book/")
	if err != nil {
		logger.Error("book dir error " + err.Error())
		return
	}
	result, err := site.BookInfo(url)
	if err != nil {
		logger.Warn(err.Error())
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(result.BookURL))
	md5s := hex.EncodeToString(md5Ctx.Sum(nil))

	site.DownloadWs(result, ctx, id, group, hookfn, "public/uploads/book/"+result.BookName+"_"+md5s+"_"+id)

	err = store.SourceConv(*result, "public/uploads/book/"+result.BookName+"_"+md5s)

	if strings.Contains(url, "bookstack.cn") {
		err = store.EPUBConv(*result, "public/uploads/book/"+result.BookName+"_"+md5s+"_"+id)
	} else {
		err = store.TXTConv(*result, "public/uploads/book/"+result.BookName+"_"+md5s+"_"+id)
	}
	if err != nil {
		logger.Warn(err.Error())
	}
}

func DownloadLog(ctx context.Context, log string, id string, group string, hookfn func(context.Context, string, string, []byte)) {
	files, err := ioutil.ReadDir("public/uploads/book/")
	logger.LogIf(err)
	bookList := make([]string, 0)
	for _, f := range files {
		//筛选
		if !strings.Contains(f.Name(), id) {
			continue
		}
		bookList = append(bookList, f.Name())
	}
	if len(bookList) == 0 {
		logger.Warn("book list is empty")
		return
	}
	src := `{"list":%v,"type":"book_list"}`

	b, _ := json.Marshal(bookList)
	go hookfn(ctx, id, group, []byte(fmt.Sprintf(src, string(b))))
}
