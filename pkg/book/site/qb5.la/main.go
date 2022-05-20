package qb5_la

import (
	"net/http"
	"net/url"
	"strings"

	"api/pkg/book/site"
	"api/pkg/book/utils"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "全本小说网",
		HomePage: "https://www.qb5.la/",
		Match: []string{
			`https://www\.qb5\.la/book_\d+/`,
			`https://www\.qb5\.la/book_\d+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//*[@id="info"]/h1/text()`,
			`//*[@id="picbox"]/div/img`,
			`//*[@id="info"]/h1/small/a/text()`,
			`//div[@class="zjbox"]/dl[@class="zjlist"]/dd/a`,
			``,
			`//div[@id='intro']`),
		Chapter: site.Type1Chapter(`//*[@id="content"]/text()`),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.qb5.la/modules/article/search.php")
				if err != nil {
					return nil
				}
				value := baseurl.Query()
				value.Add("searchtype", "all")
				value.Add("searchkey", utils.U8ToGBK(s))
				value.Add("sbt", utils.U8ToGBK("搜索"))
				baseurl.RawQuery = value.Encode()

				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					return nil
				}
				return req
			},
			`//tr`,
			`td[@class='odd'][1]/a`,
			`td[@class='odd'][2]`,
			func(r site.ChapterSearchResult) site.ChapterSearchResult {
				r.Author = strings.TrimPrefix(r.Author, "/")
				return r
			},
		),
	}
}
