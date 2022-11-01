package qianqian

// SearchJSONStruct 搜索
type SearchJSONStruct struct {
	State bool `json:"state"`
	Data  struct {
		Total     int `json:"total"`
		TypeTrack []struct {
			TSID       string `json:"TSID"`
			AlbumTitle string `json:"albumTitle"`
			Artist     []struct {
				ArtistTypeName string `json:"artistTypeName"`
				Name           string `json:"name"`
				Gender         string `json:"gender"`
				Pic            string `json:"pic"`
			} `json:"artist"`
			Title string `json:"title"`
			Pic   string `json:"pic"`
			Lyric string `json:"lyric"`
			Id    string `json:"id"`
			IsVip int    `json:"isVip"`
		} `json:"typeTrack"`
	} `json:"data"`
	Errmsg string `json:"errmsg"`
	Errno  int    `json:"errno"`
}

// GetPlayURLJSONStruct 获取播放链接
type GetPlayURLJSONStruct struct {
	State  bool   `json:"state"`
	Errmsg string `json:"errmsg"`
	Errno  int    `json:"errno"`
	Data   struct {
		Path           string `json:"path"`
		TrailAudioInfo struct {
			Path string `json:"path"`
		} `json:"trail_audio_info"`
	} `json:"data"`
}

type CommendJSONStruct struct {
	State bool `json:"state"`
	Data  struct {
		Title  string `json:"title"`
		Result []struct {
			TSID       string `json:"TSID"`
			AlbumTitle string `json:"albumTitle"`
			Artist     []struct {
				ArtistTypeName string `json:"artistTypeName"`
				Name           string `json:"name"`
				Gender         string `json:"gender"`
				Pic            string `json:"pic"`
			} `json:"artist"`
			Pic   string `json:"pic"`
			Lyric string `json:"lyric"`
			Title string `json:"title"`
			Id    string `json:"id"`
			IsVip int    `json:"isVip"`
		} `json:"result"`
	} `json:"data"`
	Errmsg string `json:"errmsg"`
	Errno  int    `json:"errno"`
}
