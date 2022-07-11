package cn_bookstack

import (
	"context"
	"fmt"
	"io"
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
		Name:     "书栈网",
		HomePage: "https://www.bookstack.cn/",
		Match: []string{
			`https://www\.bookstack\.cn/\w+/\w+`,
			`https://bookstack\.cn/\w+/\w+`,
		},
		BookInfo: site.Type1BookInfo(
			`//h1`,
			`//div[@class="row bookstack-info"]/div[2]/div/img`,
			``,
			`//div[@class="help-block"]/ul[@class="none-listyle"]/li/a`,
			``,
			`//li[@class="bookstack-description"]/div`,
			func(r *store.Store) *store.Store {
				doc, err := utils.GetWegPageDOM("https://www.bookstack.cn" + r.Volumes[0].Chapters[0].URL)
				if err != nil {
					logger.Error(err.Error())
				} else {
					r.Volumes = nil
					nodeContent := htmlquery.Find(doc, `//div[@class="article-menu-detail collapse-menu"]//ul//li/a`)
					var vol = store.Volume{
						Name:     "正文",
						Chapters: make([]store.Chapter, 0),
					}
					for _, v := range nodeContent {
						chapterURL, err := url.Parse(htmlquery.SelectAttr(v, "href"))
						if err != nil {
							logger.Error(err.Error())
						}
						cq := store.Chapter{
							Name: strings.TrimSpace(htmlquery.InnerText(v)),
							URL:  chapterURL.String(),
						}
						vol.Chapters = append(vol.Chapters, cq)
					}
					//去重
					vol.Chapters = removeDuplicateElement(vol)
					r.Volumes = append(r.Volumes, vol)
				}
				return r
			},
		),
		Chapter: Chapter(`//article[@id="page-content"]`),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.bookstack.cn/search/result")
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				value := baseurl.Query()
				value.Add("wd", s)
				baseurl.RawQuery = value.Encode()
				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				return req
			},
			`//ul/li[@class='clearfix']/div[2]`,
			`//a`,
			``,
			func(r site.ChapterSearchResult) site.ChapterSearchResult {
				//处理函数
				return r
			},
		),
	}
}

//removeDuplicateElement 去重处理
func removeDuplicateElement(languages store.Volume) []store.Chapter {
	var result []store.Chapter
	temp := map[string]struct{}{}
	for _, item := range languages.Chapters {
		if _, ok := temp[item.URL]; !ok {
			temp[item.URL] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Chapter(expr string) func(ctx context.Context) (content []string, err error) {
	return func(ctx context.Context) (content []string, err error) {
		if expr == "" {
			return
		}
		doc, err := htmlquery.Parse(ctx.Value("body").(io.Reader))
		if err != nil {
			return nil, err
		}

		var M []string
		//list
		nodeContent := htmlquery.Find(doc, expr)
		if len(nodeContent) == 0 {
			err = fmt.Errorf("no matching content")
			return nil, err
		}
		for _, v := range nodeContent {
			//删除元素
			fi := htmlquery.FindOne(v, "//blockquote")
			if fi != nil {
				v.RemoveChild(fi)
			}
			h1 := htmlquery.FindOne(v, "//h1")
			if h1 != nil {
				v.RemoveChild(h1)
			}
			toc := htmlquery.FindOne(v, `//div[@class="markdown-toc editormd-markdown-toc"]`)
			if toc != nil {
				v.RemoveChild(toc)
			}

			t := htmlquery.OutputHTML(v, false)
			//修改元素
			img := htmlquery.Find(v, "//img")
			if len(img) != 0 {
				for _, v1 := range img {
					imgUrl := htmlquery.SelectAttr(v1, "data-original")
					t = strings.Replace(t, `<img src="/static/images/loading.gif" alt="Alt text" data-original="`+imgUrl+`"`, `<img src="`+imgUrl+`" alt="Alt text" data-original="`+imgUrl+`"`, -1)
				}
			}

			if t == "" {
				continue
			}
			//错误内容 重试
			if strings.Contains(t, "BookStack") {
				err = fmt.Errorf("错误内容")
				return nil, err
			}
			M = append(M, t)
		}
		return M, nil
	}
}
