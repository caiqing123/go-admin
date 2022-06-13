package zinc

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var zinc = ZincClient{
	ZincClientConfig: &ZincClientConfig{
		ZincHost:     "http://localhost:4080",
		ZincUser:     "admin",
		ZincPassword: "admin",
	},
}

func TestZincClient_EsQuery(t *testing.T) {
	queryMap := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]string{}, //搜索全部
		},
		//"query": map[string]interface{}{
		//	"match_phrase": map[string]interface{}{
		//		"business_type": "api", //条件搜索 business_type搜索字段
		//	},
		//},
		"sort": []string{"-is_top", "-latest_replied_on"},
		"from": 0,
		"size": 10,
	}
	rsp, err := zinc.EsQuery("admin-log", queryMap)
	prtBodyBytes, _ := json.MarshalIndent(rsp, "", "   ")

	fmt.Println(string(prtBodyBytes))
	fmt.Println(err)
}

func TestZincClient_ApiQuery(t *testing.T) {
	query := "error"
	queryMap := map[string]interface{}{
		"search_type": "querystring",
		"query": map[string]interface{}{
			"term": query,
		},
		"sort_fields": []string{"-is_top", "-latest_replied_on"},
		"from":        0,
		"max_results": 10,
	}
	rsp, err := zinc.ApiQuery("paopao-log", queryMap)
	prtBodyBytes, _ := json.MarshalIndent(rsp, "", "   ")

	fmt.Println(string(prtBodyBytes))
	fmt.Println(err)
}

func TestZincClient_ExistIndex(t *testing.T) {
	exist := zinc.ExistIndex("paopao-log")
	fmt.Println(exist)
}

func TestZincClient_CreateIndex(t *testing.T) {
	exist := zinc.CreateIndex("demo", &ZincIndexProperty{
		"id": &ZincIndexPropertyT{
			Type:     "numeric",
			Index:    true,
			Store:    true,
			Sortable: true,
		},
		"user_id": &ZincIndexPropertyT{
			Type:  "numeric",
			Index: true,
			Store: true,
		},
	})
	fmt.Println(exist)
}

func TestZincClient_PutDoc(t *testing.T) {
	doc := map[string]interface{}{
		"id":                2,
		"user_id":           20,
		"comment_count":     1,
		"collection_count":  11,
		"upvote_count":      12,
		"is_top":            1,
		"is_essence":        1,
		"content":           "contentFormated",
		"tags":              "tagMaps",
		"ip_loc":            1,
		"latest_replied_on": 22,
		"attachment_price":  25,
		"created_on":        30,
		"modified_on":       31,
	}
	exist, err := zinc.PutDoc("paopao-data", 2, doc)
	fmt.Println(err)
	fmt.Println(exist)
}

func TestZincClient_BulkPushDoc(t *testing.T) {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"index": map[string]interface{}{
			"_index": "paopao-data",
			"_id":    fmt.Sprintf("%d", 1),
		},
	}, map[string]interface{}{
		"id":                1,
		"user_id":           20,
		"comment_count":     1,
		"collection_count":  11,
		"upvote_count":      12,
		"is_top":            1,
		"is_essence":        1,
		"content":           "contentFormated",
		"tags":              "tagMaps",
		"ip_loc":            1,
		"latest_replied_on": 22,
		"attachment_price":  25,
		"created_on":        30,
		"modified_on":       31,
	})
	exist, err := zinc.BulkPushDoc(data)
	fmt.Println(err)
	fmt.Println(exist)
}

func TestZincClient_BulkPutLogDoc(t *testing.T) {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"index": map[string]interface{}{
			"_index": "paopao-log",
		},
	}, map[string]interface{}{
		"time":    time.Now(),
		"message": "messagemessagemessage",
		"level":   "info",
		"data": map[string]interface{}{
			"name": "lin",
		},
	})
	exist, err := zinc.BulkPutLogDoc(data)
	fmt.Println(err)
	fmt.Println(exist)
}

func TestZincClient_DelDoc(t *testing.T) {
	err := zinc.DelDoc("paopao-data", "1")
	fmt.Println(err)
}

//代理请求测试
func TestProxy(t *testing.T) {
	//代理地址
	proxyStr := "http://165.232.119.23:80"
	proxy, err := url.Parse(proxyStr)
	// 目标网页
	pageUrl := "https://myip.ipip.net/"
	//  请求目标网页
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	req, _ := http.NewRequest("GET", pageUrl, nil)
	//req.Header.Add("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	res, err := client.Do(req)

	if err != nil {
		// 请求发生异常
		fmt.Println(err.Error())
	} else {
		defer res.Body.Close() //保证最后关闭Body

		fmt.Println("status code:", res.StatusCode) // 获取状态码

		// 有gzip压缩时,需要解压缩读取返回内容
		if res.Header.Get("Content-Encoding") == "gzip" {
			reader, _ := gzip.NewReader(res.Body) // gzip解压缩
			defer func(reader *gzip.Reader) {
				err := reader.Close()
				if err != nil {
				}
			}(reader)
			_, _ = io.Copy(os.Stdout, reader)
		}

		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	}
}
