package video

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BaseURL var BaseURL = "https://www.kuaibozy.com/api.php/provide/vod/from/kbm3u8/at/xml"
var BaseURL = "https://api.apibdzy.com/api.php/provide/vod/at/xml/"

var client = &http.Client{Timeout: 10 * time.Second}

type Resources struct {
	List []struct {
		Id   int    `xml:"id"`
		Name string `xml:"name"`
		Pic  string `xml:"pic"`
		Type string `xml:"type"`
		Des  string `xml:"des"`
		Lang string `xml:"lang"`
		Last string `xml:"last"`
		Year string `xml:"year"`
		Dd   string `xml:"dl>dd"`
	} `xml:"list>video"`
}

// QueryResources 获取资源
func QueryResources(keyWord string, t string, p string, ids string) (ret *Resources, err error) {
	fullURL := fmt.Sprintf("%s?ac=videolist&wd=%s&t=%v&pg=%v&ids=%v", BaseURL, keyWord, t, p, ids)
	req, _ := http.NewRequest("GET", fullURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request err: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body err: %v\n", err)
	}

	v := &Resources{}
	err = xml.Unmarshal(body, v)
	if err != nil {
		return nil, fmt.Errorf("获取xml信息失败")
	}
	if len(v.List) == 0 {
		return nil, fmt.Errorf("没有查询到相关资源")
	}
	return v, nil
}

type Types struct {
	Class []struct {
		Value string `xml:",chardata"`
		Id    int    `xml:"id,attr"`
	} `xml:"class>ty"`
}

// QueryTypes 获取资源分类
func QueryTypes() (ret *Types, err error) {
	req, _ := http.NewRequest("GET", BaseURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request err: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body err: %v\n", err)
	}

	v := &Types{}
	err = xml.Unmarshal(body, v)
	if err != nil {
		return nil, fmt.Errorf("获取xml信息失败")
	}
	if len(v.Class) == 0 {
		return nil, fmt.Errorf("没有查询到相关资源")
	}
	return v, nil
}
