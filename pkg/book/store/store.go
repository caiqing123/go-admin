package store

import (
	"time"
)

// Store is store yaml data file format
type Store struct {
	BookURL      string
	BookName     string
	Author       string    // 作者
	CoverURL     string    // 封面链接
	Description  string    // 介绍
	LastUpdate   time.Time `yaml:",omitempty"` // 数据更新时间
	DownloadURL  string
	CacheLoadURL string
	Volumes      []Volume
}

// Volume 卷
type Volume struct {
	Name     string
	IsVIP    bool
	Chapters []Chapter
}

// Chapter 章节
type Chapter struct {
	Name  string
	URL   string
	IsVIP bool `yaml:",omitempty"`
	Text  []string
	//MuxLock sync.Mutex `yaml:"-"`
}
