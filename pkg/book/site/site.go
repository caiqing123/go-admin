package site

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/transform"

	"api/pkg/book/store"
	"api/pkg/book/utils"
)

var SitePool []SiteA

type SiteA struct {
	Name     string // 站点名称
	HomePage string // 站点首页

	File string

	// match url, look that https://godoc.org/path#Match
	Match []string

	// search book on site
	Search func(s string) (result []ChapterSearchResult, err error) `json:"-"`

	// parse fiction info by page body
	BookInfo func(body io.Reader) (s *store.Store, err error) `json:"-"`

	// parse fiction chaper content by page body
	Chapter func(context.Context) (content []string, err error) `json:"-"`
}

type ChapterSearchResult struct {
	BookName string
	Author   string
	BookURL  string
}

func AddSite(site SiteA) {
	if site.File == "" {
		_, filename, _, _ := runtime.Caller(1)
		site.File = filename
	}
	SitePool = append(SitePool, site)
}

func (s SiteA) match(u string) (bool, error) {
	for _, v := range s.Match {
		re, err := regexp.Compile(v)
		if err != nil {
			return false, err
		}
		if re.MatchString(u) {
			return true, nil
		}
	}
	return false, nil
}

func MatchSites(pool []SiteA, u string) (*SiteA, error) {
	var result []*SiteA
	for k := range pool {
		ok, err := pool[k].match(u)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, &pool[k])
		}
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// BookInfo 获取小说信息
func BookInfo(BookURL string) (s *store.Store, err error) {
	ms, err := MatchSites(SitePool, BookURL)

	if err != nil {
		return nil, err
	}

	// Get WebPage
	resp, err := utils.RequestGet(BookURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if ms.BookInfo == nil {
		return nil, nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var body io.Reader = bytes.NewReader(bodyBytes)

	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		encode, err := utils.DetectContentCharset(bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		body = transform.NewReader(body, encode.NewDecoder())
	}
	chapter, err := ms.BookInfo(body)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(chapter.BookName) == "" {
		err = fmt.Errorf("BookInfo Name is empty")
		return
	}

	chapter.Author = strings.Replace(chapter.Author, "\u00a0", "", -1)

	chapter.BookURL = BookURL

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u1, err := resp.Request.URL.Parse(k2.URL)
			if err != nil {
				return nil, err
			}
			chapter.Volumes[v1].Chapters[v2].URL = u1.String()
		}
	}

	CoverURL, err := url.Parse(chapter.CoverURL)
	if err != nil {
		return nil, err
	}

	if chapter.CoverURL != "" {
		chapter.CoverURL = resp.Request.URL.ResolveReference(CoverURL).String()
	}

	if chapter.DownloadURL != "" {
		DownloadURL, _ := url.Parse(chapter.DownloadURL)
		chapter.DownloadURL = resp.Request.URL.ResolveReference(DownloadURL).String()
	}

	if len(chapter.Volumes) == 0 && chapter.DownloadURL == "" {
		return nil, fmt.Errorf("not match volumes")
	}
	return chapter, err
}

// Chapter 获取章节内容
func Chapter(BookURL string) (content []string, err error) {
	ms, err := MatchSites(SitePool, BookURL)
	if err != nil {
		return nil, err
	}
	// Get WebPage
	body, err := utils.GetWebPageBodyReader(BookURL)
	if err != nil {
		return nil, err
	}

	if ms.Chapter == nil {
		return nil, fmt.Errorf("site %s Chapter Func is empty", ms.Name)
	}

	bu, err := url.Parse(BookURL)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "url", bu)
	ctx = context.WithValue(ctx, "body", body)

	return ms.Chapter(ctx)
}

func Type1SearchAfter(
	getReq func(s string) *http.Request,
	resultExpr, nameExpr, authorExpr string,
	after func(r ChapterSearchResult) ChapterSearchResult) func(s string) (result []ChapterSearchResult, err error) {
	return func(s string) (result []ChapterSearchResult, err error) {
		req := getReq(s)
		if req == nil {
			return
		}
		var (
			resp *http.Response
		)
		if err = utils.Retry(5, time.Millisecond*500, func() error {
			resp, err = http.DefaultClient.Do(req)
			return err
		}); err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if req.URL.String() != resp.Request.URL.String() {
			// 单个搜索结果 直接到详情页
			bookInfo, e := BookInfo(resp.Request.URL.String())
			if e != nil {
				return nil, e
			}
			r := ChapterSearchResult{
				BookName: bookInfo.BookName,
				Author:   bookInfo.Author,
				BookURL:  resp.Request.URL.String(),
			}
			if after != nil {
				r = after(r)
			}
			result = append(result, r)
			return result, nil
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return
		}

		var body io.Reader = bytes.NewReader(bodyBytes)
		encode, err := utils.DetectContentCharset(bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		body = transform.NewReader(body, encode.NewDecoder())

		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}

		r := htmlquery.Find(doc, resultExpr)
		if len(r) == 0 {
			return nil, nil
		}
		for _, v := range r {
			s2 := htmlquery.FindOne(v, nameExpr)
			if s2 == nil {
				continue
			}

			s4 := &html.Node{}
			if authorExpr != "" {
				s4 = htmlquery.FindOne(v, authorExpr)
				if s4 == nil {
					continue
				}
			}

			u1, _ := url.Parse(htmlquery.SelectAttr(s2, "href"))
			r := ChapterSearchResult{
				BookName: htmlquery.InnerText(s2),
				Author:   htmlquery.InnerText(s4),
				BookURL:  resp.Request.URL.ResolveReference(u1).String(),
			}
			if after != nil {
				r = after(r)
			}
			result = append(result, r)
		}
		return result, nil
	}
}

