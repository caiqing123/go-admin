package netease

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cast"

	nhttp "api/pkg/http"
	"api/pkg/music/comm"
)

// Netease 网易云
func Netease(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	u, _ := url.ParseRequestURI("http://music.163.com/api/search/pc")
	query := url.Values{}
	query.Set("s", songName)
	offset := 20 * (cast.ToInt(p) - 1)
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", "30")
	query.Set("type", "1")
	u.RawQuery = query.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Host", "music.163.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "NMTID=00OgriIBhmBRHmA_UqJj15PIKu9U68AAAGCH-eGGQ")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")

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

	info := &JSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return nil, fmt.Errorf("获取json信息失败")
	}
	if info.Code == 200 && len(info.Result.Songs) != 0 {
		for index, result := range info.Result.Songs {
			downloadUrl := fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%d.mp3", result.ID)
			Data, _ := nhttp.Get("https://music.163.com/api/song/lyric?lv=-1&tv=-1&id=" + strconv.Itoa(result.ID))
			LrcData := &LrcDataJSONStruct{}
			_ = json.Unmarshal([]byte(Data), &LrcData)
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Artists[0].Name + " ]", Author: result.Artists[0].Name,
				SongName: result.Name,
				SongURL:  downloadUrl,
				LrcData:  LrcData.Lrc.Lyric,
				ImgURL:   result.Album.PicURL,
				PicURL:   result.Album.PicURL,
			})
		}
		return ret, nil
	}
	return
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	api := "https://music.163.com/api/v6/playlist/detail"
	values := url.Values{}
	values.Set("id", "19723756")
	values.Set("n", "30")

	payload := strings.NewReader(values.Encode())
	req, _ := http.NewRequest("POST", api, payload)
	req.Header.Add("Host", "music.163.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 Edg/84.0.522.48")
	req.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
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

	for index, result := range info.Playlist.Tracks {
		downloadUrl := fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%d.mp3", result.ID)
		Data, _ := nhttp.Get("https://music.163.com/api/song/lyric?lv=-1&tv=-1&id=" + strconv.Itoa(result.ID))
		LrcData := &LrcDataJSONStruct{}
		_ = json.Unmarshal([]byte(Data), &LrcData)
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Ar[0].Name + " ]", Author: result.Ar[0].Name,
			SongName: result.Name,
			SongURL:  downloadUrl,
			LrcData:  LrcData.Lrc.Lyric,
			ImgURL:   result.Al.PicUrl,
			PicURL:   result.Al.PicUrl,
		})
	}
	return ret, nil
}
