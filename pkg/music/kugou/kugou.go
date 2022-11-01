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

	nhttp "api/pkg/http"
	"api/pkg/music/comm"
)

// Kugou kugou音乐
func Kugou(songName string) (ret []comm.Result, err error) {
	fullURL := fmt.Sprintf("http://msearchcdn.kugou.com/api/v3/search/song?pagesize=30&keyword=%s", songName)
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "msearchcdn.kugou.com")
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
		log.Println("获取json信息失败")
		return nil, fmt.Errorf("获取json信息失败")
	}
	if info.Status == 1 && info.Error == "" {
		for index, result := range info.Data.Info {
			downloadUrl, _ := getPlayURL(result.Hash)
			if downloadUrl == "" {
				continue
			}
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Songname + " - [ " + result.Singername + " ]", Author: result.Singername,
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

func NewKugou(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	fullURL := fmt.Sprintf("https://songsearch.kugou.com/song_search_v2?keyword=%s&page=%s", songName, p)
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "www.kugou.com")
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

	info := &NewSearchJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return nil, fmt.Errorf("获取json信息失败")
	}
	if info.Status == 1 && info.ErrorMsg == "" {
		for index, result := range info.Data.Lists {
			downloadUrl, _ := newGetPlayURL(result.AlbumID, result.FileHash)
			Data, _ := nhttp.Get(fmt.Sprintf("https://wwwapi.kugou.com/yy/index.php?r=play/getdata&hash=%s&album_id=%s&mid=1&platid=4", result.FileHash, result.AlbumID))
			LrcData := &LrcDataJSONStruct{}
			_ = json.Unmarshal([]byte(Data), &LrcData)
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.SongName + " - [ " + result.SingerName + " ]", Author: result.SingerName,
				SongName: result.SongName,
				SongURL:  downloadUrl,
				LrcData:  LrcData.Data.Lyrics,
				ImgURL:   LrcData.Data.Img,
				PicURL:   LrcData.Data.Img,
			})
		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}

func newGetPlayURL(albumID string, fileHash string) (playURL string, err error) {
	fullURL := fmt.Sprintf("https://m.kugou.com/app/i/getSongInfo.php?cmd=playInfo&hash=%s", fileHash)
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "www.kugou.com")
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
	info := &NewSongInfoJSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("json unmarshal err: %v\n", err)
		return
	}
	if info.Status == 1 {
		if info.Url != "" {
			return info.Url, nil
		}
	}
	return "", fmt.Errorf("not found")
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	fullURL := fmt.Sprintf("https://m.kugou.com/?json=true")
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Host", "www.kugou.com")
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
	for index, result := range info.Data {
		downloadUrl, _ := newGetPlayURL(result.AlbumID, result.Hash)
		Data, _ := nhttp.Get(fmt.Sprintf("https://wwwapi.kugou.com/yy/index.php?r=play/getdata&hash=%s&album_id=%s&mid=1&platid=4", result.Hash, result.AlbumID))
		LrcData := &LrcDataJSONStruct{}
		_ = json.Unmarshal([]byte(Data), &LrcData)
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.SongName + " - [ " + result.SingerName + " ]", Author: result.SingerName,
			SongName: result.SongName,
			SongURL:  downloadUrl,
			LrcData:  LrcData.Data.Lyrics,
			ImgURL:   LrcData.Data.Img,
			PicURL:   LrcData.Data.Img,
		})
	}
	return ret, nil
}
