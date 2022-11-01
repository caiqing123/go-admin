package comm

import (
	"net/http"
	"time"
)

var (
	// Client 用于请求
	Client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// Result 所有搜索成功后都返回这个
type Result struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	SongName       string `json:"song_name"`
	SongURL        string `json:"song_url"`
	LrcData        string `json:"Lrc_data"`
	ImgURL         string `json:"img_url"`
	PicURL         string `json:"pic_url"`
	SongSilkBase64 string `json:"song_silk_base64"`
}
