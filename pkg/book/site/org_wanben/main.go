package org_wanben

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"api/pkg/book/site"
	"api/pkg/book/store"
	"api/pkg/book/utils"
	"api/pkg/logger"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "完本神站",
		HomePage: "https://www.wanben.org/",
		Match: []string{
			`https://www\.wanben\.org/\d+/`,
			`https://www\.wanben\.org/\d+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//div[@class="detailTitle"]/h1/text()`,
			`//div[@class="detailTopLeft"]/img/@src`,
			`//div[@class="detailTopMid"]/div[@class="writer"]/a/text()`,
			`//div[@class="chapter"]/ul/li/a`,
			``,
			`//td[@colspan="7"][2]/text()`,
			nil,
		),
		Search: site.Type1SearchAfter(
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.wanben.org/modules/article/ss20210414.php")
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				captcha := ""
				client := &http.Client{}
				req1, err := http.NewRequest("GET", "https://www.wanben.org/captcha.php", nil)
				if err != nil {
				}
				resp, err := client.Do(req1)
				if err != nil {
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
					}
				}(resp.Body)
				result, _ := ioutil.ReadAll(resp.Body)

				res := store.GeneralData{}
				if resp, err := store.Request(store.GeneralMethod, result); err != nil {
					fmt.Println(err.Error())
					return nil
				} else if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
					fmt.Println(err.Error())
					return nil
				} else if res.Error() {
					fmt.Println(res.Response.Error.Message)
					return nil
				} else {
					for _, sdr := range res.Response.TextDetections {
						captcha += sdr.DetectedText
					}
				}
				if len(captcha) != 4 {
					fmt.Println("验证码识别错误" + (captcha))
					return nil
				}
				fmt.Println("原验证码" + (captcha))
				fmt.Println("验证码" + strings.ToLower(captcha))

				value := baseurl.Query()
				value.Add("searchkey", utils.U8ToGBK(s))
				value.Add("captcha", strings.ToLower(captcha))
				baseurl.RawQuery = value.Encode()
				req, err := http.NewRequest("GET", baseurl.String(), nil)

				for _, co := range resp.Cookies() {
					if co.Name == "PHPSESSID" {
						req.Header.Add("Cookie", "PHPSESSID="+co.Value)
					}
				}
				if err != nil {
					logger.Error(err.Error())
					return nil
				}
				return req
			},
			`//div[@class='resultLeft']/ul/li`,
			`//div[@class='sortPhr']/a`,
			`//p[@class='author']/a`,
			func(r site.ChapterSearchResult) site.ChapterSearchResult {
				//处理函数 可为nil
				return r
			},
		),
		Chapter: site.Type1Chapter(`//div[@class="readerCon"]/p/text()`),

		//Chapter: site.Type2Chapter(`//div[@class="readerCon"]/p/text()`, func(preURL *url.URL, doc *html.Node) *html.Node {
		//	nextNode := htmlquery.FindOne(doc, `//div[@class="readPage"]/a[3]`)
		//	if nextNode == nil {
		//		return nil
		//	}
		//	nextText := htmlquery.InnerText(nextNode)
		//	// log.Printf("nextText: %v\n", nextText)
		//	fmt.Println(1)
		//	if strings.Contains(nextText, "下一章") {
		//		return nil
		//	} else if strings.Contains(nextText, "下一页") {
		//		nextURL := htmlquery.SelectAttr(nextNode, "href")
		//		// log.Printf("nextURL: %v\n", nextURL)
		//		doc, err := utils.GetWegPageDOM(nextURL)
		//		if err != nil {
		//			log.Printf("GetWegPageDOM: %s", err)
		//			return nil
		//		}
		//		return doc
		//	}
		//	return nil
		//}, func(b []string) []string {
		//	if strings.HasPrefix(b[0], "一秒记住") {
		//		b = b[1:]
		//	}
		//	return b[:len(b)-1]
		//}),
	}
}
