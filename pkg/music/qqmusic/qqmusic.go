package qqmusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"api/pkg/music"
)

// QQMusic qq音乐
func QQMusic(songName string) (ret []music.Result, err error) {
	searchURL := "https://c.y.qq.com/soso/fcgi-bin/client_search_cp"
	fullURL := fmt.Sprintf("%s?t=0&cr=1&p=1&n=30&w=%s&format=json&aggr=1&lossless=1&new_json=1", searchURL, songName)
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))
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
		log.Printf("json unmarshal err: %v\n", err)
		return
	}
	if info.Code == 0 {
		for index, result := range info.Data.Song.List {
			downloadUrl, _ := getPlayURL(result.Mid)
			ret = append(ret, music.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Singer[0].Name + " ]", Author: result.Singer[0].Name,
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
	param := "{%22req%22:%20{%22module%22:%20%22CDN.SrfCdnDispatchServer%22,%20%22method%22:%20%22GetCdnDispatch%22,%20%22param%22:%20{%22guid%22:%20%223982823384%22,%20%22calltype%22:%200,%20%22userip%22:%20%22%22}},%20%22req_0%22:%20{%22module%22:%20%22vkey.GetVkeyServer%22,%20%22method%22:%20%22CgiGetVkey%22,%20%22param%22:%20{%22guid%22:%20%223982823384%22,%20%22songmid%22:%20[%22" + mid + "%22],%20%22songtype%22:%20[0],%20%22uin%22:%20%220%22,%20%22loginflag%22:%201,%20%22platform%22:%20%2220%22}},%20%22comm%22:%20{%22uin%22:%200,%20%22format%22:%20%22json%22,%20%22ct%22:%2024,%20%22cv%22:%200}}"

	fullURL := vkeyURL + "?data=" + param
	req, _ := http.NewRequest("GET", fullURL, strings.NewReader(""))

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
