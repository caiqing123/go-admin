package demo

import (
	"net/http"
	"net/url"

	"api/pkg/book/site"
	"api/pkg/book/utils"
	"api/pkg/logger"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "demo",
		HomePage: "https://segmentfault.com/",
		Match: []string{
			`https://www\.segmentfault\.com/\w+/\d+`,
			`https://segmentfault\.com/\w+/\d+`,
		},
		BookInfo: site.Type1BookInfo(
			`//h1`,
			`//*[@class="d-flex align-items-center flex-wrap"]/a/picture/img`,
			`//*[@class="d-flex align-items-center flex-wrap"]/a/strong`,
			`//footer[@id="footer"]/div/div/dl/dd/a`,
			``,
			``),
		Chapter: site.Type1Chapter(`//article`),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://segmentfault.com/search")
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				value := baseurl.Query()
				value.Add("q", utils.U8ToGBK(s))
				baseurl.RawQuery = value.Encode()
				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				return req
			},
			`//div[@class='list-group list-group-flush']/li`,
			`//a`,
			`//div[@class='text-secondary text-truncate-1 font-size-14 mb-2']`,
			func(r site.ChapterSearchResult) site.ChapterSearchResult {
				//处理函数 可为nil
				return r
			},
		),
	}
}
