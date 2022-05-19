package migu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"api/pkg/music"
)

// Migu 咪咕音乐
func Migu(songName string) (ret []music.Result, err error) {
	u, _ := url.ParseRequestURI("http://pd.musicapp.migu.cn/MIGUM3.0/v1.0/content/search_all.do")
	query := url.Values{}
	query.Set("pageNo", "1")
	query.Set("pageSize", "30")
	query.Set("searchSwitch", `{"song":1,"album":0,"singer":0,"tagSong":0,"mvSong":0,"songlist":0,"bestShow":0}`)
	query.Set("text", songName)
	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("Referer", "https://m.music.migu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Mobile Safari/537.36")
	req.Header.Add("Host", "https://migu.cn")

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
	info := &JSONStruct{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("获取json信息失败")
		return
	}

	if info.Code == "000000" && len(info.SongResultData.Result) > 0 {
		for index, result := range info.SongResultData.Result {
			option := result.NewRateFormats[len(result.NewRateFormats)-1]
			pathname := &url.URL{}
			if option.AndroidURL != "" {
				pathname, err = url.Parse(option.AndroidURL)
			} else {
				pathname, err = url.Parse(option.URL)

			}
			downloadUrl := "https://freetyst.nf.migu.cn/" + pathname.Path

			ret = append(ret, music.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Singers[0].Name + " ]", Author: result.Singers[0].Name,
				SongName: result.Name,
				SongURL:  downloadUrl})

		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}
