package me_zxcs

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
		Name:     "知轩藏书",
		HomePage: "http://zxcs.me",
		Match: []string{
			`http://www\.zxcs\.me/\w+/\d+`,
			`http://zxcs\.me/\w+/\d+`,
		},
		BookInfo: site.Type1BookInfo(
			`//div[@id="content"]/h1`,
			`//div[@id="content"]//img`,
			`//div[@id="content"]/h1`,
			``,
			`//div[@id="content"]/div[@class="pagefujian"]/div[@class="down_2"]/a`,
			`//div[@id="content"]/p[3]`,
			func(r *store.Store) *store.Store {
				//处理函数 可为nil
				name := strings.Split(r.BookName, "作者：")
				r.BookName = name[0]
				r.Author = name[1]
				description := strings.Split(r.Description, "【内容简介】：")
				description = strings.Split(description[1], "【作者简介】：")
				r.Description = description[0]

				doc, err := utils.GetWegPageDOM(r.DownloadURL)
				if err != nil {
					logger.Error(err.Error())
				} else {
					downpour := ""
					downloadContent := htmlquery.FindOne(doc, `//span[@class="downfile"][1]/a`)
					if downloadContent != nil {
						u1, _ := url.Parse(htmlquery.SelectAttr(downloadContent, "href"))
						downpour = u1.String()
					}
					r.DownloadURL = downpour
				}
				return r
			},
		),
		Chapter: site.Type1Chapter(``),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("http://zxcs.me/index.php")
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				value := baseurl.Query()
				value.Add("keyword", s)
				baseurl.RawQuery = value.Encode()
				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				return req
			},
			`//dl[@id='plist']`,
			`//dt/a`,
			`//dt/a`,
			func(r site.ChapterSearchResult) site.ChapterSearchResult {
				//处理函数 可为nil
				name := strings.Split(r.BookName, "作者：")
				r.BookName = name[0]
				r.Author = name[1]
				return r
			},
		),
	}
}
