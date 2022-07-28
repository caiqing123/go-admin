package hotlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	urls "net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"api/pkg/cache"
)

type HotData struct {
	Code    int
	Message string
	Data    interface{}
}

type Spider struct {
	DataType string
}

// GetV2EX V2EX
func (spider Spider) GetV2EX() []map[string]interface{} {
	url := "https://www.v2ex.com/?tab=hot"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	document.Find(".item_title").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://www.v2ex.com" + url})
		}
	})
	return allData
}

func (spider Spider) GetITHome() []map[string]interface{} {
	url := "https://www.ithome.com/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	document.Find("#rank #d-1 li").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": url})
		}
	})
	return allData
}

// GetZhiHu 知乎
func (spider Spider) GetZhiHu() []map[string]interface{} {
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	url := "https://www.zhihu.com/hot"
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1`)

	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".App-main div div").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.ParentFiltered("a").Attr("href")
		text := selection.Find("h1").Text()
		if boolUrl && text != "" {
			allData = append(allData, map[string]interface{}{"title": text, "url": url})
		}
	})
	return allData
}

func (spider Spider) GetWeiBo() []map[string]interface{} {
	url := "https://s.weibo.com/top/summary"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36`)
	request.Header.Add("cookie", `SUB=_2AkMVGMOLf8NxqwJRmfAXymziboh_ywvEieKjRDJQJRMxHRl-yT8Xqh0btRB6PpjtZItcsP16OHrUEpeyexCJTs118TOt;`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("tbody tr").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://s.weibo.com" + url})
		}
	})
	if len(allData) > 1 {
		return allData[1:]
	}
	return allData
}

// GetDouBan 豆瓣
func (spider Spider) GetDouBan() []map[string]interface{} {
	url := "https://www.douban.com/group/explore"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `PostmanRuntime/7.28.4`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Referer", `https://www.douban.com/group/explore`)
	request.Header.Add("Host", `www.douban.com`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".channel-item").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("h3 a").Attr("href")
		text := selection.Find("h3 a").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": url})
		}
	})
	return allData
}

// GetTianYa 天涯
func (spider Spider) GetTianYa() []map[string]interface{} {
	url := "http://bbs.tianya.cn/list.jsp?item=funinfo&grade=3&order=1"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Referer", `http://bbs.tianya.cn/list.jsp?item=funinfo&grade=3&order=1`)
	request.Header.Add("Host", `bbs.tianya.cn`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("table tr").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("td a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "http://bbs.tianya.cn/" + url})
		}
	})
	return allData
}

// GetHuPu 虎扑
func (spider Spider) GetHuPu() []map[string]interface{} {
	url := "https://bbs.hupu.com/all-gambia"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Referer", `https://bbs.hupu.com/`)
	request.Header.Add("Host", `bbs.hupu.com`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".text-list-model .list-item-wrap").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find(".t-info a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://bbs.hupu.com" + url})
		}
	})
	return allData
}

// GetGitHub Github
func (spider Spider) GetGitHub() []map[string]interface{} {
	url := "https://github.com/trending"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}

	document.Find(".Box article").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find(".lh-condensed a")
		//desc := selection.Find(".col-9 .text-gray .my-1 .pr-4")
		//descText := desc.Text()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		descText := selection.Find("p").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "desc": descText, "url": "https://github.com" + url})
		}
	})
	return allData
}

func (spider Spider) GetBaiDu() []map[string]interface{} {
	url := "http://top.baidu.com/buzz?b=341&c=513&fr=topbuzz_b1"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `top.baidu.com`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".SN-WEB-waterfall-item .row-start-center").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find(".one-line-ellipsis")
		text := s.Text()
		if text != "" {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://m.baidu.com/s?word=" + urls.QueryEscape(text) + "&sa=fyb_news"})
		}
	})
	return allData

}

func (spider Spider) Get36Kr() []map[string]interface{} {
	url := "https://36kr.com/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `36kr.com`)
	request.Header.Add("Referer", `https://36kr.com/`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".hotlist-item-toptwo").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := selection.Find("a p").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://36kr.com" + url})
		}
	})
	document.Find(".hotlist-item-other-info").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://36kr.com" + url})
		}
	})
	return allData

}

func (spider Spider) GetGuoKr() []map[string]interface{} {
	url := "https://www.guokr.com/scientific/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `www.guokr.com`)
	request.Header.Add("Referer", `https://www.guokr.com/scientific/`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("div .article").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h3 a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": text, "url": url})
			}
		}
	})
	return allData
}

func (spider Spider) GetHuXiu() []map[string]interface{} {
	url := "https://www.huxiu.com/article"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `www.huxiu.com`)
	request.Header.Add("Referer", `https://www.huxiu.com/channel/107.html`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".article-item-wrap").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text, _ := s.Find("img").Attr("alt")
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": text, "url": "https://www.huxiu.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetZHDaily() []map[string]interface{} {
	url := "http://daily.zhihu.com/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".row .box").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("span").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": text, "url": "https://daily.zhihu.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetSegmentfault() []map[string]interface{} {
	url := "https://segmentfault.com/hottest"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".news-list .news__item-info").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a:nth-child(2)").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("h4").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": text, "url": "https://segmentfault.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetHacPai() []map[string]interface{} {
	url := "https://hacpai.com/domain/play"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".hotkey li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h2 a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": text, "url": url})
			}
		}
	})
	return allData
}

