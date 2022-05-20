package book

import (
	"reflect"
	"runtime"

	"api/pkg/book/site"
	qb5la "api/pkg/book/site/qb5.la"
)

type siteFunc func() site.SiteA

func addSiteFunc(fn siteFunc) {
	s := fn()
	s.File = runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	site.AddSite(s)
}

func InitSites() {
	//addSiteFunc(demo.Site)
	addSiteFunc(qb5la.Site)
}
