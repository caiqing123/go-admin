package cc_b520

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"

	"api/pkg/book/site"
	"api/pkg/book/store"
	"api/pkg/book/utils"
	"api/pkg/logger"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "笔趣阁",
		HomePage: "https://www.b520.cc/",
		Match: []string{
			`http://www\.b520\.cc/\d+_\d+/*`,
			`http://www\.b520\.cc/\d+_\d+/\d+\.html/*`,
		},
		BookInfo: site.Type1BookInfo(
			`//*[@id="info"]/h1/text()`,
			`//*[@id="fmimg"]/img`,
			`//*[@id="info"]/p[1]/text()`,
			`//div[@id="list"]/dl/dd/a`,
			``,
			`//div[@id='intro']`,
			func(r *store.Store) *store.Store {
				//处理函数 可为nil
				r.Author = strings.TrimLeft(r.Author, "作\u00a0\u00a0\u00a0\u00a0者：")
				r.Volumes[0].Chapters = r.Volumes[0].Chapters[9:]

				baseurl, _ := url.Parse("http://downnovel.com/search.htm")
				value := baseurl.Query()
				value.Add("keyword", r.BookName)
				baseurl.RawQuery = value.Encode()
				doc, err := utils.GetWegPageDOM(baseurl.String())
				if err != nil {
					logger.Error(err.Error())
				} else {
					downpour := ""
					downloadContent := htmlquery.FindOne(doc, `//li[1]/a`)
					if downloadContent != nil {
						u1, _ := url.Parse(htmlquery.SelectAttr(downloadContent, "href"))
						downpour = u1.String()
					}
					doc1, err1 := utils.GetWegPageDOM("http://downnovel.com" + downpour)
					if err1 == nil {
						downBook := ""
						downloadContent1 := htmlquery.FindOne(doc1, `//a[@class="btn_b"]`)
						if downloadContent1 != nil {
							u1, _ := url.Parse(htmlquery.SelectAttr(downloadContent1, "href"))
							downBook = u1.String()
						}
						r.DownloadURL = downBook
					}
				}
				return r
			},
		),
		Chapter: site.Type1Chapter(`//*[@id="content"]/p`),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("http://www.b520.cc/modules/article/search.php")
				if err != nil {
					return nil
				}
				value := baseurl.Query()
				value.Add("searchkey", s)
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
