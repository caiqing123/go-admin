package migu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"api/pkg/music/comm"
)

// Migu 咪咕音乐
func Migu(songName string, p string) (ret []comm.Result, err error) {
	if songName == "" {
		return commend()
	}
	u, _ := url.ParseRequestURI("http://pd.musicapp.migu.cn/MIGUM3.0/v1.0/content/search_all.do")
	query := url.Values{}
	query.Set("pageNo", p)
	query.Set("pageSize", "30")
	query.Set("searchSwitch", `{"song":1,"album":0,"singer":0,"tagSong":0,"mvSong":0,"songlist":0,"bestShow":0}`)
	query.Set("text", songName)
	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("Referer", "https://m.music.migu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Mobile Safari/537.36")
	req.Header.Add("Host", "https://migu.cn")

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

	if info.Code == "000000" && len(info.SongResultData.Result) > 0 {
		for index, result := range info.SongResultData.Result {
			list := make(map[int]string)
			for i, v := range result.NewRateFormats {
				if v.FileType == "mp3" {
					list[i] = v.URL
				}
			}
			option := list[len(list)-1]
			pathname := &url.URL{}
			pathname, err = url.Parse(option)
			downloadUrl := "https://freetyst.nf.migu.cn/" + pathname.Path
			ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.Name + " - [ " + result.Singers[0].Name + " ]", Author: result.Singers[0].Name,
				SongName: result.Name,
				SongURL:  downloadUrl})

		}
		return ret, nil
	}
	return nil, fmt.Errorf("not found")
}

//commend 推荐
func commend() (ret []comm.Result, err error) {
	fullURL := fmt.Sprintf("https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/querycontentbyId.do?columnId=27553319&needAll=0")
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Referer", "https://m.music.migu.cn/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Mobile Safari/537.36")
	req.Header.Add("Host", "https://migu.cn")
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
	for index, result := range info.ColumnInfo.Contents {
		list := make(map[int]string)
		for i, v := range result.ObjectInfo.NewRateFormats {
			if v.FileType == "mp3" {
				list[i] = v.URL
			}
		}
		option := list[len(list)-1]
		pathname := &url.URL{}
		pathname, err = url.Parse(option)
		downloadUrl := "https://freetyst.nf.migu.cn/" + pathname.Path
		ret = append(ret, comm.Result{Title: strconv.Itoa(index+1) + ". " + result.ObjectInfo.SongName + " - [ " + result.ObjectInfo.Singer + " ]", Author: result.ObjectInfo.Singer,
			SongName: result.ObjectInfo.SongName,
			SongURL:  downloadUrl})
	}
	return ret, nil
}
