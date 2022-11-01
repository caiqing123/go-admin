package qqmusic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"api/pkg/music/comm"
)

// QQMusic qq音乐
func QQMusic(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	searchURL := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	param := `{"comm":{"format":"json","inCharset":"utf-8","outCharset":"utf-8","notice":0,"platform":"h5","needNewCode":1,"ct":23,"cv":0},"req_0":{"method":"DoSearchForQQMusicDesktop","module":"music.search.SearchCgiService","param":{"remoteplace":"txt.mqq.all","search_type":0,"query":"` + songName + `","page_num":` + p + `,"num_per_page":30}}}`
	param = url.QueryEscape(param)
	fullURL := searchURL + "?data=" + param
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
	req.Header.Add("Referer", "http://y.qq.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")
	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return nil, fmt.Errorf("请求错误")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return nil, fmt.Errorf("读取数据错误")
	}
	info := &SearchJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return nil, fmt.Errorf("获取json信息失败")
	}
	if info.Code == 0 {
		for index, result := range info.Req0.Data.Body.Song.List {
			downloadUrl, _ := getPlayURL(result.Mid)
			img := "https://y.gtimg.cn/music/photo_new/T002R300x300M000" + result.Album.Mid + ".jpg"
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Singer[0].Name + " ]", Author: result.Singer[0].Name,
				SongName: result.Name,
				SongURL:  downloadUrl,
				LrcData:  getLyric(result.Mid),
				ImgURL:   img,
				PicURL:   img,
			})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}

func QuickMusic(songName string) (ret []comm.Result, err error) {
	searchURL := "https://c.y.qq.com/splcloud/fcgi-bin/smartbox_new.fcg"
	fullURL := fmt.Sprintf("%s?key=%s&g_tk=5381", searchURL, songName)
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
	req.Header.Add("Referer", "http://y.qq.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")
	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return nil, fmt.Errorf("请求错误")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return nil, fmt.Errorf("读取数据错误")
	}
	info := &SearchQuickJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return nil, fmt.Errorf("获取json信息失败")
	}

	if info.Code == 0 {
		for index, result := range info.Data.Song.Itemlist {
			downloadUrl, _ := getPlayURL(result.Mid)
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Singer + " ]", Author: result.Singer,
				SongName: result.Name,
				SongURL:  downloadUrl})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}

func getPlayURL(mid string) (playURL string, err error) {
	if mid == "" {
		return "", nil
	}
	vkeyURL := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	param := `{"req": {"module": "CDN.SrfCdnDispatchServer", "method": "GetCdnDispatch", "param": {"guid": "3982823384", "calltype": 0, "userip": ""}}, "req_0": {"module": "vkey.GetVkeyServer", "method": "CgiGetVkey", "param": {"guid": "3982823384", "songmid": ["` + mid + `"], "songtype": [0], "uin": "925648047", "loginflag": 1, "platform": "20"}}, "comm": {"uin": "925648047", "format": "json", "ct": 24, "cv": 0}}`
	param = url.QueryEscape(param)

	fullURL := vkeyURL + "?data=" + param
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return
	}

	info := &VkeyJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return
	}
	if info.Code == 0 {
		for _, result := range info.Req0.Data.Midurlinfo {
			if purl := result.Purl; len(purl) > 0 {
				return info.Req0.Data.Sip[0] + purl, nil
			}
		}
	}
	return "", fmt.Errorf("not found")
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	u := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	param := `{"comm": {"ct": 24,"cv": 0},"detail": {"method": "GetDetail","module": "musicToplist.ToplistInfoServer","param": {"topId": 62,"offset": 0,"num":30,"period":""}}}`
	param = url.QueryEscape(param)
	fullURL := u + "?data=" + param
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
	req.Header.Add("Referer", "http://y.qq.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")

	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return nil, fmt.Errorf("请求错误")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return nil, fmt.Errorf("读取数据错误")
	}
	info := &CommendJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return nil, fmt.Errorf("获取json信息失败")
	}

	for index, result := range info.Detail.Data.SongInfoList {
		downloadUrl, _ := getPlayURL(result.Mid)
		img := "https://y.gtimg.cn/music/photo_new/T002R300x300M000" + result.Album.Mid + ".jpg"
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Title + " - [ " + result.Singer[0].Name + " ]", Author: result.Singer[0].Name,
			SongName: result.Title,
			SongURL:  downloadUrl,
			LrcData:  getLyric(result.Mid),
			ImgURL:   img,
			PicURL:   img,
		})
	}
	return ret, nil
}

func getNewPlayURL(mid string) (playURL string, err error) {
	if mid == "" {
		return "", nil
	}
	vkeyURL := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	//m4a C400 m4a |mp3 M500  128 |mp3 M800  320 |ape A000  ape |flac F000 flac
	//filename = M500 + mid + mid+ .mp3(文件格式)
	//param := `{"req_0":{"module":"vkey.GetVkeyServer","method":"CgiGetVkey","param":{"filename":["M800000EApX10oJiqD000EApX10oJiqD.mp3"],"guid":"guid","songmid":["000EApX10oJiqD"],"songtype":[0],"uin":"925648047","loginflag":1,"platform":"20"}},"loginUin":"925648047","comm":{"uin":"925648047","format":"json","ct":19,"cv":0,"authst":"Q_H_L_5GTtxvRC2IeYomCgXdpjZ2w-Vk98iUBca0v9sUWtTwAYyNmkRt_2uCw"}}`
	param := `{"req_0":{"module":"vkey.GetVkeyServer","method":"CgiGetVkey","param":{"filename":[],"guid":"guid","songmid":["` + mid + `"],"songtype":[0],"uin":"","loginflag":1,"platform":"20"}},"loginUin":"","comm":{"uin":"","format":"json","ct":19,"cv":0,"authst":""}}`
	param = url.QueryEscape(param)

	fullURL := vkeyURL + "?data=" + param
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return
	}

	info := &VkeyJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return
	}
	if info.Code == 0 {
		for _, result := range info.Req0.Data.Midurlinfo {
			if purl := result.Purl; len(purl) > 0 {
				return info.Req0.Data.Sip[0] + purl, nil
			}
		}
	}
	return "", fmt.Errorf("not found")
}

func getLyric(mid string) (lyric string) {
	if mid == "" {
		return ""
	}
	fullURL := fmt.Sprintf("https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg?songmid=%s&g_tk=5381&loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&platform=yqq", mid)
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Referer", "https://y.qq.com/portal/player.html")
	resp, err := comm.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	info := &LyricJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return ""
	}
	if info.Code == 0 {
		Lyric, _ := base64.StdEncoding.DecodeString(info.Lyric)
		return string(Lyric)
	}
	return ""
}
