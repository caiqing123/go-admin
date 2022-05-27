package book

import (
	"context"
	"reflect"
	"runtime"

	"api/pkg/book/site"
	"api/pkg/book/site/cc_b520"
	"api/pkg/book/site/me_zxcs"
	"api/pkg/book/site/org_wanben"
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
	addSiteFunc(org_wanben.Site)
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
	site.DownloadWs(result, ctx, id, group, hookfn)
	err = store.Conv(*result, "public/uploads/book/"+result.BookName+"_"+id+".txt")
	if err != nil {
		logger.Warn(err.Error())
	}
}
