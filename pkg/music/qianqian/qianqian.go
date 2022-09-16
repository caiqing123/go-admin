package qianqian

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"api/pkg/music/comm"
)

// Qianqian 千千音乐
func Qianqian(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	u, _ := url.ParseRequestURI("https://music.91q.com/v1/search")
	md5Ctx := md5.New()
	query := url.Values{}
	query.Set("appid", "16073360")
	query.Set("pageSize", "30")
	query.Set("type", "1")
	query.Set("pageNo", p)
	query.Set("word", songName)
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	get, _ := url.QueryUnescape(query.Encode())
	md5Ctx.Write([]byte(get + "0b50b02fd0d73a9c4c8c3a781c30845f"))

	query.Set("sign", hex.EncodeToString(md5Ctx.Sum(nil)))
	u.RawQuery = query.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Referer", "https://music.91q.com/")
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
		return nil, fmt.Errorf("获取数据失败")
	}

	info := &SearchJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return nil, fmt.Errorf("解析信息失败")
	}
	if info.State == true {
		for index, result := range info.Data.TypeTrack {
			downloadUrl, _ := getPlayURL(result.TSID)
			Author := ""
			if len(result.Artist) == 0 {
				Author = "未知"
			} else {
				Author = result.Artist[0].Name
			}
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Title + " - [ " + Author + " ]", Author: Author,
				SongName: result.Title,
				SongURL:  downloadUrl})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("获取数据失败")
}

func getPlayURL(rid string) (playURL string, err error) {
	u, _ := url.ParseRequestURI("https://music.taihe.com/v1/song/tracklink")
	md5Ctx := md5.New()
	query := url.Values{}
	query.Set("appid", "16073360")
	query.Set("TSID", rid)
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	md5Ctx.Write([]byte(query.Encode() + "0b50b02fd0d73a9c4c8c3a781c30845f"))
	query.Set("sign", hex.EncodeToString(md5Ctx.Sum(nil)))
	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Referer", "ttps://music.taihe.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")

	resp, err := comm.Client.Do(req)
	if err != nil {
		log.Printf("request err: %v\n", err)
		return "", fmt.Errorf("请求错误")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body err: %v\n", err)
		return "", fmt.Errorf("获取数据失败")
	}

	info := &GetPlayURLJSONStruct{}
	err = json.Unmarshal(body, &info)

	if info.State == true {
		if info.Data.Path != "" {
			return info.Data.Path, nil
		}
		return info.Data.TrailAudioInfo.Path, nil
	}

	return "", fmt.Errorf("获取歌曲链接失败")
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	u, _ := url.ParseRequestURI("https://music.91q.com/v1/bd/list")
	md5Ctx := md5.New()
	query := url.Values{}
	query.Set("appid", "16073360")
	query.Set("bdid", "257852")
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	get, _ := url.QueryUnescape(query.Encode())
	md5Ctx.Write([]byte(get + "0b50b02fd0d73a9c4c8c3a781c30845f"))

	query.Set("sign", hex.EncodeToString(md5Ctx.Sum(nil)))
	u.RawQuery = query.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Referer", "https://music.91q.com/")
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
	info := &CommendJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return nil, fmt.Errorf("获取json信息失败")
	}

	for index, result := range info.Data.Result {
		downloadUrl, _ := getPlayURL(result.TSID)
		Author := ""
		if len(result.Artist) == 0 {
			Author = "未知"
		} else {
			Author = result.Artist[0].Name
		}
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Title + " - [ " + Author + " ]", Author: Author,
			SongName: result.Title,
			SongURL:  downloadUrl})
	}
	return ret, nil
}