func (spider Spider) GetWYNews() []map[string]interface{} {
	url := "https://gw.m.163.com/gentie-web/api/v2/products/a2869674571f77b5a0867c3d71db5856/rankDocs/all/list?ibc=newsapph5&limit=30"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var allData []map[string]interface{}

	type JSONStruct struct {
		Data struct {
			CmtDocs []struct {
				DocTitle string `json:"doc_title"`
				DocId    string `json:"docId"`
			} `json:"cmtDocs"`
		} `json:"data"`
		Code int `json:"code"`
	}

	info := &JSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("获取json信息失败" + spider.DataType)
		return []map[string]interface{}{}
	}
	if info.Code == 0 && len(info.Data.CmtDocs) != 0 {
		for _, result := range info.Data.CmtDocs {
			allData = append(allData, map[string]interface{}{"title": result.DocTitle, "url": "https://c.m.163.com/news/a/" + result.DocId + ".html"})
		}
	}
	return allData
}

func (spider Spider) GetWaterAndWood() []map[string]interface{} {
	url := "https://www.newsmth.net/nForum/mainpage?ajax"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//sss,_ := GbkToUtf8([]byte(string(str)))
	//fmt.Println(string(sss))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	// topics
	document.Find("#top10 li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a:nth-child(2)").First()
		url, boolUrl := s.Attr("href")
		text, _ := GbkToUtf8([]byte(s.Text()))
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.newsmth.net" + url})
				}
			}
		}
	})
	document.Find(".topics").Find("li").Each(func(i int, selection *goquery.Selection) {
		if i > 10 {
			s := selection.Find("a:nth-child(2)").First()
			url, boolUrl := s.Attr("href")
			text, _ := GbkToUtf8([]byte(s.Text()))
			if len(text) != 0 {
				if boolUrl {
					if len(allData) <= 100 {
						allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.newsmth.net" + url})
					}
				}
			}
		}
	})
	return allData
}

// http://nga.cn/

func (spider Spider) GetNGA() []map[string]interface{} {
	url := "http://nga.cn/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("h2").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": url})
				}
			}
		}
	})
	return allData
}

func (spider Spider) GetCSDN() []map[string]interface{} {
	url := "https://www.csdn.net/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".headswiper-content .headswiper-item").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": url})
				}
			}
		}
	})
	return allData
}

// GetWeiXin https://weixin.sogou.com/?pid=sogou-wsse-721e049e9903c3a7&kw=
func (spider Spider) GetWeiXin() []map[string]interface{} {
	url := "https://weixin.sogou.com/?pid=sogou-wsse-721e049e9903c3a7&kw="
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".news-list li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h3 a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": url})
				}
			}
		}
	})
	return allData
}

// GetKD 凯迪
func (spider Spider) GetKD() []map[string]interface{} {
	url := "https://www.9kd.com/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".hot-news-list li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": url})
				}
			}
		}
	})
	return allData
}

// https://www.chiphell.com/

func (spider Spider) GetChiphell() []map[string]interface{} {
	url := "https://www.chiphell.com/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("#frameP3M3F6 li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	return allData
}

// http://jandan.net/

func (spider Spider) GetJianDan() []map[string]interface{} {
	url := "http://jandan.net/"
	timeout := 5 * time.Second //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("h2").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": text, "url": url})
				}
			}
		}
	})
	return allData
}

//GbkToUtf8 部分热榜标题需要转码
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//ExecGetData 执行每个分类数据
func ExecGetData(spider Spider) {
	var count = 0
LOOP:
	if count < 3 {
		start := time.Now()
		reflectValue := reflect.ValueOf(spider)
		dataType := reflectValue.MethodByName("Get" + spider.DataType)
		data := dataType.Call(nil)
		originData := data[0].Interface().([]map[string]interface{})
		seconds := time.Since(start).Seconds()
		if len(originData) > 0 {
			group.Done()
			cache.Caches(spider.DataType, originData, time.Minute*30)
			fmt.Printf("耗费 %.2fs 秒完成抓取%s", seconds, spider.DataType)
			fmt.Println()
		} else {
			fmt.Printf("耗费 %.2fs 秒,抓取失败%s", seconds, spider.DataType)
			fmt.Println()
			//失败重试2次
			time.Sleep(time.Duration(count+1) * 5 * time.Second)
			count = count + 1
			goto LOOP
		}
	} else {
		group.Done()
	}
}

var group sync.WaitGroup

func All(isUpdate bool) {
	allData := []string{
		"ZhiHu",
		"WeiBo",
		"DouBan",
		"TianYa",
		"HuPu",
		"BaiDu",
		"36Kr",
		"GuoKr",
		"HuXiu",
		"ZHDaily",
		"Segmentfault",
		"WYNews",
		"WaterAndWood",
		"HacPai",
		"KD",
		"NGA",
		"WeiXin",
		"Chiphell",
		"JianDan",
		"ITHome",
		"CSDN",
	}
	fmt.Println("开始抓取" + strconv.Itoa(len(allData)) + "种数据类型")
	group.Add(len(allData))
	var spider Spider
	for _, value := range allData {
		if isUpdate == true {
			cache.Forget(value)
		}
		if cache.Get(value) != "" {
			group.Done()
			continue
		}
		spider = Spider{DataType: value}
		go ExecGetData(spider)
	}
	group.Wait()
}
