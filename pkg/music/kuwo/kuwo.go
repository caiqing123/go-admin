package kuwo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	nhttp "api/pkg/http"
	"api/pkg/music/comm"
)

// Kuwo kuwo 音乐
func Kuwo(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	u, _ := url.ParseRequestURI("http://kuwo.cn/api/www/search/searchMusicBykeyWord")
	query := url.Values{}
	query.Set("key", songName)
	query.Set("pn", p)
	query.Set("rn", "30")
	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Host", "kuwo.cn")
	req.Header.Add("Referer", "http://kuwo.cn/")
	req.Header.Add("csrf", "4IT871VN3DA")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")
	req.Header.Add("Cookie", "kw_token=4IT871VN3DA")

	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return nil, fmt.Errorf("酷我网络请求失败")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return nil, fmt.Errorf("酷我获取数据失败")
	}

	info := &SearchJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return nil, fmt.Errorf("获取酷我歌曲信息失败")
	}
	if info.Code == 200 {
		for index, result := range info.Data.List {
			downloadUrl, _ := getPlayURL(result.Rid)
			LrcData, _ := nhttp.Get("https://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=" + strconv.Itoa(result.Rid))
			ret = append(ret, comm.Result{
				Title:    strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Artist + " ]",
				Author:   result.Artist,
				SongName: result.Name,
				SongURL:  downloadUrl,
				LrcData:  LrcData,
				ImgURL:   result.Pic120,
				PicURL:   result.Pic,
			})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("获取酷我数据失败")
}

func getPlayURL(rid int) (playURL string, err error) {
	formats := []string{"320kmp3", "192kmp3", "128kmp3"}
	// kuwo有点奇怪，api时不时就没响应
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	for _, format := range formats {
		u := fmt.Sprintf("https://www.kuwo.cn/api/v1/www/music/playUrl?format=mp3&mid=%d&type=convert_url3&br=%s", rid, format)
		body, err := get(client, u)
		if err != nil {
			continue
		}
		info := &GetPlayURLJSONStruct{}
		err = json.Unmarshal(body, &info)
		if err != nil {
			continue
		}
		if info.Code == 200 {
			return info.Data.URL, nil
		}
	}
	return "", fmt.Errorf("获取歌曲链接失败")
}

func get(client *http.Client, u string) (body []byte, err error) {
	req, _ := http.NewRequest("GET", u, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	fullURL := fmt.Sprintf("http://www.kuwo.cn/api/www/bang/bang/musicList?bangId=93&pn=1&rn=30&httpsStatus=1")
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "kuwo.cn")
	req.Header.Add("Referer", "http://kuwo.cn/")
	req.Header.Add("csrf", "4IT871VN3DA")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")
	req.Header.Add("Cookie", "kw_token=4IT871VN3DA")
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
	for index, result := range info.Data.MusicList {
		downloadUrl, _ := getPlayURL(result.Rid)
		LrcData, _ := nhttp.Get("https://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=" + strconv.Itoa(result.Rid))
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Artist + " ]", Author: result.Artist,
			SongName: result.Name,
			SongURL:  downloadUrl,
			LrcData:  LrcData,
			ImgURL:   result.Pic120,
			PicURL:   result.Pic,
		})
	}
	return ret, nil
}
