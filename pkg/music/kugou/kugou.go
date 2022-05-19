package kugou

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"api/pkg/music"
)

// Kugou kugou音乐
func Kugou(songName string) (ret []music.Result, err error) {
	fullURL := fmt.Sprintf("http://msearchcdn.kugou.com/api/v3/search/song?pagesize=30&keyword=%s", songName)
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "msearchcdn.kugou.com")
	resp, err := music.Client.Do(req)
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

	info := &SearchJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return
	}
	if info.Status == 1 && info.Error == "" {
		for index, result := range info.Data.Info {
			downloadUrl, _ := getPlayURL(result.Hash)
			ret = append(ret, music.Result{Title: strconv.Itoa(index+1) + ". " + result.Songname + " - [ " + result.Singername + " ]", Author: result.Singername,
				SongName: result.Songname,
				SongURL:  downloadUrl})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}

func getPlayURL(hash string) (playURL string, err error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(hash + "kgcloudv2"))
	fullURL := fmt.Sprintf("http://trackercdn.kugou.com/i/v2/?hash=%s&key=%s&pid=3&behavior=play&cmd=25&version=8990", hash, hex.EncodeToString(md5Ctx.Sum(nil)))
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "trackercdn.kugou.com")
	resp, err := music.Client.Do(req)
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

	info := &SongInfoJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return
	}
	if info.Status == 1 {
		return info.Url[0], nil
	}
	return "", fmt.Errorf("not found")
}

//func getPlayURL(hash string) (playURL string, err error) {
//	fullURL := fmt.Sprintf("https://wwwapi.kugou.com/yy/index.php?r=play/getdata&dfid=2kuKRO3GStCZ0VBY9V12pXeT&mid=f679eeece44cf6bec74d2867be4901f7&hash=%s", hash)
//
//	req, _ := http.NewRequest("GET", fullURL, nil)
//	req.Header.Add("Host", "wwwapi.kugou.com")
//
//	resp, err := music.Client.Do(req)
//	if err != nil {
//		log.Printf("request err: %v\n", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Printf("read body err: %v\n", err)
//		return
//	}
//
//	info := &SongInfoJSONStruct{}
//	err = json.Unmarshal(body, &info)
//	if err != nil {
//		fmt.Println(string(body))
//		log.Printf("json unmarshal err: %v\n", err)
//		return
//	}
//	if info.Status == 1 && info.ErrCode == 0 {
//		return info.Data.PlayURL, nil
//	}
//	return "", fmt.Errorf("not found")
//}