func Type1BookInfo(nameExpr, coverExpr, authorExpr, chapterExpr, DownloadExpr, DescriptionExpr string, after func(r *store.Store) *store.Store) func(body io.Reader) (s *store.Store, err error) {
	return func(body io.Reader) (s *store.Store, err error) {
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}
		s = &store.Store{}
		var tmpNode *html.Node

		tmpNode = htmlquery.FindOne(doc, nameExpr)
		if tmpNode == nil {
			err = fmt.Errorf("no matching bookName")
			return
		}
		s.BookName = htmlquery.InnerText(tmpNode)

		if coverExpr != "" {
			coverNode := htmlquery.FindOne(doc, coverExpr)
			if coverNode == nil {
				err = fmt.Errorf("no matching cover")
				return
			}
			if cu, err := url.Parse(strings.TrimSpace(htmlquery.SelectAttr(coverNode, "src"))); err != nil {
				log.Printf("Cover Image URL Error:" + err.Error())
			} else {
				s.CoverURL = cu.String()
			}
		}

		// Author
		if authorExpr != "" {
			authorContent := htmlquery.FindOne(doc, authorExpr)
			if authorContent == nil {
				err = fmt.Errorf("no matching author")
				return
			}
			s.Author = strings.TrimSpace(htmlquery.InnerText(authorContent))
		}

		if DownloadExpr != "" {
			// DownloadExpr
			downloadContent := htmlquery.FindOne(doc, DownloadExpr)
			if downloadContent == nil {
				err = fmt.Errorf("no matching downloadContent")
				return
			}
			u1, _ := url.Parse(htmlquery.SelectAttr(downloadContent, "href"))
			s.DownloadURL = u1.String()
		}

		if DescriptionExpr != "" {
			// DescriptionExpr
			descriptionContent := htmlquery.FindOne(doc, DescriptionExpr)
			if descriptionContent == nil {
				err = fmt.Errorf("no matching DescriptionExpr")
				return
			}
			s.Description = htmlquery.InnerText(descriptionContent)
		}

		// Contents
		if chapterExpr != "" {
			nodeContent := htmlquery.Find(doc, chapterExpr)
			if len(nodeContent) == 0 && s.DownloadURL == "" {
				err = fmt.Errorf("no matching contents")
				return
			}

			var vol = store.Volume{
				Name:     "正文",
				Chapters: make([]store.Chapter, 0),
			}
			for _, v := range nodeContent {
				chapterURL, err := url.Parse(htmlquery.SelectAttr(v, "href"))
				if err != nil {
					return nil, err
				}
				vol.Chapters = append(vol.Chapters, store.Chapter{
					Name: strings.TrimSpace(htmlquery.InnerText(v)),
					URL:  chapterURL.String(),
				})
			}
			s.Volumes = append(s.Volumes, vol)
		}
		if after != nil {
			s = after(s)
		}

		return
	}
}

// Type1Chapter 章节段落匹配
func Type1Chapter(expr string) func(ctx context.Context) (content []string, err error) {
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
			t := htmlquery.InnerText(v)
			t = strings.TrimSpace(t)

			if t == "" {
				continue
			}

			M = append(M, t)
		}
		return M, nil
	}
}

// Type2Chapter 章节匹配2：单章分多页,
// next函数返回下一个页面的DOM
// block函数用于屏蔽多余的段落
func Type2Chapter(
	expr string,
	next func(preURL *url.URL, doc *html.Node) *html.Node,
	block func([]string) []string,
) func(context.Context) (content []string, err error) {
	return func(ctx context.Context) (content []string, err error) {
		doc, err := htmlquery.Parse(ctx.Value("body").(io.Reader))
		if err != nil {
			return nil, err
		}
		var M []string
		if block == nil {
			block = func(a []string) []string { return a }
		}
		for {
			//list
			nodeContent := htmlquery.Find(doc, expr)
			if len(nodeContent) == 0 {
				err = fmt.Errorf("no matching content")
				return nil, err
			}
			var MM []string
		loopContent:
			for _, v := range nodeContent {
				t := htmlquery.InnerText(v)
				t = strings.TrimSpace(t)

				if t == "" {
					continue loopContent
				}
				MM = append(MM, t)
			}
			M = append(M, block(MM)...)

			if next == nil {
				return M, nil
			}
			doc = next(ctx.Value("url").(*url.URL), doc)
			if doc == nil {
				break
			}
		}
		return M, nil
	}
}
